package dbConfig

import (
	"arjunomray/eventer/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	var err error
	for i := range 5 {

		Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
