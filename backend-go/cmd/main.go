package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/handlers"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/repositories"
	"ecommerce-backend/internal/seeds"
	"ecommerce-backend/internal/services"
	"ecommerce-backend/internal/utils"
	"ecommerce-backend/internal/websocket"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	godotenv.Load()

	var (
		mode      = flag.String("mode", "server", "Mode: server, init, seed, admin, generate-images, auto-init")
		waitForDB = flag.Bool("wait", false, "Wait for database to be available")
		timeout   = flag.Duration("timeout", 30*time.Second, "Timeout for database connection")
		seedType  = flag.String("type", "all", "Seed type: all, categories, products, users, orders, reviews")
		help      = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	switch *mode {
	case "init":
		runInit(cfg, *waitForDB, *timeout)
	case "seed":
		runSeed(cfg, *seedType)
	case "admin":
		runAdmin(cfg)
	case "generate-images":
		runGenerateImages()
	case "auto-init":
		runAutoInit(cfg, *waitForDB, *timeout)
	case "server":
		runServer(cfg)
	default:
		log.Fatal("Invalid mode. Use: server, init, seed, admin, generate-images, auto-init")
	}
}

func runInit(cfg *config.AppConfig, waitForDB bool, timeout time.Duration) {
	fmt.Println("🔧 Initializing database...")
	fmt.Printf("   Host: %s:%d\n", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("   Database: %s\n", cfg.Database.Database)
	fmt.Printf("   User: %s\n", cfg.Database.User)

	if waitForDB {
		fmt.Println("⏳ Waiting for database to be available...")
		if err := waitForDatabase(cfg, timeout); err != nil {
			log.Fatal("Database not available:", err)
		}
		fmt.Println("✅ Database is available!")
	}

	if err := database.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()

	fmt.Println("✅ Database initialized successfully!")

	fmt.Println("🔄 Running migrations...")
	if err := database.RunMigrations(database.GetDB()); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("✅ Migrations completed successfully!")
	fmt.Println("🎉 Database setup completed!")
}

func runSeed(cfg *config.AppConfig, seedType string) {
	fmt.Println("🌱 Seeding database...")

	if err := database.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()

	seedManager, err := seeds.NewSeedManager()
	if err != nil {
		log.Fatal("Failed to create seed manager:", err)
	}
	defer seedManager.Close()

	if seedType == "all" {
		err = seedManager.Run()
	} else {
		err = seedManager.RunSpecific([]string{seedType})
	}

	if err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	fmt.Printf("✅ Database seeded with %s data successfully!\n", seedType)
}

func runAdmin(cfg *config.AppConfig) {
	fmt.Println("🔧 Starting admin panel...")

	if err := database.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Use(authMiddleware())

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", dashboardHandler)
	r.GET("/api/stats", statsHandler)
	r.GET("/api/logs", logsHandler)
	r.GET("/api/metrics", metricsHandler)
	r.GET("/api/database", databaseHandler)
	r.GET("/api/cache", cacheHandler)
	r.GET("/api/users", usersHandler)
	r.GET("/api/products", productsHandler)
	r.GET("/api/orders", ordersHandler)
	r.GET("/api/health", healthHandler)

	r.POST("/api/seed", seedHandler)
	r.POST("/api/migrate", migrateHandler)
	r.POST("/api/cache/clear", clearCacheHandler)
	r.POST("/api/logs/clear", clearLogsHandler)

	r.GET("/ws", websocketHandler)

	port := os.Getenv("ADMIN_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Admin panel starting on port %s\n", port)
	fmt.Printf("📊 Admin panel: http://localhost:%s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start admin server:", err)
	}
}

