package main

import (
	"log/slog"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"arjunomray/eventer/config"
	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/handlers"
	"arjunomray/eventer/middleware"
	"arjunomray/eventer/utils"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	utils.SetJWTSecret(cfg.JWTSecret)
	dbConfig.InitDatabase(cfg.Database)

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

	// Routes
	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginUser)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/events", handlers.CreateEvent)
	protected.GET("/events", handlers.GetEvents)
	protected.GET("/events/:id", handlers.GetEvent)
	protected.POST("/tickets", handlers.CreateTicket)
	protected.GET("/tickets", handlers.GetTickets)
	protected.GET("/tickets/:id", handlers.GetUserTickets)
	protected.GET("/users/:id", handlers.GetUser)
	protected.GET("/users/:id/events", handlers.GetUserEvents)

	slog.Info("server started", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
