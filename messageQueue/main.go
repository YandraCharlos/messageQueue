package main

import (
	"messageQueue/controller"
	"messageQueue/repository"
	"messageQueue/service"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	rdb, err := repository.NewRedisDatabase(&repository.RedisConfig{
		URL:       "192.168.99.100",
		Port:      "6379",
		MaxIdle:   10,
		MaxActive: 1000,
		Timeout:   1,
		Wait:      false,
	})
	if err != nil {
		os.Exit(1)
	}

	queueRepo := repository.NewQueueRepository(rdb)
	queueService := service.NewService(queueRepo)
	controller.ApplyController(e, queueService)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":8000"))
}
