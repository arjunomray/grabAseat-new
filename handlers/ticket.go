package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"arjunomray/eventer/api"
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/middleware"
	"arjunomray/eventer/models"
)

var (
	errEventPast = errors.New("event has already passed")
	errSoldOut   = errors.New("event is sold out")
	errDuplicate = errors.New("ticket already exists")
)

type CreateTicketInput struct {
	EventID uuid.UUID `json:"event_id" binding:"required"`
}

func CreateTicket(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok {
		api.Unauthorized(c, "Unauthorized")
		return
	}

	var input CreateTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	var ticket models.Ticket
	err := dbConfig.Db.Transaction(func(tx *gorm.DB) error {
		// Lock the event row so concurrent bookings for the same event
		// serialize, making the seat check and decrement atomic.
		var event models.Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", input.EventID).
			First(&event).Error; err != nil {
			return err
		}

		if event.Time.Before(time.Now()) {
			return errEventPast
		}
		if event.SeatsAvailable <= 0 {
			return errSoldOut
		}

		var existing models.Ticket
		err := tx.Where("event_id = ? AND user_id = ?", event.ID, claims.UserID).First(&existing).Error
		if err == nil {
			return errDuplicate
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := tx.Model(&models.Event{}).
			Where("id = ?", event.ID).
			Update("seats_available", gorm.Expr("seats_available - 1")).Error; err != nil {
			return err
		}

		ticket = models.Ticket{
			ID:      uuid.New(),
			EventID: event.ID,
			UserID:  claims.UserID,
		}
		return tx.Create(&ticket).Error
	})

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		api.NotFound(c, "Event not found")
		return
	case errors.Is(err, errEventPast):
		api.BadRequest(c, "Event has already passed")
		return
	case errors.Is(err, errSoldOut):
		api.Conflict(c, "Event is sold out")
		return
	case errors.Is(err, errDuplicate), errors.Is(err, gorm.ErrDuplicatedKey):
		api.Conflict(c, "Ticket already exists")
		return
	case err != nil:
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func GetTickets(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok {
		api.Unauthorized(c, "Unauthorized")
		return
	}

	var tickets []models.Ticket
	if err := dbConfig.Db.Where("user_id = ?", claims.UserID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		api.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func DeleteTicket(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok {
		api.Unauthorized(c, "Unauthorized")
		return
	}

	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		api.BadRequest(c, "Invalid ticket ID")
		return
	}

	var ticket models.Ticket
	if err := dbConfig.Db.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.NotFound(c, "Ticket not found")
			return
		}
		api.Internal(c, err)
		return
	}

	if ticket.UserID != claims.UserID {
		api.Forbidden(c, "You can only delete your own tickets")
		return
	}

	if err := dbConfig.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&ticket).Error; err != nil {
			return err
		}
		return tx.Model(&models.Event{}).
			Where("id = ?", ticket.EventID).
			Update("seats_available", gorm.Expr("seats_available + 1")).Error
	}); err != nil {
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

	claims, ok := middleware.ClaimsFrom(c)
	if !ok {
		api.Unauthorized(c, "Unauthorized")
		return
	}
	if claims.UserID != userID {
		api.Forbidden(c, "You can only access your own resources")
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
