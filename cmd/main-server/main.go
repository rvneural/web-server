package main

import (
	"WebServer/internal/pkg/app"
	"log"
)

func main() {
	server := app.New()
	log.Println("Starting application...")
	server.Run()
}
