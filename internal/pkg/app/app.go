package app

import (
	"log/slog"
	"os"

	endpoint "WebServer/internal/endpoint/app"

	imageGenerationPage "WebServer/internal/server/handlers/pages/image-generation"
	upscalePage "WebServer/internal/server/handlers/pages/image-upscaler"
	recognitionFromFilePage "WebServer/internal/server/handlers/pages/recognition-from-file"
	textProcessingPage "WebServer/internal/server/handlers/pages/text-processing"

	audioFormHandler "WebServer/internal/server/handlers/forms/audio"
	imageFormHandler "WebServer/internal/server/handlers/forms/img/generator"

	bgRemover "WebServer/internal/server/handlers/forms/img/rembg"
	imageUpscalerFormHandler "WebServer/internal/server/handlers/forms/img/upscale"

	photopea "WebServer/internal/server/handlers/forms/photopea"
	textFormHandler "WebServer/internal/server/handlers/forms/text"

	result "WebServer/internal/server/handlers/pages/results"

	notFound "WebServer/internal/server/handlers/pages/404"

	newID "WebServer/internal/server/operations"

	notFoundOperationPage "WebServer/internal/server/handlers/pages/results/notfound"
	progressOperationPage "WebServer/internal/server/handlers/pages/results/progress"

	dbWorker "WebServer/internal/services/db"

	adminOperationList "WebServer/internal/server/handlers/pages/admin/operations"

	saveSystem "WebServer/internal/server/handlers/forms/saving"

	newsPage "WebServer/internal/server/handlers/pages/feed"
	rssFeed "WebServer/internal/server/handlers/pages/rss"

	"WebServer/internal/server/handlers/pages/stats"
)

type App struct {
	Endpoint endpoint.App
	login    string
	password string
	tlsMode  bool
	idMaxLen int
	logger   *slog.Logger
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
	a.logger = a.Endpoint.GetLogger()

	if a.login != "" && a.password != "" {
		a.Endpoint.SetBasicAuth(a.login, a.password)
	}

	a.Endpoint.RegisterPageWithCache("/", recognitionFromFilePage.New())
	a.Endpoint.RegisterPageWithCache("/image", imageGenerationPage.New())
	a.Endpoint.RegisterPageWithCache("/text", textProcessingPage.New())
	a.Endpoint.RegisterPageWithCache("/imgprocess", upscalePage.New())
	a.Endpoint.RegisterPageNoCache("/news", newsPage.New(a.logger))
	a.Endpoint.RegisterPageNoCache("/rss", rssFeed.New(a.logger))

	a.Endpoint.RegisterAdminPageNoCahce("/stats", stats.New())

	dataBaseWorker := dbWorker.New(a.logger)

	a.Endpoint.RegisterAdminPageNoCahce("/operations", adminOperationList.New(dataBaseWorker))

	a.Endpoint.RegisterResultNoCache("/get", newID.New(dataBaseWorker))
	a.Endpoint.RegisterResultNoCache("/:id", result.New(notFoundOperationPage.New(), progressOperationPage.New(), dataBaseWorker, a.logger))

	a.Endpoint.RegisterForm("/recognize", audioFormHandler.New(dataBaseWorker, a.logger))
	a.Endpoint.RegisterForm("/rewriteFromWeb", textFormHandler.New("{{ rewrite }}", dataBaseWorker, a.logger))
	a.Endpoint.RegisterForm("/processTextFromWeb", textFormHandler.New("", dataBaseWorker, a.logger))
	a.Endpoint.RegisterForm("/generateImage", imageFormHandler.New(dataBaseWorker, a.logger))
	a.Endpoint.RegisterForm("/upscaleImage", imageUpscalerFormHandler.New(a.logger))
	a.Endpoint.RegisterForm("/removeBackground", bgRemover.New(a.logger))
	a.Endpoint.RegisterForm("/photopea", photopea.New(a.logger))

	a.Endpoint.RegisterResultFormHandler("/saveOperation", saveSystem.New(dataBaseWorker, a.logger))
	a.Endpoint.RegisterResultFormHandler("/getVersion", saveSystem.NewVersionSystem(dataBaseWorker, a.logger))

	a.Endpoint.Register404Page(notFound.New())
}

func (a *App) Run() {

	a.init()

	if a.tlsMode {
		a.logger.Info("Starting endpoint with TLS...")
		a.Endpoint.StartTLS()
	} else {
		a.logger.Info("Starting endpoint without TLS...")
		a.Endpoint.StartLocal()
	}
}
