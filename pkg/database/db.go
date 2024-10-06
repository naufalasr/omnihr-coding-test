package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"omnihr-coding-test/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBInterface interface {
	Offset(offset int) *gorm.DB
	Limit(limit int) *gorm.DB
	Find(interface{}, ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
}

func NewDatabase() *gorm.DB {
	var database *gorm.DB
	var err error

	db_hostname := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	db_user := os.Getenv("POSTGRES_USER")
	db_pass := os.Getenv("POSTGRES_PASSWORD")
	db_port := os.Getenv("POSTGRES_PORT")

	dbURl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_pass, db_hostname, db_port, db_name)
	for i := 1; i <= 3; i++ {
		database, err = gorm.Open(postgres.Open(dbURl), &gorm.Config{})
		if err == nil {
			break
		} else {
			log.Printf("Attempt %d: Failed to initialize database. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}
	database.AutoMigrate(&models.Employee{})
	database.AutoMigrate(&models.User{})

	return database
}
