package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"omnihr-coding-test/pkg/cache"
	"omnihr-coding-test/pkg/database"
	"testing"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHealthcheck tests the Healthcheck handler.
func TestHealthcheck(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup the router
	r := gin.Default()
	r.GET("/", Healthcheck)

	// Create a request to pass to our handler.
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, w.Code, "Expected response code to be 200")

	// Check the response body is what we expect.
	assert.Equal(t, "\"ok\"", w.Body.String(), "Expected response body to be 'ok'")
}

func TestFindEmployees(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a router
	r := gin.Default()

	// Setup middleware to inject app context
	r.Use(func(c *gin.Context) {
		ctx := context.Background()
		mockRedis := new(cache.MockRedisClient)
		mockDB := new(database.MockDB)

		// Setup the mock for Redis Get call
		key := "employees_offset_0_limit_10"
		mockRedis.On("Get", ctx, key).Return(redis.NewStringResult("[]", nil)) // Adjust this as necessary

		appCtx := &AppContext{
			RedisClient: mockRedis,
			DB:          mockDB,
			Ctx:         &ctx,
		}
		c.Set("appCtx", appCtx)
		c.Next()
	})

	r.GET("/employees", FindEmployees)

	// Prepare and send the request
	req, _ := http.NewRequest(http.MethodGet, "/employees?offset=0&limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Expected response code to match")
	assert.JSONEq(t, `{"data":[]}`, w.Body.String(), "Expected response body to match JSON")
}
