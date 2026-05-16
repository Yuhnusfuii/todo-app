package main

import (
	"fmt"
	"os"
	"todo-app/config"
	"todo-app/delivery/http"
	"todo-app/delivery/http/middleware"
	"todo-app/domain"
	"todo-app/repository/postgres"
	"todo-app/usecase"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	tasks  = make(map[int]domain.Task, 0)
	nextID = 1
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Use(middleware.Recovery())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(middleware.RateLimiter())
	db := config.InitDB(cfg.Database)

	userRepo := postgres.NewUserRepository(db)
	userUC := usecase.NewUserUsecase(userRepo, cfg.JWT)
	http.NewUserHandler(app, userUC)

	taskRepo := postgres.NewTaskPostgresRepo(db)
	taskUC := usecase.NewTaskUsecase(taskRepo)
	http.NewTaskHandler(app, taskUC, cfg.JWT.Secret)

	app.Listen(":3000")
}
