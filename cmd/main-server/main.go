package main

import (
	"WebServer/internal/pkg/app"
)

func main() {
	server := app.New()
	server.Run()
}
