package handlers

import (
	"errors"
	"net/http"
	"time"

	"arjunomray/eventer/api"
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/middleware"
	"arjunomray/eventer/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventInput struct {
	Name           string    `json:"name" binding:"required"`
	Description    string    `json:"description"`
	Time           time.Time `json:"time" binding:"required"`
	Location       string    `json:"location" binding:"required"`
	PriceCents     int64     `json:"priceCents" binding:"gte=0"`
	SeatsAvailable int       `json:"seatsAvailable" binding:"required,gte=1"`
	Tags           []string  `json:"tags"`
}

func CreateEvent(c *gin.Context) {
	var eventInput EventInput
	if err := c.ShouldBindJSON(&eventInput); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	if !eventInput.Time.After(time.Now()) {
		api.BadRequest(c, "Event time must be in the future")
		return
	}

	// The event owner comes from the JWT, never from the request body.
	userClaims, ok := middleware.ClaimsFrom(c)
	if !ok {
		api.Unauthorized(c, "Unauthorized")
		return
	}

	event := models.Event{
		ID:             uuid.New(),
		Name:           eventInput.Name,
		Description:    eventInput.Description,
		Time:           eventInput.Time,
		Location:       eventInput.Location,
		PriceCents:     eventInput.PriceCents,
		SeatsAvailable: eventInput.SeatsAvailable,
		Tags:           eventInput.Tags,
		CreatedBy:      userClaims.UserID,
	}

	if err := dbConfig.Db.Create(&event).Error; err != nil {
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusCreated, event)
}

func GetEvents(c *gin.Context) {
	var events []models.Event
	if err := dbConfig.Db.Find(&events).Error; err != nil {
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, events)
}

func GetEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.BadRequest(c, "Invalid event ID")
		return
	}

	var event models.Event
	if err := dbConfig.Db.First(&event, eventID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.NotFound(c, "Event not found")
			return
		}
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, event)
}
