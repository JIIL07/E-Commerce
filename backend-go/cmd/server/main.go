package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/handlers"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/repositories"
	"ecommerce-backend/internal/services"
	"ecommerce-backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.Load()

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
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.MetricsMiddleware())
	r.Use(middleware.RateLimitMiddleware(100, time.Minute))

	db := database.GetDB()

	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	cartRepo := repositories.NewCartRepository(db)

	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo, categoryRepo, reviewRepo)
	reviewService := services.NewReviewService(reviewRepo)
	cartService := services.NewCartService(cartRepo, productRepo)

	authHandler := handlers.NewAuthHandler(userService, cfg)
	productHandler := handlers.NewProductHandler(productService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	cartHandler := handlers.NewCartHandler(cartService)
	uploadHandler := handlers.NewUploadHandler("./uploads")

	wsHub := websocket.NewHub()
	go wsHub.Run()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    time.Since(time.Now()).String(),
			"metrics":   middleware.GlobalMetrics.GetStats(),
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "E-Commerce API Server",
			"version": "2.0.0",
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
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/profile", middleware.AuthMiddleware(cfg), authHandler.Profile)
		auth.PUT("/profile", middleware.AuthMiddleware(cfg), authHandler.UpdateProfile)
	}

	products := r.Group("/api/products")
	{
		products.GET("/", productHandler.GetProducts)
		products.GET("/featured", productHandler.GetFeaturedProducts)
		products.GET("/search", productHandler.SearchProducts)
		products.GET("/:id", productHandler.GetProduct)
	}

	categories := r.Group("/api/categories")
	{
		categories.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Categories endpoint - to be implemented"})
		})
	}

	cart := r.Group("/api/cart")
	cart.Use(middleware.AuthMiddleware(cfg))
	{
		cart.GET("/", cartHandler.GetCart)
		cart.POST("/", cartHandler.AddToCart)
		cart.PUT("/:id", cartHandler.UpdateCartItem)
		cart.DELETE("/:id", cartHandler.RemoveFromCart)
		cart.DELETE("/", cartHandler.ClearCart)
	}

	orders := r.Group("/api/orders")
	orders.Use(middleware.AuthMiddleware(cfg))
	{
		orders.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Orders endpoint - to be implemented"})
		})
	}

	reviews := r.Group("/api/reviews")
	{
		reviews.GET("/product/:productId", reviewHandler.GetProductReviews)
		reviews.GET("/user", middleware.AuthMiddleware(cfg), reviewHandler.GetUserReviews)
		reviews.GET("/user/:productId", middleware.AuthMiddleware(cfg), reviewHandler.GetUserReviewForProduct)
		reviews.POST("/", middleware.AuthMiddleware(cfg), reviewHandler.CreateReview)
		reviews.PUT("/:id", middleware.AuthMiddleware(cfg), reviewHandler.UpdateReview)
		reviews.DELETE("/:id", middleware.AuthMiddleware(cfg), reviewHandler.DeleteReview)
	}

	payments := r.Group("/api/payments")
	payments.Use(middleware.AuthMiddleware(cfg))
	{
		payments.POST("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Payments endpoint - to be implemented"})
		})
	}

	wishlist := r.Group("/api/wishlist")
	wishlist.Use(middleware.AuthMiddleware(cfg))
	{
		wishlist.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Wishlist endpoint - to be implemented"})
		})
	}

	uploads := r.Group("/api/uploads")
	{
		uploads.POST("/", middleware.AuthMiddleware(cfg), uploadHandler.UploadImage)
		uploads.DELETE("/:filename", middleware.AuthMiddleware(cfg), uploadHandler.DeleteImage)
		uploads.GET("/:filename", uploadHandler.ServeImage)
	}

	ws := r.Group("/ws")
	ws.Use(middleware.AuthMiddleware(cfg))
	{
		ws.GET("/", websocket.HandleWebSocket(wsHub))
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Route not found",
			"path":    c.Request.URL.Path,
		})
	})

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
		log.Printf("📊 Health check: http://localhost:%s/api/health", cfg.Server.Port)
		log.Printf("🔐 Auth API: http://localhost:%s/api/auth", cfg.Server.Port)
		log.Printf("🛍️  Products API: http://localhost:%s/api/products", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
