package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

func ValidationErrorResponse(c *gin.Context, errors ValidationErrors) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Error:   "Validation failed",
		Data:    errors,
	})
}

func PaginatedResponse(c *gin.Context, data interface{}, total int64, page, limit int) {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func GetIntQuery(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	return defaultValue
}

func GetIntParam(c *gin.Context, key string) (int, error) {
	value := c.Param(key)
	return strconv.Atoi(value)
}

func GetStringQuery(c *gin.Context, key, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetBoolQuery(c *gin.Context, key string, defaultValue bool) bool {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	if boolValue, err := strconv.ParseBool(value); err == nil {
		return boolValue
	}

	return defaultValue
}

func GetFloatQuery(c *gin.Context, key string, defaultValue float64) float64 {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}

	return defaultValue
}

func GetStringSliceQuery(c *gin.Context, key string) []string {
	value := c.Query(key)
	if value == "" {
		return []string{}
	}

	return strings.Split(value, ",")
}

func GetIntSliceQuery(c *gin.Context, key string) ([]int, error) {
	values := GetStringSliceQuery(c, key)
	result := make([]int, 0, len(values))

	for _, value := range values {
		if intValue, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			result = append(result, intValue)
		} else {
			return nil, fmt.Errorf("invalid integer value: %s", value)
		}
	}

	return result, nil
}

func GetFloatSliceQuery(c *gin.Context, key string) ([]float64, error) {
	values := GetStringSliceQuery(c, key)
	result := make([]float64, 0, len(values))

	for _, value := range values {
		if floatValue, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			result = append(result, floatValue)
		} else {
			return nil, fmt.Errorf("invalid float value: %s", value)
		}
	}

	return result, nil
}

func GetSortQuery(c *gin.Context, defaultSort string) (string, string) {
	sort := GetStringQuery(c, "sort", defaultSort)

	parts := strings.Split(sort, ",")
	if len(parts) != 2 {
		return defaultSort, "asc"
	}

	field := strings.TrimSpace(parts[0])
	direction := strings.ToLower(strings.TrimSpace(parts[1]))

	if direction != "asc" && direction != "desc" {
		direction = "asc"
	}

	return field, direction
}

func GetDateRangeQuery(c *gin.Context, startKey, endKey string) (time.Time, time.Time, error) {
	startStr := c.Query(startKey)
	endStr := c.Query(endKey)

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start date format: %s", startStr)
		}
	} else {
		start = time.Now().AddDate(0, 0, -30)
	}

	if endStr != "" {
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end date format: %s", endStr)
		}
	} else {
		end = time.Now()
	}

	return start, end, nil
}

func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", fmt.Errorf("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID type")
	}

	return userIDStr, nil
}

func GetUserRole(c *gin.Context) (string, error) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", fmt.Errorf("user role not found in context")
	}

	userRoleStr, ok := userRole.(string)
	if !ok {
		return "", fmt.Errorf("invalid user role type")
	}

	return userRoleStr, nil
}

func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if !exists {
		return ""
	}

	requestIDStr, ok := requestID.(string)
	if !ok {
		return ""
	}

	return requestIDStr
}

func GetValidatedData[T any](c *gin.Context) (T, error) {
	var zero T

	data, exists := c.Get("validated_data")
	if !exists {
		return zero, fmt.Errorf("validated data not found in context")
	}

	typedData, ok := data.(T)
	if !ok {
		return zero, fmt.Errorf("invalid data type")
	}

	return typedData, nil
}

func GetPaginationParams(c *gin.Context) (int, int, int) {
	page := GetIntQuery(c, "page", 1)
	limit := GetIntQuery(c, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	return page, limit, offset
}

func SetCacheHeaders(c *gin.Context, ttl time.Duration) {
	c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int(ttl.Seconds())))
	c.Header("Expires", time.Now().Add(ttl).Format(time.RFC1123))
}

func SetNoCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

func SetETag(c *gin.Context, etag string) {
	c.Header("ETag", etag)
}

func CheckETag(c *gin.Context, etag string) bool {
	ifNoneMatch := c.GetHeader("If-None-Match")
	return ifNoneMatch == etag
}

func SendNotModified(c *gin.Context) {
	c.Status(http.StatusNotModified)
}

func SendFile(c *gin.Context, filepath, filename string) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.File(filepath)
}

func SendJSONFile(c *gin.Context, data interface{}, filename string) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to marshal JSON")
		return
	}

	c.Data(http.StatusOK, "application/json", jsonData)
}

func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip == "" {
		ip = c.GetHeader("X-Real-IP")
	}
	if ip == "" {
		ip = c.ClientIP()
	}

	if strings.Contains(ip, ",") {
		ip = strings.Split(ip, ",")[0]
	}

	return strings.TrimSpace(ip)
}

func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

func IsMobile(c *gin.Context) bool {
	userAgent := GetUserAgent(c)
	mobileKeywords := []string{"Mobile", "Android", "iPhone", "iPad", "iPod", "BlackBerry", "Windows Phone"}

	for _, keyword := range mobileKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}

	return false
}

func IsBot(c *gin.Context) bool {
	userAgent := GetUserAgent(c)
	botKeywords := []string{"bot", "crawler", "spider", "scraper", "curl", "wget"}

	userAgentLower := strings.ToLower(userAgent)
	for _, keyword := range botKeywords {
		if strings.Contains(userAgentLower, keyword) {
			return true
		}
	}

	return false
}
