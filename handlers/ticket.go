package handlers

import (
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTicket(c *gin.Context) {
	var ticket models.Ticket

	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingTicket models.Ticket
	if err := dbConfig.Db.Where("event_id = ? AND user_id = ?", ticket.EventID, ticket.UserID).First(&existingTicket).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket already exists"})
		return
	}

	ticket.ID = uuid.New()
	if err := dbConfig.Db.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func GetTickets(c *gin.Context) {
	var tickets []models.Ticket
	if err := dbConfig.Db.Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func DeleteTickets(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dbConfig.Db.Delete(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted"})
}

func GetUserTickets(c *gin.Context) {
	log.Print("Here")
	userIDParam := c.Param("id")

	userID, err := uuid.Parse(userIDParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := dbConfig.Db.Preload("Tickets").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	log.Print(user.Tickets)

	c.JSON(http.StatusOK, user.Tickets)
}
