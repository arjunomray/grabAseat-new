package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Models
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Tickets   []Ticket  `json:"tickets" gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

type Event struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedBy      uuid.UUID      `json:"createdBy" gorm:"type:uuid"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Time           time.Time      `json:"time" gorm:"type:timestamptz"`
	Location       string         `json:"location"`
	PriceCents     int64          `json:"priceCents" gorm:"column:price_cents"`
	SeatsAvailable int            `json:"seatsAvailable"`
	Tags           pq.StringArray `json:"tags" gorm:"type:text[]"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Tickets        []Ticket       `json:"tickets" gorm:"foreignKey:EventID"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

type Ticket struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	EventID   uuid.UUID `json:"event_id" gorm:"type:uuid"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Event     *Event    `json:"event,omitempty" gorm:"foreignKey:EventID"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
