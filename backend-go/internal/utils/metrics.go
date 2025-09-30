package utils

import (
	"sync"
	"time"
)

type MetricType string

const (
	Counter   MetricType = "counter"
	Gauge     MetricType = "gauge"
	Histogram MetricType = "histogram"
	Summary   MetricType = "summary"
)

type Metric struct {
	Name   string
	Type   MetricType
	Value  float64
	Labels map[string]string
	Time   time.Time
}

type MetricsCollector struct {
	metrics map[string]*Metric
	mutex   sync.RWMutex
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*Metric),
	}
}

func (mc *MetricsCollector) IncrementCounter(name string, labels map[string]string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	key := mc.getKey(name, labels)
	if metric, exists := mc.metrics[key]; exists {
		metric.Value++
		metric.Time = time.Now()
	} else {
		mc.metrics[key] = &Metric{
			Name:   name,
			Type:   Counter,
			Value:  1,
			Labels: labels,
			Time:   time.Now(),
		}
	}
}

func (mc *MetricsCollector) SetGauge(name string, value float64, labels map[string]string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	key := mc.getKey(name, labels)
	mc.metrics[key] = &Metric{
		Name:   name,
		Type:   Gauge,
		Value:  value,
		Labels: labels,
		Time:   time.Now(),
	}
}

func (mc *MetricsCollector) ObserveHistogram(name string, value float64, labels map[string]string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	key := mc.getKey(name, labels)
	if metric, exists := mc.metrics[key]; exists {
		metric.Value = (metric.Value + value) / 2
		metric.Time = time.Now()
	} else {
		mc.metrics[key] = &Metric{
			Name:   name,
			Type:   Histogram,
			Value:  value,
			Labels: labels,
			Time:   time.Now(),
		}
	}
}

func (mc *MetricsCollector) GetMetric(name string, labels map[string]string) (*Metric, bool) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	key := mc.getKey(name, labels)
	metric, exists := mc.metrics[key]
	return metric, exists
}

func (mc *MetricsCollector) GetAllMetrics() map[string]*Metric {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	result := make(map[string]*Metric)
	for key, metric := range mc.metrics {
		result[key] = metric
	}
	return result
}

func (mc *MetricsCollector) Reset() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics = make(map[string]*Metric)
}

func (mc *MetricsCollector) getKey(name string, labels map[string]string) string {
	key := name
	for k, v := range labels {
		key += ":" + k + "=" + v
	}
	return key
}

type PerformanceMonitor struct {
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
}

func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		startTime: time.Now(),
	}
}

func (pm *PerformanceMonitor) Start() {
	pm.startTime = time.Now()
}

func (pm *PerformanceMonitor) Stop() {
	pm.endTime = time.Now()
	pm.duration = pm.endTime.Sub(pm.startTime)
}

func (pm *PerformanceMonitor) GetDuration() time.Duration {
	if pm.endTime.IsZero() {
		return time.Since(pm.startTime)
	}
	return pm.duration
}

func (pm *PerformanceMonitor) GetDurationMs() float64 {
	return float64(pm.GetDuration().Nanoseconds()) / 1e6
}

type HealthChecker struct {
	checks map[string]func() error
	mutex  sync.RWMutex
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]func() error),
	}
}

func (hc *HealthChecker) AddCheck(name string, check func() error) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	hc.checks[name] = check
}

func (hc *HealthChecker) RemoveCheck(name string) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	delete(hc.checks, name)
}

func (hc *HealthChecker) CheckAll() map[string]error {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()

	results := make(map[string]error)
	for name, check := range hc.checks {
		results[name] = check()
	}
	return results
}

func (hc *HealthChecker) IsHealthy() bool {
	results := hc.CheckAll()
	for _, err := range results {
		if err != nil {
			return false
		}
	}
	return true
}

type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
	mutex    sync.RWMutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, req := range requests {
			if req.After(cutoff) {
				validRequests = append(validRequests, req)
			}
		}
		rl.requests[key] = validRequests

		if len(validRequests) >= rl.limit {
			return false
		}
	}

	rl.requests[key] = append(rl.requests[key], now)
	return true
}

func (rl *RateLimiter) Reset(key string) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	delete(rl.requests, key)
}

func (rl *RateLimiter) ResetAll() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	rl.requests = make(map[string][]time.Time)
}

var globalMetrics = NewMetricsCollector()
var globalHealthChecker = NewHealthChecker()

func IncrementCounter(name string, labels map[string]string) {
	globalMetrics.IncrementCounter(name, labels)
}

func SetGauge(name string, value float64, labels map[string]string) {
	globalMetrics.SetGauge(name, value, labels)
}

func ObserveHistogram(name string, value float64, labels map[string]string) {
	globalMetrics.ObserveHistogram(name, value, labels)
}

func GetMetrics() map[string]*Metric {
	return globalMetrics.GetAllMetrics()
}

func AddHealthCheck(name string, check func() error) {
	globalHealthChecker.AddCheck(name, check)
}

func CheckHealth() map[string]error {
	return globalHealthChecker.CheckAll()
}

func IsHealthy() bool {
	return globalHealthChecker.IsHealthy()
}
