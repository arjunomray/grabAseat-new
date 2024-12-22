package dbConfig

import (
	"arjunomray/eventer/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDatabase() {
	var err error
	Db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	Db.AutoMigrate(&models.User{}, &models.Event{}, &models.Ticket{})
	log.Print("Database Connected")
}
