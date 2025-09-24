package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Metrics struct {
	RequestCount   int64
	ResponseTime   time.Duration
	ErrorCount     int64
	ActiveRequests int64
	mutex          sync.RWMutex
}

var GlobalMetrics = &Metrics{}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		GlobalMetrics.mutex.Lock()
		GlobalMetrics.ActiveRequests++
		GlobalMetrics.mutex.Unlock()

		c.Next()

		duration := time.Since(start)

		GlobalMetrics.mutex.Lock()
		GlobalMetrics.ActiveRequests--
		GlobalMetrics.RequestCount++
		GlobalMetrics.ResponseTime += duration

		if c.Writer.Status() >= 400 {
			GlobalMetrics.ErrorCount++
		}
		GlobalMetrics.mutex.Unlock()
	}
}

func (m *Metrics) GetStats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	avgResponseTime := time.Duration(0)
	if m.RequestCount > 0 {
		avgResponseTime = m.ResponseTime / time.Duration(m.RequestCount)
	}

	return map[string]interface{}{
		"request_count":     m.RequestCount,
		"active_requests":   m.ActiveRequests,
		"error_count":       m.ErrorCount,
		"avg_response_time": avgResponseTime.String(),
	}
}
