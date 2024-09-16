package app

import (
	endpoint "WebServer/internal/endpoint/app"
	imageGenerationPage "WebServer/internal/server/handlers/pages/image-generation"
	recognitionFromFilePage "WebServer/internal/server/handlers/pages/recognition-from-file"
	textProcessingPage "WebServer/internal/server/handlers/pages/text-processing"
	rewritePage "WebServer/internal/server/handlers/pages/text-rewriting"

	audioFormHandler "WebServer/internal/services/formHandlers/audio"
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

	a.Endpoint.RegisterForm("/recognize", audioFormHandler.New())
	a.Endpoint.RegisterForm("/rewriteFromWeb", textFormHandler.New("{{ rewrite }}"))
	a.Endpoint.RegisterForm("/processTextFromWeb", textFormHandler.New(""))
	a.Endpoint.RegisterForm("/generateImage", imageFormHandler.New())

	log.Println("Starting endpoint...")
	a.Endpoint.Start()
}
