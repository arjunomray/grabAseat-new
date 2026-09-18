// Package router wires the HTTP routes and middleware together.
package router

import (
	"arjunomray/eventer/config"
	"arjunomray/eventer/handlers"
	"arjunomray/eventer/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// New builds the application router.
func New(cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Infrastructure endpoints (unversioned).
	r.GET("/healthz", handlers.Healthz)
	r.GET("/readyz", handlers.Readyz)

	// API v1
	v1 := r.Group("/api/v1")
	v1.POST("/register", handlers.RegisterUser)
	v1.POST("/login", handlers.LoginUser)

	// Public: seat availability stream (EventSource cannot send headers).
	v1.GET("/stream", handlers.StreamEvents)

	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/events", handlers.CreateEvent)
	protected.GET("/events", handlers.GetEvents)
	protected.GET("/events/:id", handlers.GetEvent)
	protected.POST("/tickets", handlers.CreateTicket)
	protected.GET("/tickets", handlers.GetTickets)
	protected.GET("/tickets/me", handlers.GetTickets)
	protected.DELETE("/tickets/:id", handlers.DeleteTicket)
	protected.GET("/users/:id", handlers.GetUser)
	protected.GET("/users/:id/events", handlers.GetUserEvents)
	protected.GET("/users/:id/tickets", handlers.GetUserTickets)

	return r
}
