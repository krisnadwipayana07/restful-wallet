package infrastructure

import (
	"errors"
	"fmt"
	"log"

	"github.com/krisnadwipayana07/restful-fintech/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseConnection(config *configs.Config) (*gorm.DB, error) {
	dsn := config.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, errors.New("failed to connect to database")
	}

	fmt.Println("Database connected!")
	return db, nil
}
