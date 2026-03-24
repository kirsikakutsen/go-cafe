package routes

import (
	"go-cafe/internal/handlers"

	"github.com/labstack/echo/v5"
)

func SetupRouter(e *echo.Echo) {
	api := e.Group("/api")

	api.POST("/auth/login", handlers.Login)
}