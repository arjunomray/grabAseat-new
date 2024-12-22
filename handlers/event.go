package handlers

import (
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.ID = uuid.New()

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
