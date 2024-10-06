package main

import (
	"context"
	"log"
	"omnihr-coding-test/pkg/api"
	"omnihr-coding-test/pkg/cache"
	"omnihr-coding-test/pkg/config"
	"omnihr-coding-test/pkg/database"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func main() {
	redisClient := cache.NewRedisClient()
	db := database.NewDatabase()
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	appConfig, err := config.LoadConfig("pkg/config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	defer logger.Sync()

	gin.SetMode(gin.DebugMode)

	r := api.NewRouter(logger, db, redisClient, &ctx, appConfig)
	if err := r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
