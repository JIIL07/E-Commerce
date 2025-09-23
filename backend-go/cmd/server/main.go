package main

import (
	"log"
	"os"

	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/handlers"
	"ecommerce-backend/internal/middleware"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if err := database.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())

	r.Use(logger.SetLogger())

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"timestamp": "2024-01-01T00:00:00Z",
			"uptime":    "0s",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "E-Commerce API Server",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": gin.H{
				"products":   "/api/products",
				"categories": "/api/categories",
				"auth":       "/api/auth",
				"cart":       "/api/cart",
				"orders":     "/api/orders",
				"reviews":    "/api/reviews",
				"payments":   "/api/payments",
				"wishlist":   "/api/wishlist",
				"health":     "/api/health",
			},
		})
	})

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), handlers.Profile)
	}

	products := r.Group("/api/products")
	{
		products.GET("/", handlers.GetProducts)
		products.GET("/featured", handlers.GetFeaturedProducts)
		products.GET("/:id", handlers.GetProduct)
	}

	categories := r.Group("/api/categories")
	{
		categories.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Categories endpoint - to be implemented"})
		})
	}

	cart := r.Group("/api/cart")
	cart.Use(middleware.AuthMiddleware())
	{
		cart.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Cart endpoint - to be implemented"})
		})
	}

	orders := r.Group("/api/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Orders endpoint - to be implemented"})
		})
	}

	reviews := r.Group("/api/reviews")
	{
		reviews.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Reviews endpoint - to be implemented"})
		})
	}

	payments := r.Group("/api/payments")
	payments.Use(middleware.AuthMiddleware())
	{
		payments.POST("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Payments endpoint - to be implemented"})
		})
	}

	wishlist := r.Group("/api/wishlist")
	wishlist.Use(middleware.AuthMiddleware())
	{
		wishlist.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Wishlist endpoint - to be implemented"})
		})
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Route not found",
			"path":    c.Request.URL.Path,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/api/health", port)
	log.Printf("🔐 Auth API: http://localhost:%s/api/auth", port)
	log.Printf("🛍️  Products API: http://localhost:%s/api/products", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
