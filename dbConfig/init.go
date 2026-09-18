package dbConfig

import (
	"arjunomray/eventer/config"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

// Connect opens a database connection without retrying or exiting, so tests
// and tooling can handle failures themselves.
func Connect(cfg config.Database) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{TranslateError: true})
}

func InitDatabase(cfg config.Database) {
	var err error
	for i := range 5 {
		Db, err = Connect(cfg)
		if err == nil {
			break
		}
		slog.Warn(
			"failed to connect to database, retrying",
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

	slog.Info("database connected")
}