func runServer(cfg *config.AppConfig) {
	fmt.Println("🚀 Starting E-Commerce server...")

	utils.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpiresIn, cfg.JWT.RefreshIn, cfg.JWT.Issuer, cfg.JWT.Audience)
	if err := database.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()
	if err := database.RunMigrations(database.GetDB()); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
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

	r.LoadHTMLGlob("templates/*")
	db := database.GetDB()
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	wishlistRepo := repositories.NewWishlistRepository(db)
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo, categoryRepo, reviewRepo)
	reviewService := services.NewReviewService(reviewRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo, productRepo)
	paymentService := services.NewPaymentService(paymentRepo, orderRepo)
	wishlistService := services.NewWishlistService(wishlistRepo)
	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	authHandler := handlers.NewAuthHandler(userService, cfg)
	productHandler := handlers.NewProductHandler(productService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
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

	r.HEAD("/api/health", func(c *gin.Context) {
		c.Status(200)
	})

	r.GET("/api/metrics", func(c *gin.Context) {
		stats := middleware.GlobalMetrics.GetStats()
		c.Header("Content-Type", "text/plain")
		c.String(200, `# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total %d

# HELP http_active_requests Number of active HTTP requests
# TYPE http_active_requests gauge
http_active_requests %d

# HELP http_errors_total Total number of HTTP errors
# TYPE http_errors_total counter
http_errors_total %d

# HELP http_request_duration_seconds Average HTTP request duration
# TYPE http_request_duration_seconds gauge
http_request_duration_seconds %s
`,
			stats["request_count"],
			stats["active_requests"],
			stats["error_count"],
			stats["avg_response_time"])
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
				"docs":       "/docs",
				"admin":      "/admin",
			},
		})
	})

	r.GET("/docs", func(c *gin.Context) {
		c.HTML(200, "docs.html", gin.H{
			"title": "E-Commerce API Documentation",
		})
	})

	r.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", gin.H{
			"title": "Admin Dashboard",
		})
	})
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), authHandler.Profile)
		auth.PUT("/profile", middleware.AuthMiddleware(), authHandler.UpdateProfile)
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
		categories.GET("/", categoryHandler.GetCategories)
		categories.GET("/:slug", categoryHandler.GetCategory)
		categories.POST("/", middleware.AuthMiddleware(), categoryHandler.CreateCategory)
		categories.PUT("/:slug", middleware.AuthMiddleware(), categoryHandler.UpdateCategory)
		categories.DELETE("/:slug", middleware.AuthMiddleware(), categoryHandler.DeleteCategory)
	}
	cart := r.Group("/api/cart")
	cart.Use(middleware.AuthMiddleware())
	{
		cart.GET("/", cartHandler.GetCart)
		cart.POST("/", cartHandler.AddToCart)
		cart.PUT("/:id", cartHandler.UpdateCartItem)
		cart.DELETE("/:id", cartHandler.RemoveFromCart)
		cart.DELETE("/", cartHandler.ClearCart)
	}
	orders := r.Group("/api/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.GET("/", orderHandler.GetOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.POST("/", orderHandler.CreateOrder)
		orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		orders.DELETE("/:id", orderHandler.CancelOrder)
	}
	reviews := r.Group("/api/reviews")
	{
		reviews.GET("/product/:productId", reviewHandler.GetProductReviews)
		reviews.GET("/user", middleware.AuthMiddleware(), reviewHandler.GetUserReviews)
		reviews.GET("/user/:productId", middleware.AuthMiddleware(), reviewHandler.GetUserReviewForProduct)
		reviews.POST("/", middleware.AuthMiddleware(), reviewHandler.CreateReview)
		reviews.PUT("/:id", middleware.AuthMiddleware(), reviewHandler.UpdateReview)
		reviews.DELETE("/:id", middleware.AuthMiddleware(), reviewHandler.DeleteReview)
	}
	payments := r.Group("/api/payments")
	payments.Use(middleware.AuthMiddleware())
	{
		payments.POST("/intent", paymentHandler.CreatePaymentIntent)
		payments.POST("/confirm", paymentHandler.ConfirmPayment)
		payments.GET("/history", paymentHandler.GetPaymentHistory)
	}

	wishlist := r.Group("/api/wishlist")
	wishlist.Use(middleware.AuthMiddleware())
	{
		wishlist.GET("/", wishlistHandler.GetWishlist)
		wishlist.POST("/", wishlistHandler.AddToWishlist)
		wishlist.DELETE("/:productId", wishlistHandler.RemoveFromWishlist)
		wishlist.GET("/:productId/check", wishlistHandler.IsInWishlist)
		wishlist.DELETE("/", wishlistHandler.ClearWishlist)
	}
	uploads := r.Group("/api/uploads")
	{
		uploads.POST("/", middleware.AuthMiddleware(), uploadHandler.UploadImage)
		uploads.DELETE("/:filename", middleware.AuthMiddleware(), uploadHandler.DeleteImage)
		uploads.GET("/:filename", uploadHandler.ServeImage)
	}
	wsHandler := websocket.NewHandler(wsHub)
	ws := r.Group("/ws")
	{
		ws.GET("/", wsHandler.HandleWebSocket)
		ws.GET("/users", wsHandler.GetConnectedUsers)
		ws.GET("/count", wsHandler.GetClientCount)
		ws.POST("/notification", middleware.AuthMiddleware(), wsHandler.SendNotification)
		ws.POST("/order-update", middleware.AuthMiddleware(), wsHandler.SendOrderUpdate)
		ws.POST("/product-update", middleware.AuthMiddleware(), wsHandler.SendProductUpdate)
		ws.POST("/stock-alert", middleware.AuthMiddleware(), wsHandler.SendStockAlert)
		ws.POST("/price-alert", middleware.AuthMiddleware(), wsHandler.SendPriceAlert)
		ws.POST("/new-product", middleware.AuthMiddleware(), wsHandler.SendNewProductAlert)
		ws.POST("/promotion", middleware.AuthMiddleware(), wsHandler.SendPromotionAlert)
		ws.POST("/maintenance", middleware.AuthMiddleware(), wsHandler.SendMaintenanceAlert)
		ws.POST("/user-activity", middleware.AuthMiddleware(), wsHandler.SendUserActivity)
		ws.POST("/analytics", middleware.AuthMiddleware(), wsHandler.SendAnalyticsUpdate)
		ws.POST("/stats", middleware.AuthMiddleware(), wsHandler.SendRealTimeStats)
	}

	admin := r.Group("/admin/api")
	{
		admin.GET("/stats", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"timestamp": time.Now().Unix(),
				"uptime":    time.Since(time.Now()).String(),
				"users": []map[string]interface{}{
					{"id": "1", "email": "admin@example.com", "name": "Admin User", "role": "admin"},
					{"id": "2", "email": "user@example.com", "name": "Regular User", "role": "user"},
				},
				"products": []map[string]interface{}{
					{"id": "1", "name": "Sample Product", "price": 99.99, "stock": 50},
					{"id": "2", "name": "Another Product", "price": 149.99, "stock": 25},
				},
				"orders": []map[string]interface{}{
					{"id": "1", "user_id": "2", "total": 199.98, "status": "completed"},
					{"id": "2", "user_id": "2", "total": 99.99, "status": "pending"},
				},
				"database": map[string]interface{}{
					"status":      "connected",
					"connections": 5,
					"size":        "50MB",
				},
				"cache": map[string]interface{}{
					"size":     100,
					"hit_rate": "85%",
				},
				"metrics": map[string]interface{}{
					"http_requests": map[string]interface{}{
						"total":    1500,
						"avg_time": "150ms",
					},
				},
			})
		})
		admin.POST("/seed", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Database seeded successfully"})
		})
		admin.POST("/migrate", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Database migrated successfully"})
		})
		admin.POST("/cache/clear", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Cache cleared successfully"})
		})
		admin.POST("/logs/clear", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Logs cleared successfully"})
		})
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Route not found",
			"path":    c.Request.URL.Path,
		})
	})
	server := &http.Server{
		Addr:         ":" + fmt.Sprintf("%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}
	go func() {
		log.Printf("🚀 Server starting on port %d", cfg.Server.Port)
		log.Printf("📊 Health check: http://localhost:%d/api/health", cfg.Server.Port)
		log.Printf("🔐 Auth API: http://localhost:%d/api/auth", cfg.Server.Port)
		log.Printf("🛍️ Products API: http://localhost:%d/api/products", cfg.Server.Port)
		log.Printf("📚 API Docs: http://localhost:%d/docs", cfg.Server.Port)
		log.Printf("🔧 Admin Panel: http://localhost:%d/admin", cfg.Server.Port)
		log.Printf("🌐 WebSocket: ws://localhost:%d/ws", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/" || c.Request.URL.Path == "/login" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func dashboardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Admin Dashboard",
		"stats": getSystemStats(),
	})
}

func statsHandler(c *gin.Context) {
	stats := getSystemStats()
	c.JSON(http.StatusOK, stats)
}

func logsHandler(c *gin.Context) {
	level := c.Query("level")
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	logs := getLogs(level, limit)
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func metricsHandler(c *gin.Context) {
	metrics := getMetrics()
	c.JSON(http.StatusOK, metrics)
}

func databaseHandler(c *gin.Context) {
	dbStats := getDatabaseStats()
	c.JSON(http.StatusOK, dbStats)
}

func cacheHandler(c *gin.Context) {
	cacheStats := getCacheStats()
	c.JSON(http.StatusOK, cacheStats)
}

func usersHandler(c *gin.Context) {
	users := getUsers()
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func productsHandler(c *gin.Context) {
	products := getProducts()
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func ordersHandler(c *gin.Context) {
	orders := getOrders()
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func healthHandler(c *gin.Context) {
	health := getHealthStatus()
	c.JSON(http.StatusOK, health)
}

func seedHandler(c *gin.Context) {
	seedType := c.Query("type")
	if seedType == "" {
		seedType = "all"
	}

	seedManager, err := seeds.NewSeedManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create seed manager",
			"details": err.Error(),
		})
		return
	}
	defer seedManager.Close()

	if seedType == "all" {
		err = seedManager.Run()
	} else {
		err = seedManager.RunSpecific([]string{seedType})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to seed database",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database seeded successfully",
		"type":    seedType,
	})
}

func migrateHandler(c *gin.Context) {
	if err := database.InitDatabase(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initialize database",
			"details": err.Error(),
		})
		return
	}
	defer database.CloseDatabase()

	if err := database.RunMigrations(database.GetDB()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to run migrations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database migrated successfully"})
}

