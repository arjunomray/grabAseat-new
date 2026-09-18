package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"arjunomray/eventer/api"
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/models"
)

func CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	var existingTicket models.Ticket
	err := dbConfig.Db.Where("event_id = ? AND user_id = ?", ticket.EventID, ticket.UserID).First(&existingTicket).Error
	switch {
	case err == nil:
		api.Conflict(c, "Ticket already exists")
		return
	case !errors.Is(err, gorm.ErrRecordNotFound):
		api.Internal(c, err)
		return
	}

	ticket.ID = uuid.New()
	if err := dbConfig.Db.Create(&ticket).Error; err != nil {
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func GetTickets(c *gin.Context) {
	var tickets []models.Ticket
	if err := dbConfig.Db.Find(&tickets).Error; err != nil {
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func DeleteTickets(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	if err := dbConfig.Db.Delete(&ticket).Error; err != nil {
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted"})
}

func GetUserTickets(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.BadRequest(c, "Invalid user ID")
		return
	}

	var user models.User
	if err := dbConfig.Db.Preload("Tickets").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.NotFound(c, "User not found")
			return
		}
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, user.Tickets)
}
