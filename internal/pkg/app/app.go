package app

import (
	cfg "WebServer/internal/config/app"
	"os"
	"strings"

	endpoint "WebServer/internal/endpoint/app"
	authPage "WebServer/internal/server/handlers/pages/auth"
	imageGenerationPage "WebServer/internal/server/handlers/pages/image-generation"
	recognitionFromFilePage "WebServer/internal/server/handlers/pages/recognition-from-file"
	textProcessingPage "WebServer/internal/server/handlers/pages/text-processing"
	rewritePage "WebServer/internal/server/handlers/pages/text-rewriting"

	audioFormHandler "WebServer/internal/services/formHandlers/audio"
	authHandler "WebServer/internal/services/formHandlers/auth"
	imageFormHandler "WebServer/internal/services/formHandlers/img"
	textFormHandler "WebServer/internal/services/formHandlers/text"

	"log"
)

type App struct {
	Endpoint endpoint.App
}

func New() *App {
	return &App{
		Endpoint: *endpoint.New(),
	}
}

//http.HandleFunc("/text", processTextPage)

func (a *App) Run() {
	a.Endpoint.RegisterPage("/", recognitionFromFilePage.New())
	a.Endpoint.RegisterPage("/image", imageGenerationPage.New())
	a.Endpoint.RegisterPage("/rewrite", rewritePage.New())
	a.Endpoint.RegisterPage("/text", textProcessingPage.New())
	a.Endpoint.RegisterPage("/auth", authPage.New())

	a.Endpoint.RegisterForm("/recognize", audioFormHandler.New())
	a.Endpoint.RegisterForm("/rewriteFromWeb", textFormHandler.New("{{ rewrite }}"))
	a.Endpoint.RegisterForm("/processTextFromWeb", textFormHandler.New(""))
	a.Endpoint.RegisterForm("/generateImage", imageFormHandler.New())
	a.Endpoint.RegisterForm("/login", authHandler.New(cfg.LOGIN, cfg.PASSWORD))

	log.Println("Starting endpoint...")
	tlsMode := os.Getenv("TLS_MODE")
	if strings.ToLower(tlsMode) == "true" {
		log.Println("Starting endpoint with TLS...")
		a.Endpoint.StartTLS() // adjust cert and key files as needed
	} else {
		log.Println("Starting endpoint locally...")
		a.Endpoint.StartLocal()
	}
}
