package main

import (
	"WebServer/internal/pkg/app"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	server := app.New()
	server.Run()
}
