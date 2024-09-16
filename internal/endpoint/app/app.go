package app

import (
	"log"
	"net/http"
)

type App struct {
}

type PageHandler interface {
	GetPage(w http.ResponseWriter, r *http.Request)
}

type FormHandler interface {
	HandleForm(w http.ResponseWriter, r *http.Request)
}

func New() *App {
	return &App{}
}

func (a *App) RegisterPage(pattern string, handler PageHandler) {
	log.Println("Registering handler for pattern", pattern)
	http.HandleFunc(pattern, handler.GetPage)
}

func (a *App) RegisterForm(pattern string, handler FormHandler) {
	log.Println("Registering form handler for pattern", pattern)
	http.HandleFunc(pattern, handler.HandleForm)
}

func (a *App) RegisterStatic(pattern string, handler http.Handler) {
	log.Println("Registering static handler for pattern", pattern)
	http.Handle(pattern, handler)
}

func (a *App) Start() {

	log.Println("Starting server...")

	go func() {
		if err := a.startRedirection(); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	log.Fatal(a.startTLSServer())
}
