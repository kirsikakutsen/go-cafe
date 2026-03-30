package main

import (
	"go-cafe/internal/database"
	"go-cafe/internal/routes"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}
	database.Connect()

	e := echo.New()

	routes.SetupRouter(e)

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}