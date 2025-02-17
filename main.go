package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"arjunomray/eventer/dbConfig"
	"arjunomray/eventer/handlers"
	"arjunomray/eventer/middleware"
)

func main() {
	dbConfig.InitDatabase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Adjust this to match your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger setup

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

	log.Print("Server started on port 8080")
	r.Run()
}
