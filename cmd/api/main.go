package main

import (
	"go-cafe/internal/database"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}
	database.Connect()
}