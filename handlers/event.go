package handlers

import (
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/models"
	"arjunomray/eventer/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventInput struct {
	Name           string   `json:"name"`
	Time           string   `json:"time"`
	Location       string   `json:"location"`
	Price          float64  `json:"price"`
	SeatsAvailable int      `json:"seatsAvailable"`
	Tags           []string `json:"tags"`
	CreatedBy      string   `json:"createdBy"`
}

func CreateEvent(c *gin.Context) {
	var eventInput EventInput
	if err := c.ShouldBindJSON(&eventInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := models.Event{
		Name:           eventInput.Name,
		Time:           eventInput.Time,
		Location:       eventInput.Location,
		Price:          eventInput.Price,
		SeatsAvailable: eventInput.SeatsAvailable,
		Tags:           eventInput.Tags,
		CreatedBy:      eventInput.CreatedBy,
	}

	// Get user ID from JWT token
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims := claims.(*utils.Claims)

	event.ID = uuid.New()
	event.CreatedBy = userClaims.UserID.String() // Set the creator's ID

	if err := dbConfig.Db.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

func GetEvents(c *gin.Context) {
	var events []models.Event
	if err := dbConfig.Db.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func GetEvent(c *gin.Context) {
	eventIDParam := c.Param("id")
	eventID := uuid.MustParse(eventIDParam)
	var event models.Event
	if err := dbConfig.Db.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}