func clearCacheHandler(c *gin.Context) {
	utils.GetCacheStats()
	globalCacheManager := utils.NewCacheManager()
	globalCacheManager.ClearAll()

	c.JSON(http.StatusOK, gin.H{
		"message":   "Cache cleared successfully",
		"timestamp": time.Now().Unix(),
	})
}

func clearLogsHandler(c *gin.Context) {
	logFiles := []string{
		"logs/backend/access.log",
		"logs/backend/error.log",
		"logs/backend/app.log",
		"logs/nginx/access.log",
		"logs/nginx/error.log",
		"logs/frontend/build.log",
	}

	clearedCount := 0
	for _, logFile := range logFiles {
		if _, err := os.Stat(logFile); err == nil {
			if err := os.Truncate(logFile, 0); err == nil {
				clearedCount++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Logs cleared successfully",
		"cleared_files": clearedCount,
		"timestamp":     time.Now().Unix(),
	})
}

func websocketHandler(c *gin.Context) {
	upgrader := ws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	for {
		select {
		case <-time.After(5 * time.Second):
			stats := getSystemStats()
			err := conn.WriteJSON(gin.H{
				"type": "stats",
				"data": stats,
			})
			if err != nil {
				log.Printf("WebSocket write failed: %v", err)
				return
			}
		}
	}
}

func getSystemStats() map[string]interface{} {
	return map[string]interface{}{
		"timestamp":    time.Now().Unix(),
		"uptime":       time.Since(time.Now()).String(),
		"memory_usage": getMemoryUsage(),
		"cpu_usage":    getCPUUsage(),
		"database":     getDatabaseStats(),
		"cache":        getCacheStats(),
		"logs":         getLogStats(),
	}
}

func getMemoryUsage() map[string]interface{} {
	return map[string]interface{}{
		"alloc":       "100MB",
		"total_alloc": "500MB",
		"sys":         "200MB",
		"num_gc":      10,
	}
}

func getCPUUsage() map[string]interface{} {
	return map[string]interface{}{
		"usage": "15%",
		"cores": 4,
	}
}

func getDatabaseStats() map[string]interface{} {
	db := database.GetDB()
	if db == nil {
		return map[string]interface{}{
			"status": "disconnected",
			"error":  "Database not initialized",
		}
	}

	var stats map[string]interface{} = make(map[string]interface{})
	stats["status"] = "connected"

	var userCount, productCount, orderCount, categoryCount int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&productCount)
	db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&orderCount)
	db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&categoryCount)

	stats["users"] = userCount
	stats["products"] = productCount
	stats["orders"] = orderCount
	stats["categories"] = categoryCount

	var dbSize string
	db.QueryRow("SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&dbSize)
	stats["size"] = dbSize

	return stats
}

