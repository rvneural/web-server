package app

import (
	endpoint "WebServer/internal/endpoint/app"
	"log"
)

type App struct {
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	server := endpoint.New()
	log.Println("Starting endpoint...")
	server.Start()
}
