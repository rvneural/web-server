package app

import (
	"os"

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

	result "WebServer/internal/server/handlers/pages/results"

	notFound "WebServer/internal/server/handlers/pages/404"

	newID "WebServer/internal/server/operations"

	"WebServer/internal/services/idgenerator"

	notFoundOperationPage "WebServer/internal/server/handlers/pages/results/notfound"
	progressOperationPage "WebServer/internal/server/handlers/pages/results/progress"

	dbConfig "WebServer/internal/config/db"
	dbWorker "WebServer/internal/services/db"

	adminOperationList "WebServer/internal/server/handlers/pages/admin/operations"

	"WebServer/internal/server/handlers/pages/stats"

	"log"
)

type App struct {
	Endpoint endpoint.App
	login    string
	password string
	tlsMode  bool
	idMaxLen int
}

func New() *App {
	return &App{
		Endpoint: *endpoint.New(),
		idMaxLen: 35,
	}
}

func (a *App) init() {

	a.login = os.Getenv("LOGIN")
	a.password = os.Getenv("PASSWORD")
	a.tlsMode = os.Getenv("TLS_MODE") == "true"

	a.Endpoint.RegisterPageWithCache("/", recognitionFromFilePage.New())
	a.Endpoint.RegisterPageWithCache("/image", imageGenerationPage.New())
	a.Endpoint.RegisterPageWithCache("/rewrite", rewritePage.New())
	a.Endpoint.RegisterPageWithCache("/text", textProcessingPage.New())
	a.Endpoint.RegisterPageWithCache("/upscale", upscalePage.New())

	a.Endpoint.RegisterAdminPageNoCahce("/stats", stats.New())

	notFoundOperationPageP := notFoundOperationPage.New()
	progressOperationPageP := progressOperationPage.New()

	dataBaseWorker := dbWorker.New(
		dbConfig.HOST, dbConfig.PORT, dbConfig.LOGIN, dbConfig.PASSWORD, dbConfig.DB_NAME, dbConfig.RESULT_TABLE_NAME,
	)

	a.Endpoint.RegisterAdminPageNoCahce("/operations", adminOperationList.New(dataBaseWorker))

	a.Endpoint.RegisterResultNoCache("/get", newID.New(idgenerator.New(a.idMaxLen)))
	a.Endpoint.RegisterResultWithCache("/:id", result.New(notFoundOperationPageP, progressOperationPageP, dataBaseWorker))

	a.Endpoint.RegisterForm("/recognize", audioFormHandler.New(dataBaseWorker))
	a.Endpoint.RegisterForm("/rewriteFromWeb", textFormHandler.New("{{ rewrite }}", dataBaseWorker))
	a.Endpoint.RegisterForm("/processTextFromWeb", textFormHandler.New("", dataBaseWorker))
	a.Endpoint.RegisterForm("/generateImage", imageFormHandler.New(dataBaseWorker))
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
