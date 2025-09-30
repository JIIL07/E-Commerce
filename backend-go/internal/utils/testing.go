package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestHelper struct {
	Router *gin.Engine
	Client *http.Client
}

func NewTestHelper() *TestHelper {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &TestHelper{
		Router: router,
		Client: client,
	}
}

func (th *TestHelper) AddMiddleware(middleware gin.HandlerFunc) {
	th.Router.Use(middleware)
}

func (th *TestHelper) AddRoute(method, path string, handler gin.HandlerFunc) {
	th.Router.Handle(method, path, handler)
}

func (th *TestHelper) GET(path string, handler gin.HandlerFunc) {
	th.Router.GET(path, handler)
}

func (th *TestHelper) POST(path string, handler gin.HandlerFunc) {
	th.Router.POST(path, handler)
}

func (th *TestHelper) PUT(path string, handler gin.HandlerFunc) {
	th.Router.PUT(path, handler)
}

func (th *TestHelper) DELETE(path string, handler gin.HandlerFunc) {
	th.Router.DELETE(path, handler)
}

func (th *TestHelper) MakeRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody io.Reader

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	th.Router.ServeHTTP(w, req)

	return w
}

func (th *TestHelper) AssertResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	assert.Equal(t, expectedStatus, w.Code)

	if expectedBody != nil {
		var expectedJSON []byte
		if str, ok := expectedBody.(string); ok {
			expectedJSON = []byte(str)
		} else {
			expectedJSON, _ = json.Marshal(expectedBody)
		}

		assert.JSONEq(t, string(expectedJSON), w.Body.String())
	}
}

func (th *TestHelper) AssertResponseContains(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedText string) {
	assert.Equal(t, expectedStatus, w.Code)
	assert.Contains(t, w.Body.String(), expectedText)
}

func (th *TestHelper) AssertResponseHeader(t *testing.T, w *httptest.ResponseRecorder, headerName, expectedValue string) {
	assert.Equal(t, expectedValue, w.Header().Get(headerName))
}

func (th *TestHelper) ParseResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), target)
	require.NoError(t, err)
}

type MockDatabase struct {
	Data map[string]interface{}
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		Data: make(map[string]interface{}),
	}
}

func (mdb *MockDatabase) Set(key string, value interface{}) {
	mdb.Data[key] = value
}

func (mdb *MockDatabase) Get(key string) (interface{}, bool) {
	value, exists := mdb.Data[key]
	return value, exists
}

func (mdb *MockDatabase) Delete(key string) {
	delete(mdb.Data, key)
}

func (mdb *MockDatabase) Clear() {
	mdb.Data = make(map[string]interface{})
}

func (mdb *MockDatabase) Size() int {
	return len(mdb.Data)
}

type MockCache struct {
	Data map[string]interface{}
	TTL  map[string]time.Time
}

func NewMockCache() *MockCache {
	return &MockCache{
		Data: make(map[string]interface{}),
		TTL:  make(map[string]time.Time),
	}
}

func (mc *MockCache) Set(key string, value interface{}, ttl time.Duration) {
	mc.Data[key] = value
	if ttl > 0 {
		mc.TTL[key] = time.Now().Add(ttl)
	}
}

func (mc *MockCache) Get(key string) (interface{}, bool) {
	value, exists := mc.Data[key]
	if !exists {
		return nil, false
	}

	if ttl, hasTTL := mc.TTL[key]; hasTTL && time.Now().After(ttl) {
		delete(mc.Data, key)
		delete(mc.TTL, key)
		return nil, false
	}

	return value, true
}

func (mc *MockCache) Delete(key string) {
	delete(mc.Data, key)
	delete(mc.TTL, key)
}

func (mc *MockCache) Clear() {
	mc.Data = make(map[string]interface{})
	mc.TTL = make(map[string]time.Time)
}

func (mc *MockCache) Size() int {
	return len(mc.Data)
}

