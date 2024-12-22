package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Models
type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"password"`
	Tickets  []Ticket  `json:"tickets" gorm:"foreignKey:UserID"`
	Events   []Event   `json:"events" gorm:"many2many:user_events;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

type Event struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name"`
	Time      string    `json:"time"`
	CreatedBy string    `json:"createdBy"`
	Users     []User    `json:"users" gorm:"many2many:user_events;"`
	Tickets   []Ticket  `json:"tickets" gorm:"foreignKey:EventID"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

type Ticket struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	EventID uuid.UUID `json:"event_id" gorm:"type:uuid"`
	UserID  uuid.UUID `json:"user_id" gorm:"type:uuid"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}