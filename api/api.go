// Package api holds the shared HTTP response helpers: a consistent error
// envelope, status helpers, and the request-ID accessor.
package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestIDKey is the Gin context key under which the request ID is stored.
const RequestIDKey = "request_id"

// ErrorBody is the JSON shape returned for every error response.
type ErrorBody struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// RequestID returns the request ID attached by the RequestID middleware.
func RequestID(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}

// Error writes a consistent error response and aborts the handler chain.
func Error(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, ErrorBody{Error: message, Code: code})
}

// BadRequest reports malformed input.
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "bad_request", message)
}

// Unauthorized reports a missing or invalid credential.
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "unauthorized", message)
}

// Forbidden reports an authenticated caller acting on a resource they do not own.
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "forbidden", message)
}

// NotFound reports a missing resource.
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "not_found", message)
}

// Conflict reports a request that clashes with existing state.
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, "conflict", message)
}

// Unavailable reports that a dependency needed to serve the request is down.
func Unavailable(c *gin.Context, message string) {
	Error(c, http.StatusServiceUnavailable, "unavailable", message)
}

// Internal logs the underlying error and returns a generic message so database
// and other implementation details never reach the client.
func Internal(c *gin.Context, err error) {
	slog.Error("internal error",
		"error", err,
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"request_id", RequestID(c),
	)
	Error(c, http.StatusInternalServerError, "internal_error", "Something went wrong")
}
