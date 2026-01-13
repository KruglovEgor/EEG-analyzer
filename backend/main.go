// Package main EEG Analyzer API
// @title EEG Analyzer API
// @version 1.0
// @description REST API for analyzing EEG (electroencephalography) data
// @description Supports single-file and multi-file analysis with various brain rhythm bands
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @host localhost:3000
// @BasePath /
package main

import (
	"log"
	"os"

	_ "eeg-analyzer/docs" // Import generated docs
	"eeg-analyzer/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Add memory monitoring middleware
	router.Use(MemoryMonitorMiddleware())

	// Enable gzip compression
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Configure CORS for frontend
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://vad1mchk.github.io"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Analysis endpoints
	router.POST("/analyze", handlers.AnalyzeMultipart) // Multipart/form-data (recommended)
	router.POST("/analyze-json", handlers.AnalyzeEEG)  // JSON with base64 (legacy/Swagger)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting EEG Analyzer API on port %s", port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