type MockLogger struct {
	Logs []LogEntry
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		Logs: make([]LogEntry, 0),
	}
}

func (ml *MockLogger) Debug(message string, data ...any) {
	ml.log(DEBUG, message, data...)
}

func (ml *MockLogger) Info(message string, data ...any) {
	ml.log(INFO, message, data...)
}

func (ml *MockLogger) Warn(message string, data ...any) {
	ml.log(WARN, message, data...)
}

func (ml *MockLogger) Error(message string, data ...any) {
	ml.log(ERROR, message, data...)
}

func (ml *MockLogger) Fatal(message string, data ...any) {
	ml.log(FATAL, message, data...)
}

func (ml *MockLogger) log(level LogLevel, message string, data ...any) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     getLevelString(level),
		Message:   message,
	}

	if len(data) > 0 {
		entry.Data = data
	}

	ml.Logs = append(ml.Logs, entry)
}

func (ml *MockLogger) GetLogs() []LogEntry {
	return ml.Logs
}

func (ml *MockLogger) GetLogsByLevel(level LogLevel) []LogEntry {
	var filtered []LogEntry
	levelStr := getLevelString(level)

	for _, log := range ml.Logs {
		if log.Level == levelStr {
			filtered = append(filtered, log)
		}
	}

	return filtered
}

func (ml *MockLogger) Clear() {
	ml.Logs = make([]LogEntry, 0)
}

func (ml *MockLogger) Count() int {
	return len(ml.Logs)
}

func (ml *MockLogger) HasLog(level LogLevel, message string) bool {
	levelStr := getLevelString(level)

	for _, log := range ml.Logs {
		if log.Level == levelStr && log.Message == message {
			return true
		}
	}

	return false
}

func getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type TestData struct {
	Users    []map[string]interface{}
	Products []map[string]interface{}
	Orders   []map[string]interface{}
}

func NewTestData() *TestData {
	return &TestData{
		Users: []map[string]interface{}{
			{
				"id":       "1",
				"email":    "test@example.com",
				"name":     "Test User",
				"role":     "user",
				"password": "password123",
			},
			{
				"id":       "2",
				"email":    "admin@example.com",
				"name":     "Admin User",
				"role":     "admin",
				"password": "admin123",
			},
		},
		Products: []map[string]interface{}{
			{
				"id":          "1",
				"name":        "Test Product",
				"description": "A test product",
				"price":       99.99,
				"stock":       10,
				"category_id": "1",
			},
			{
				"id":          "2",
				"name":        "Another Product",
				"description": "Another test product",
				"price":       199.99,
				"stock":       5,
				"category_id": "2",
			},
		},
		Orders: []map[string]interface{}{
			{
				"id":         "1",
				"user_id":    "1",
				"total":      99.99,
				"status":     "pending",
				"created_at": time.Now(),
			},
			{
				"id":         "2",
				"user_id":    "2",
				"total":      199.99,
				"status":     "completed",
				"created_at": time.Now().Add(-24 * time.Hour),
			},
		},
	}
}

func (td *TestData) GetUser(id string) (map[string]interface{}, bool) {
	for _, user := range td.Users {
		if user["id"] == id {
			return user, true
		}
	}
	return nil, false
}

func (td *TestData) GetProduct(id string) (map[string]interface{}, bool) {
	for _, product := range td.Products {
		if product["id"] == id {
			return product, true
		}
	}
	return nil, false
}

func (td *TestData) GetOrder(id string) (map[string]interface{}, bool) {
	for _, order := range td.Orders {
		if order["id"] == id {
			return order, true
		}
	}
	return nil, false
}

func (td *TestData) AddUser(user map[string]interface{}) {
	td.Users = append(td.Users, user)
}

func (td *TestData) AddProduct(product map[string]interface{}) {
	td.Products = append(td.Products, product)
}

func (td *TestData) AddOrder(order map[string]interface{}) {
	td.Orders = append(td.Orders, order)
}

