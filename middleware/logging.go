package middleware

import (
	"arjunomray/eventer/api"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID ensures every request has an ID. It reuses the incoming
// X-Request-ID header when present and generates one otherwise, storing it in
// the request context and echoing it back on the response.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(api.RequestIDKey, id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// Logger writes one structured log line per request.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		attrs := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", api.RequestID(c),
			"client_ip", c.ClientIP(),
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		if c.Writer.Status() >= http.StatusInternalServerError {
			slog.Error("request", attrs...)
			return
		}
		slog.Info("request", attrs...)
	}
}

// Recovery logs panics through slog and returns a generic 500 response instead
// of Gin's default HTML/text output.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error(
					"panic recovered",
					"error", recovered,
					"path", c.Request.URL.Path,
					"request_id", api.RequestID(c),
				)
				api.Error(c, http.StatusInternalServerError, "internal_error", "internal server error")
			}
		}()
		c.Next()
	}
}
