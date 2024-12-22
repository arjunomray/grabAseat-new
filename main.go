package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"arjunomray/eventer/dbConfig"
	_ "arjunomray/eventer/docs"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
