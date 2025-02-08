package main

import (
	"fmt"

	"github.com/krisnadwipayana07/restful-fintech/configs"
	"github.com/krisnadwipayana07/restful-fintech/internal/delivery/http"
	"github.com/krisnadwipayana07/restful-fintech/internal/domain/repository"
	"github.com/krisnadwipayana07/restful-fintech/internal/domain/service"
	"github.com/krisnadwipayana07/restful-fintech/internal/infrastructure"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load Config
	config, err := configs.InitConfig()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	// Initialize dependencies
	db, err := infrastructure.NewDatabaseConnection(&config)
	if err != nil {
		panic(err)
	}

	redis, err := infrastructure.InitRedisConnection(&config)
	if err != nil {
		panic(err)
	}

	// Initialize repository
	repo, err := repository.New(db)
	if err != nil {
		panic(err)
	}

	// Initialize service
	service, err := service.New(repo, db, redis)
	if err != nil {
		panic(err)
	}

	// Setup routes
	http.InitHandler(e, service)

	// Start server
	port := fmt.Sprintf(":%s", config.Port)
	e.Logger.Fatal(e.Start(port))

}