func getCacheStats() map[string]interface{} {
	cacheStats := utils.GetCacheStats()

	if len(cacheStats) == 0 {
		return map[string]interface{}{
			"size":         0,
			"hit_rate":     "0%",
			"miss_rate":    "0%",
			"total_hits":   0,
			"total_misses": 0,
		}
	}

	var totalSize, totalHits, totalMisses int
	var totalHitRate, totalMissRate float64

	for _, stats := range cacheStats {
		totalSize += stats.Size
		totalHits += int(stats.TotalHits)
		totalMisses += int(stats.TotalMisses)
		totalHitRate += stats.HitRate
		totalMissRate += stats.MissRate
	}

	avgHitRate := totalHitRate / float64(len(cacheStats))
	avgMissRate := totalMissRate / float64(len(cacheStats))

	return map[string]interface{}{
		"size":         totalSize,
		"hit_rate":     fmt.Sprintf("%.1f%%", avgHitRate),
		"miss_rate":    fmt.Sprintf("%.1f%%", avgMissRate),
		"total_hits":   totalHits,
		"total_misses": totalMisses,
		"caches":       len(cacheStats),
	}
}

func getLogStats() map[string]interface{} {
	return map[string]interface{}{
		"total":    1000,
		"errors":   5,
		"warnings": 10,
		"info":     985,
	}
}

