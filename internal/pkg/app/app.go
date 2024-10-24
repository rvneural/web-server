package app

import (
	"os"
	"strings"

	endpoint "WebServer/internal/endpoint/app"

	imageGenerationPage "WebServer/internal/server/handlers/pages/image-generation"
	upscalePage "WebServer/internal/server/handlers/pages/image-upscaler"
	recognitionFromFilePage "WebServer/internal/server/handlers/pages/recognition-from-file"
	textProcessingPage "WebServer/internal/server/handlers/pages/text-processing"
	rewritePage "WebServer/internal/server/handlers/pages/text-rewriting"

	audioFormHandler "WebServer/internal/server/handlers/forms/audio"
	imageFormHandler "WebServer/internal/server/handlers/forms/img/generator"
	imageUpscalerFormHandler "WebServer/internal/server/handlers/forms/img/upscale"
	photopea "WebServer/internal/server/handlers/forms/photopea"
	textFormHandler "WebServer/internal/server/handlers/forms/text"

	audioResult "WebServer/internal/server/handlers/pages/results/audio"
	imageResult "WebServer/internal/server/handlers/pages/results/image"
	textResult "WebServer/internal/server/handlers/pages/results/text"

	notFound "WebServer/internal/server/handlers/pages/404"

	"log"
)

type App struct {
	Endpoint endpoint.App
	login    string
	password string
	tlsMode  bool
}

func New() *App {
	return &App{
		Endpoint: *endpoint.New(),
	}
}

func (a *App) init() {

	a.login = os.Getenv("LOGIN")
	a.password = os.Getenv("PASSWORD")
	tlsMode := os.Getenv("TLS_MODE")

	if strings.ToLower(tlsMode) == "true" {
		a.tlsMode = true
	} else {
		a.tlsMode = false
	}

	a.Endpoint.RegisterPage("/", recognitionFromFilePage.New())
	a.Endpoint.RegisterPage("/image", imageGenerationPage.New())
	a.Endpoint.RegisterPage("/rewrite", rewritePage.New())
	a.Endpoint.RegisterPage("/text", textProcessingPage.New())
	a.Endpoint.RegisterPage("/upscale", upscalePage.New())

	a.Endpoint.RegisterResult("/audio", audioResult.New())
	a.Endpoint.RegisterResult("/text", textResult.New())
	a.Endpoint.RegisterResult("/image", imageResult.New())

	a.Endpoint.RegisterForm("/recognize", audioFormHandler.New())
	a.Endpoint.RegisterForm("/rewriteFromWeb", textFormHandler.New("{{ rewrite }}"))
	a.Endpoint.RegisterForm("/processTextFromWeb", textFormHandler.New(""))
	a.Endpoint.RegisterForm("/generateImage", imageFormHandler.New())
	a.Endpoint.RegisterForm("/upscaleImage", imageUpscalerFormHandler.New())
	a.Endpoint.RegisterForm("/photopea", photopea.New())

	a.Endpoint.Register404Page(notFound.New())
}

func (a *App) Run() {

	a.init()

	if a.login != "" && a.password != "" {
		a.Endpoint.SetBasicAuth(a.login, a.password)
	}

	if a.tlsMode {
		log.Println("Starting endpoint with TLS...")
		a.Endpoint.StartTLS()
	} else {
		log.Println("Starting endpoint locally...")
		a.Endpoint.StartLocal()
	}
}
