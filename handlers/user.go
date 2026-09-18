package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"arjunomray/eventer/api"
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/middleware"
	"arjunomray/eventer/models"
	"arjunomray/eventer/utils"
)

// requireSelf returns the authenticated user, or writes a 401/403 and reports
// false when the caller is not the owner of the requested user resource.
func requireSelf(c *gin.Context, userID uuid.UUID) bool {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok {
		api.Unauthorized(c, "Unauthorized")
		return false
	}
	if claims.UserID != userID {
		api.Forbidden(c, "You can only access your own resources")
		return false
	}
	return true
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.BadRequest(c, "Invalid user ID")
		return
	}
	if !requireSelf(c, id) {
		return
	}

	var user models.User
	if err := dbConfig.Db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.NotFound(c, "User not found")
			return
		}
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUserEvents(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.BadRequest(c, "Invalid user ID")
		return
	}
	if !requireSelf(c, userID) {
		return
	}

	// Get events created by user
	var createdEvents []models.Event
	if err := dbConfig.Db.Where("created_by = ?", userID).Find(&createdEvents).Error; err != nil {
		api.Internal(c, err)
		return
	}

	// Get events user is registered for through tickets
	var registeredEvents []models.Event
	if err := dbConfig.Db.
		Joins("JOIN tickets ON tickets.event_id = events.id").
		Where("tickets.user_id = ?", userID).
		Find(&registeredEvents).Error; err != nil {
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"created_events":    createdEvents,
		"registered_events": registeredEvents,
	})
}

func RegisterUser(c *gin.Context) {
	var userInput RegisterUserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	email := strings.ToLower(strings.TrimSpace(userInput.Email))

	var existingUser models.User
	err := dbConfig.Db.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		api.Conflict(c, "User already exists")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		api.Internal(c, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		api.Internal(c, err)
		return
	}

	user := models.User{
		ID:       uuid.New(),
		Name:     userInput.Name,
		Email:    email,
		Password: string(hash),
	}

	if err := dbConfig.Db.Create(&user).Error; err != nil {
		// The unique constraint is the source of truth: handle the race where
		// another request created the same email between the check and insert.
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			api.Conflict(c, "User already exists")
			return
		}
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func LoginUser(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	var user models.User
	err := dbConfig.Db.Where("email = ?", strings.ToLower(strings.TrimSpace(credentials.Email))).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		api.Unauthorized(c, "Invalid email or password")
		return
	}
	if err != nil {
		api.Internal(c, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		api.Unauthorized(c, "Invalid email or password")
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token, "userId": user.ID})
}