func (td *TestData) Clear() {
	td.Users = make([]map[string]interface{}, 0)
	td.Products = make([]map[string]interface{}, 0)
	td.Orders = make([]map[string]interface{}, 0)
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	assert.Equal(t, expected, actual)
}

func AssertNotEqual(t *testing.T, expected, actual interface{}) {
	assert.NotEqual(t, expected, actual)
}

func AssertTrue(t *testing.T, condition bool) {
	assert.True(t, condition)
}

func AssertFalse(t *testing.T, condition bool) {
	assert.False(t, condition)
}

func AssertNil(t *testing.T, object interface{}) {
	assert.Nil(t, object)
}

func AssertNotNil(t *testing.T, object interface{}) {
	assert.NotNil(t, object)
}

func AssertError(t *testing.T, err error) {
	assert.Error(t, err)
}

func AssertNoError(t *testing.T, err error) {
	assert.NoError(t, err)
}

func AssertContains(t *testing.T, container, item interface{}) {
	assert.Contains(t, container, item)
}

func AssertNotContains(t *testing.T, container, item interface{}) {
	assert.NotContains(t, container, item)
}

func AssertEmpty(t *testing.T, object interface{}) {
	assert.Empty(t, object)
}

func AssertNotEmpty(t *testing.T, object interface{}) {
	assert.NotEmpty(t, object)
}

func AssertLen(t *testing.T, object interface{}, length int) {
	assert.Len(t, object, length)
}

func AssertGreater(t *testing.T, e1, e2 interface{}) {
	assert.Greater(t, e1, e2)
}

func AssertLess(t *testing.T, e1, e2 interface{}) {
	assert.Less(t, e1, e2)
}

func AssertGreaterOrEqual(t *testing.T, e1, e2 interface{}) {
	assert.GreaterOrEqual(t, e1, e2)
}

func AssertLessOrEqual(t *testing.T, e1, e2 interface{}) {
	assert.LessOrEqual(t, e1, e2)
}

func AssertWithinDuration(t *testing.T, expected, actual time.Time, delta time.Duration) {
	assert.WithinDuration(t, expected, actual, delta)
}

func AssertRegexp(t *testing.T, rx interface{}, str interface{}) {
	assert.Regexp(t, rx, str)
}

func AssertNotRegexp(t *testing.T, rx interface{}, str interface{}) {
	assert.NotRegexp(t, rx, str)
}

func AssertJSONEq(t *testing.T, expected, actual string) {
	assert.JSONEq(t, expected, actual)
}

func AssertYAMLEq(t *testing.T, expected, actual string) {
	assert.YAMLEq(t, expected, actual)
}

func AssertXMLEq(t *testing.T, expected, actual string) {
	assert.Equal(t, expected, actual)
}

func AssertHTTPBodyContains(t *testing.T, handler gin.HandlerFunc, method, path, body string, expectedText string) {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	router := gin.New()
	router.Handle(method, path, handler)
	router.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), expectedText)
}

func AssertHTTPStatus(t *testing.T, handler gin.HandlerFunc, method, path, body string, expectedStatus int) {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	router := gin.New()
	router.Handle(method, path, handler)
	router.ServeHTTP(w, req)

	assert.Equal(t, expectedStatus, w.Code)
}

func AssertHTTPHeader(t *testing.T, handler gin.HandlerFunc, method, path, body, headerName, expectedValue string) {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	router := gin.New()
	router.Handle(method, path, handler)
	router.ServeHTTP(w, req)

	assert.Equal(t, expectedValue, w.Header().Get(headerName))
}

func BenchmarkHandler(b *testing.B, handler gin.HandlerFunc, method, path, body string) {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))

	router := gin.New()
	router.Handle(method, path, handler)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkFunction(b *testing.B, fn func()) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fn()
	}
}

func BenchmarkFunctionWithParam(b *testing.B, fn func(interface{}), param interface{}) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fn(param)
	}
}
