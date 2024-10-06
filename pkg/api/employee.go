package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"omnihr-coding-test/pkg/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Healthcheck
func Healthcheck(g *gin.Context) {
	g.JSON(http.StatusOK, "ok")
}

// FindEmployees
func FindEmployees(c *gin.Context) {
	// Extract AppContext from Gin context
	appCtx, exists := c.MustGet("appCtx").(*AppContext)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	companyID := c.GetInt64("company_id")

	// Get query params with defaults
	offset, limit, err := getPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get filter params
	filters, err := getEmployeeFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate cache key based on all query params
	cacheKey := generateCacheKey(companyID, offset, limit, filters)

	// Try fetching the data from Redis
	var employees []models.Employee
	if cachedEmployees, err := getFromCache(appCtx, cacheKey); err == nil {
		if err = json.Unmarshal([]byte(cachedEmployees), &employees); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": employees})
			return
		}
	}

	// Query database with filters if cache misses
	if err = queryEmployees(appCtx, companyID, offset, limit, filters, &employees); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Store the result in Redis cache
	if err = cacheResult(appCtx, cacheKey, employees); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set cache"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employees})
}

// Helper functions for get employee
func getPaginationParams(c *gin.Context) (int, int, error) {
	offsetQuery := c.DefaultQuery("offset", "0")
	limitQuery := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid offset format")
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid limit format")
	}

	return offset, limit, nil
}

func getEmployeeFilters(c *gin.Context) (map[string]string, error) {
	statusQuery := c.DefaultQuery("status", "")
	validStatuses := map[string]bool{
		"Active":      true,
		"Not Started": true,
		"Terminated":  true,
	}

	if statusQuery != "" && !validStatuses[statusQuery] {
		return nil, fmt.Errorf("invalid status")
	}

	return map[string]string{
		"department": c.DefaultQuery("department", ""),
		"position":   c.DefaultQuery("position", ""),
		"location":   c.DefaultQuery("location", ""),
		"status":     statusQuery,
	}, nil
}

func generateCacheKey(companyID int64, offset, limit int, filters map[string]string) string {
	return fmt.Sprintf("employees_company_%v_offset_%d_limit_%d_department_%s_position_%s_location_%s_status_%s",
		companyID, offset, limit, filters["department"], filters["position"], filters["location"], filters["status"])
}

func getFromCache(appCtx *AppContext, cacheKey string) (string, error) {
	return appCtx.RedisClient.Get(*appCtx.Ctx, cacheKey).Result()
}

func queryEmployees(appCtx *AppContext, companyID int64, offset, limit int, filters map[string]string, employees *[]models.Employee) error {
	query := appCtx.DB.Offset(offset).Limit(limit).Where("company_id = ?", companyID)

	if department, ok := filters["department"]; ok && department != "" {
		query = query.Where("department = ?", department)
	}
	if position, ok := filters["position"]; ok && position != "" {
		query = query.Where("position = ?", position)
	}
	if location, ok := filters["location"]; ok && location != "" {
		query = query.Where("location = ?", location)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	// get column based on config
	var selectedColumns []string
	for _, orgConfig := range appCtx.Config.Companies {
		if orgConfig.ID == companyID {
			selectedColumns = orgConfig.Columns
			break
		}
	}
	selectedColumnsString := "*"
	if len(selectedColumns) > 0 {
		selectedColumnsString = strings.Join(selectedColumns, ", ")
	}

	// Execute the query
	return query.Select(selectedColumnsString).Find(employees).Error
}

func cacheResult(appCtx *AppContext, cacheKey string, employees []models.Employee) error {
	serializedEmployees, err := json.Marshal(employees)
	if err != nil {
		return err
	}

	return appCtx.RedisClient.Set(*appCtx.Ctx, cacheKey, serializedEmployees, time.Minute).Err()
}
