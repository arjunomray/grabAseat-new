package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"arjunomray/eventer/api"
	"arjunomray/eventer/dbConfig"
)

// Healthz reports that the process is running (liveness probe).
func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Readyz reports whether the service can serve traffic (readiness probe).
// Currently that means the database is reachable.
func Readyz(c *gin.Context) {
	sqlDB, err := dbConfig.Db.DB()
	if err != nil {
		slog.Warn("readiness check failed", "error", err, "request_id", api.RequestID(c))
		api.Unavailable(c, "database unavailable")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		slog.Warn("readiness check failed", "error", err, "request_id", api.RequestID(c))
		api.Unavailable(c, "database unavailable")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