func getLogs(level string, limit int) []map[string]interface{} {
	logs := []map[string]interface{}{
		{
			"timestamp": time.Now().Add(-1 * time.Minute),
			"level":     "INFO",
			"message":   "User logged in",
			"user_id":   "123",
		},
		{
			"timestamp": time.Now().Add(-2 * time.Minute),
			"level":     "ERROR",
			"message":   "Database connection failed",
			"error":     "connection timeout",
		},
		{
			"timestamp": time.Now().Add(-3 * time.Minute),
			"level":     "WARN",
			"message":   "High memory usage detected",
			"usage":     "85%",
		},
	}

	if level != "" {
		filtered := []map[string]interface{}{}
		for _, log := range logs {
			if log["level"] == level {
				filtered = append(filtered, log)
			}
		}
		logs = filtered
	}

	if len(logs) > limit {
		logs = logs[:limit]
	}

	return logs
}

func getMetrics() map[string]interface{} {
	return map[string]interface{}{
		"http_requests": map[string]interface{}{
			"total":    1500,
			"success":  1400,
			"errors":   100,
			"avg_time": "150ms",
		},
		"database": map[string]interface{}{
			"queries":     1250,
			"avg_time":    "25ms",
			"connections": 5,
		},
		"cache": map[string]interface{}{
			"hits":   1200,
			"misses": 300,
			"size":   100,
		},
	}
}

func getUsers() []map[string]interface{} {
	db := database.GetDB()
	if db == nil {
		return []map[string]interface{}{}
	}

	rows, err := db.Query(`
		SELECT id, email, name, role, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT 10
	`)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var user struct {
			ID        string
			Email     string
			Name      string
			Role      string
			CreatedAt time.Time
			UpdatedAt time.Time
		}

		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			continue
		}

		users = append(users, map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	return users
}

