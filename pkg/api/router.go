package api

import (
	"context"
	"omnihr-coding-test/pkg/cache"
	"omnihr-coding-test/pkg/database"
	"omnihr-coding-test/pkg/middleware"
	"omnihr-coding-test/pkg/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"golang.org/x/time/rate"
)

// AppContext holds shared resources like database and Redis client
type AppContext struct {
	DB          database.DBInterface
	RedisClient cache.CacheInterface
	Ctx         *context.Context
	Config      *models.Config
}

// NewAppContext creates a new AppContext
func NewAppContext(db database.DBInterface, redisClient cache.CacheInterface, ctx *context.Context, config *models.Config) *AppContext {
	return &AppContext{
		DB:          db,
		RedisClient: redisClient,
		Ctx:         ctx,
		Config:      config,
	}
}

func ContextMiddleware(appCtx *AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("appCtx", appCtx)
		c.Next()
	}
}

func NewRouter(logger *zap.Logger, db database.DBInterface, redisClient cache.CacheInterface, ctx *context.Context, cfg *models.Config) *gin.Engine {
	appCtx := NewAppContext(db, redisClient, ctx, cfg)

	r := gin.Default()
	r.Use(ContextMiddleware(appCtx))

	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", Healthcheck)
		v1.GET("/employees", middleware.JWTAuth(), FindEmployees)

		v1.POST("/login", middleware.APIKeyAuth(), LoginHandler)
		v1.POST("/register", middleware.APIKeyAuth(), RegisterHandler)
	}

	return r
}
