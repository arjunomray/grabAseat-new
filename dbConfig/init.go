package dbConfig

import (
	"log"
	"time"

	"arjunomray/eventer/config"
	"arjunomray/eventer/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDatabase(cfg config.Database) {
	var err error
	for i := range 5 {
		Db, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect db, retrying in 5 seconds (attempt %d/5)", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	Db.AutoMigrate(&models.User{}, &models.Event{}, &models.Ticket{})
	log.Print("Database Connected")
}