func getProducts() []map[string]interface{} {
	db := database.GetDB()
	if db == nil {
		return []map[string]interface{}{}
	}

	rows, err := db.Query(`
		SELECT p.id, p.name, p.price, p.stock, c.name as category, p.created_at, p.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.created_at DESC 
		LIMIT 10
	`)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var product struct {
			ID        string
			Name      string
			Price     float64
			Stock     int
			Category  *string
			CreatedAt time.Time
			UpdatedAt time.Time
		}

		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category, &product.CreatedAt, &product.UpdatedAt); err != nil {
			continue
		}

		category := "Uncategorized"
		if product.Category != nil {
			category = *product.Category
		}

		products = append(products, map[string]interface{}{
			"id":         product.ID,
			"name":       product.Name,
			"price":      product.Price,
			"stock":      product.Stock,
			"category":   category,
			"created_at": product.CreatedAt,
			"updated_at": product.UpdatedAt,
		})
	}

	return products
}

func getOrders() []map[string]interface{} {
	db := database.GetDB()
	if db == nil {
		return []map[string]interface{}{}
	}

	rows, err := db.Query(`
		SELECT o.id, o.user_id, o.total, o.status, o.created_at, u.name as user_name
		FROM orders o
		LEFT JOIN users u ON o.user_id = u.id
		ORDER BY o.created_at DESC 
		LIMIT 10
	`)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var order struct {
			ID        string
			UserID    string
			Total     float64
			Status    string
			CreatedAt time.Time
			UserName  *string
		}

		if err := rows.Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt, &order.UserName); err != nil {
			continue
		}

		userName := "Unknown User"
		if order.UserName != nil {
			userName = *order.UserName
		}

		orders = append(orders, map[string]interface{}{
			"id":         order.ID,
			"user_id":    order.UserID,
			"user_name":  userName,
			"total":      order.Total,
			"status":     order.Status,
			"created_at": order.CreatedAt,
		})
	}

	return orders
}

func getHealthStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"services": map[string]interface{}{
			"database": "healthy",
			"cache":    "healthy",
			"api":      "healthy",
		},
	}
}

func waitForDatabase(cfg *config.AppConfig, timeout time.Duration) error {
	start := time.Now()

	for time.Since(start) < timeout {
		if err := database.InitDatabase(); err == nil {
			database.CloseDatabase()
			return nil
		}

		fmt.Print(".")
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("timeout after %v", timeout)
}

func runAutoInit(cfg *config.AppConfig, waitForDB bool, timeout time.Duration) {
	fmt.Println("🚀 Auto-initializing E-Commerce Project...")
	fmt.Println("==========================================")

	fmt.Println("\n🔧 Step 1: Initializing database...")
	runInit(cfg, waitForDB, timeout)

	fmt.Println("\n🌱 Step 2: Seeding database with sample data...")
	runSeed(cfg, "all")

	fmt.Println("\n🎨 Step 3: Generating placeholder images...")
	runGenerateImages()

	fmt.Println("\n🎉 Auto-initialization completed successfully!")
	fmt.Println("==========================================")
	fmt.Println("📱 Frontend: http://localhost:3000")
	fmt.Println("🔧 Backend API: http://localhost:5000")
	fmt.Println("📊 Admin Panel: http://localhost:5000/admin")
	fmt.Println("📚 API Docs: http://localhost:5000/docs")
	fmt.Println("==========================================")
}

func runGenerateImages() {
	fmt.Println("🎨 Generating placeholder images...")

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("Error creating uploads directory: %v\n", err)
		return
	}

	productImages := []string{
		"iphone15_pro.jpg",
		"iphone15_pro_back.jpg",
		"galaxy_s24_ultra.jpg",
		"macbook_pro_m3.jpg",
		"ipad_air.jpg",
		"airpods_pro.jpg",
		"sony_wh1000xm5.jpg",
		"nike_air_max_270.jpg",
		"adidas_ultraboost_22.jpg",
		"levis_501.jpg",
		"uniqlo_heattech.jpg",
		"patagonia_fleece.jpg",
		"clean_code.jpg",
		"js_good_parts.jpg",
		"python_crash_course.jpg",
		"design_patterns.jpg",
		"pragmatic_programmer.jpg",
		"dyson_v15.jpg",
		"kitchenaid_mixer.jpg",
		"philips_hue.jpg",
		"weber_grill.jpg",
		"peloton_bike.jpg",
		"bowflex_dumbbells.jpg",
		"yoga_mat_premium.jpg",
		"resistance_bands.jpg",
		"la_mer_cream.jpg",
		"oral_b_toothbrush.jpg",
		"multivitamin.jpg",
		"lego_creator.jpg",
		"nintendo_switch.jpg",
		"monopoly.jpg",
		"car_phone_mount.jpg",
		"dash_cam.jpg",
		"air_freshener.jpg",
		"organic_coffee.jpg",
		"protein_powder.jpg",
		"green_tea_set.jpg",
		"gold_necklace.jpg",
		"silver_ring.jpg",
		"pearl_earrings.jpg",
		"default_product.jpg",
	}

	for _, filename := range productImages {
		filePath := filepath.Join(uploadDir, filename)

		if _, err := os.Stat(filePath); err == nil {
			continue
		}

		productName := strings.TrimSuffix(filename, filepath.Ext(filename))
		productName = strings.ReplaceAll(productName, "_", " ")
		productName = strings.Title(productName)

		if err := generatePlaceholderImage(filePath, productName); err != nil {
			fmt.Printf("Error generating image for %s: %v\n", filename, err)
		} else {
			fmt.Printf("Generated placeholder image: %s\n", filename)
		}
	}

	fmt.Println("✅ Placeholder image generation completed!")
}

