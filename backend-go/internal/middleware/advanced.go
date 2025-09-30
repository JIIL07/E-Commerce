package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ecommerce-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logEntry := map[string]interface{}{
			"timestamp":  param.TimeStamp.Format(time.RFC3339),
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"method":     param.Method,
			"path":       param.Path,
			"user_agent": param.Request.UserAgent(),
			"error":      param.ErrorMessage,
		}

		jsonData, _ := json.Marshal(logEntry)
		return string(jsonData) + "\n"
	})
}

// MetricsMiddleware is defined in metrics.go

// RateLimitMiddleware is defined in rate_limiter.go

func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		cacheKey := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.String())

		if cached, exists := utils.GetCache("http").Get(cacheKey); exists {
			c.Header("X-Cache", "HIT")
			c.JSON(http.StatusOK, cached)
			c.Abort()
			return
		}

		c.Header("X-Cache", "MISS")
		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = utils.GenerateUUID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		c.Next()
	}
}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "Request timeout",
			})
			c.Abort()
		}
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		utils.Error("Panic recovered", "error", recovered, "path", c.Request.URL.Path)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
}

// CORSMiddleware is defined in cors.go

// AuthMiddleware is defined in auth.go

// AdminMiddleware is defined in auth.go

func ValidationMiddleware[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data T

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request data",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("validated_data", data)
		c.Next()
	}
}

func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := utils.GetIntQuery(c, "page", 1)
		limit := utils.GetIntQuery(c, "limit", 10)

		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 10
		}

		offset := (page - 1) * limit

		c.Set("page", page)
		c.Set("limit", limit)
		c.Set("offset", offset)

		c.Next()
	}
}

func CompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptEncoding := c.GetHeader("Accept-Encoding")

		if strings.Contains(acceptEncoding, "gzip") {
			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")
		}

		c.Next()
	}
}

func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			healthStatus := utils.CheckHealth()
			isHealthy := utils.IsHealthy()

			status := http.StatusOK
			if !isHealthy {
				status = http.StatusServiceUnavailable
			}

			c.JSON(status, gin.H{
				"status":    "ok",
				"timestamp": time.Now().Format(time.RFC3339),
				"healthy":   isHealthy,
				"checks":    healthStatus,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
