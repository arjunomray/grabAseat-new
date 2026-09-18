package handlers

import (
	"arjunomray/eventer/api"
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/middleware"
	"arjunomray/eventer/models"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// isEventAdmin reports whether userID is the event's creator or an added admin.
func isEventAdmin(eventID, userID uuid.UUID) (bool, error) {
	var event models.Event
	if err := dbConfig.Db.First(&event, "id = ?", eventID).Error; err != nil {
		return false, err
	}
	if event.CreatedBy == userID {
		return true, nil
	}

	var count int64
	if err := dbConfig.Db.Model(&models.EventAdmin{}).
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// loadEvent fetches an event or writes a 404/500 and returns false.
func loadEvent(c *gin.Context, eventID uuid.UUID) (models.Event, bool) {
	var event models.Event
	if err := dbConfig.Db.First(&event, "id = ?", eventID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.NotFound(c, "Event not found")
		} else {
			api.Internal(c, err)
		}
		return event, false
	}
	return event, true
}

type eventAdminView struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	AddedAt time.Time `json:"addedAt"`
}

// ListEventAdmins returns the creator plus the added admins. Only admins of the
// event may view the list.
func ListEventAdmins(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok {
		api.Unauthorized(c, "Unauthorized")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.BadRequest(c, "Invalid event ID")
		return
	}

	event, ok := loadEvent(c, eventID)
	if !ok {
		return
	}

	admin, err := isEventAdmin(eventID, claims.UserID)
	if err != nil {
		api.Internal(c, err)
		return
	}
	if !admin {
		api.Forbidden(c, "You are not an admin of this event")
		return
	}

	var creator models.User
	if err := dbConfig.Db.First(&creator, "id = ?", event.CreatedBy).Error; err != nil {
		api.Internal(c, err)
		return
	}

	var admins []models.EventAdmin
	if err := dbConfig.Db.Preload("User").
		Where("event_id = ?", eventID).
		Order("created_at").
		Find(&admins).Error; err != nil {
		api.Internal(c, err)
		return
	}

	views := make([]eventAdminView, 0, len(admins))
	for _, admin := range admins {
		if admin.User == nil {
			continue
		}
		views = append(views, eventAdminView{
			ID:      admin.User.ID,
			Name:    admin.User.Name,
			Email:   admin.User.Email,
			AddedAt: admin.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"creator":   gin.H{"id": creator.ID, "name": creator.Name, "email": creator.Email},
		"admins":    views,
		"canManage": event.CreatedBy == claims.UserID,
	})
}

type AddEventAdminInput struct {
	Email string `json:"email" binding:"required,email"`
}

// AddEventAdmin adds an existing user, by email, as an admin of the event.
// Only the event creator may do this.
func AddEventAdmin(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok {
		api.Unauthorized(c, "Unauthorized")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.BadRequest(c, "Invalid event ID")
		return
	}

	event, ok := loadEvent(c, eventID)
	if !ok {
		return
	}
	if event.CreatedBy != claims.UserID {
		api.Forbidden(c, "Only the event creator can manage admins")
		return
	}

	var input AddEventAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		api.BadRequest(c, err.Error())
		return
	}
	email := strings.ToLower(strings.TrimSpace(input.Email))

	var user models.User
	if err := dbConfig.Db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.NotFound(c, "No user with that email")
			return
		}
		api.Internal(c, err)
		return
	}
	if user.ID == event.CreatedBy {
		api.Conflict(c, "That user is the event creator and is already an admin")
		return
	}

	admin := models.EventAdmin{EventID: eventID, UserID: user.ID}
	if err := dbConfig.Db.Create(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			api.Conflict(c, "That user is already an admin")
			return
		}
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusCreated, eventAdminView{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		AddedAt: admin.CreatedAt,
	})
}

// RemoveEventAdmin removes an added admin. Only the event creator may do this,
// and the creator can never be removed.
func RemoveEventAdmin(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok {
		api.Unauthorized(c, "Unauthorized")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.BadRequest(c, "Invalid event ID")
		return
	}
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		api.BadRequest(c, "Invalid user ID")
		return
	}

	event, ok := loadEvent(c, eventID)
	if !ok {
		return
	}
	if event.CreatedBy != claims.UserID {
		api.Forbidden(c, "Only the event creator can manage admins")
		return
	}
	if userID == event.CreatedBy {
		api.BadRequest(c, "The event creator cannot be removed")
		return
	}

	res := dbConfig.Db.Where("event_id = ? AND user_id = ?", eventID, userID).Delete(&models.EventAdmin{})
	if res.Error != nil {
		api.Internal(c, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		api.NotFound(c, "That user is not an admin of this event")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin removed"})
}