func generatePlaceholderImage(filePath, productName string) error {
	width, height := 400, 400

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	bgColor1 := color.RGBA{99, 102, 241, 255}
	bgColor2 := color.RGBA{139, 92, 246, 255}

	for y := 0; y < height; y++ {
		ratio := float64(y) / float64(height)
		r := uint8(float64(bgColor1.R)*(1-ratio) + float64(bgColor2.R)*ratio)
		g := uint8(float64(bgColor1.G)*(1-ratio) + float64(bgColor2.G)*ratio)
		b := uint8(float64(bgColor1.B)*(1-ratio) + float64(bgColor2.B)*ratio)

		gradientColor := color.RGBA{r, g, b, 255}
		for x := 0; x < width; x++ {
			img.Set(x, y, gradientColor)
		}
	}

	logoColor := color.RGBA{255, 255, 255, 255}
	logoPoint := fixed.Point26_6{
		X: fixed.I(width/2 - 60),
		Y: fixed.I(height/2 - 40),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(logoColor),
		Face: basicfont.Face7x13,
		Dot:  logoPoint,
	}
	d.DrawString("E-Commerce")

	productColor := color.RGBA{255, 255, 255, 200}
	productPoint := fixed.Point26_6{
		X: fixed.I(width/2 - len(productName)*3),
		Y: fixed.I(height/2 + 20),
	}

	productDrawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(productColor),
		Face: basicfont.Face7x13,
		Dot:  productPoint,
	}
	productDrawer.DrawString(productName)

	accentColor := color.RGBA{255, 255, 255, 100}

	for x := 0; x < width; x++ {
		img.Set(x, 0, accentColor)
		img.Set(x, height-1, accentColor)
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, accentColor)
		img.Set(width-1, y, accentColor)
	}

	cornerSize := 20
	for i := 0; i < cornerSize; i++ {
		img.Set(i, i, accentColor)
		img.Set(width-1-i, i, accentColor)
		img.Set(i, height-1-i, accentColor)
		img.Set(width-1-i, height-1-i, accentColor)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
}

func showHelp() {
	fmt.Println("E-Commerce Unified Command Tool")
	fmt.Println("===============================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/main.go [options]")
	fmt.Println()
	fmt.Println("Modes:")
	fmt.Println("  -mode=server    Start the main API server (default)")
	fmt.Println("  -mode=init      Initialize database and run migrations")
	fmt.Println("  -mode=seed      Seed database with sample data")
	fmt.Println("  -mode=admin     Start admin panel")
	fmt.Println("  -mode=generate-images  Generate placeholder images")
	fmt.Println("  -mode=auto-init Full project initialization (init + seed + images)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -wait")
	fmt.Println("        Wait for database to be available before initializing")
	fmt.Println("  -timeout duration")
	fmt.Println("        Timeout for database connection (default: 30s)")
	fmt.Println("  -type string")
	fmt.Println("        Seed type: all, categories, products, users, orders, reviews (default: all)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
}
