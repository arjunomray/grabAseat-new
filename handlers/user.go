package handlers

import (
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/models"
	"arjunomray/eventer/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RegisterUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := dbConfig.Db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUserEvents(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get events created by user
	var createdEvents []models.Event
	if err := dbConfig.Db.Where("created_by = ?", userID).Find(&createdEvents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching created events"})
		return
	}

	// Get events user is registered for through tickets
	var registeredEvents []models.Event
	if err := dbConfig.Db.
		Joins("JOIN tickets ON tickets.event_id = events.id").
		Where("tickets.user_id = ?", userID).
		Find(&registeredEvents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching registered events"})
		return
	}

	response := gin.H{
		"created_events":    createdEvents,
		"registered_events": registeredEvents,
	}

	c.JSON(http.StatusOK, response)
}

func RegisterUser(c *gin.Context) {
	var userInput RegisterUserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	err := dbConfig.Db.Where("email = ?", userInput.Email).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		ID:       uuid.New(),
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}

	if err := dbConfig.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func LoginUser(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := dbConfig.Db.Where("email = ? AND password = ?", credentials.Email, credentials.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token, "userId": user.ID})
}
