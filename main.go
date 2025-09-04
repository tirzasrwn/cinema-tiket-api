package main

import (
	"cinema-ticket-api/config"
	"cinema-ticket-api/handlers"
	"cinema-ticket-api/middleware"
	"log"
	"os"

	_ "cinema-ticket-api/docs" // Import generated docs

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Cinema Ticket API
//	@version		1.0
//	@description	API untuk sistem pembelian tiket bioskop online
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.email	support@cinema.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:4000
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	config.InitDB()
	defer config.DB.Close()

	// Initialize router
	router := gin.Default()

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/login", handlers.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Screening CRUD routes
		protected.POST("/screenings", handlers.CreateScreening)
		protected.GET("/screenings", handlers.GetScreenings)
		protected.GET("/screenings/:id", handlers.GetScreening)
		protected.PUT("/screenings/:id", handlers.UpdateScreening)
		protected.DELETE("/screenings/:id", handlers.DeleteScreening)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
