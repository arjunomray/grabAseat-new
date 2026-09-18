package dbConfig

import (
	"log/slog"
	"os"
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
		slog.Warn("failed to connect to database, retrying",
			"attempt", i+1,
			"max_attempts", 5,
			"error", err,
		)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	Db.AutoMigrate(&models.User{}, &models.Event{}, &models.Ticket{})
	slog.Info("database connected")
}
