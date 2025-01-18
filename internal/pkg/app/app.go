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

	indexPage "WebServer/internal/server/handlers/pages/index"

	"WebServer/internal/services/auth"
	dbWorker "WebServer/internal/services/db"

	imageOperationList "WebServer/internal/server/handlers/pages/admin/images"
	adminOperationList "WebServer/internal/server/handlers/pages/admin/operations"
	adminUserPage "WebServer/internal/server/handlers/pages/admin/user"
	adminUsersPage "WebServer/internal/server/handlers/pages/admin/user/allusers"

	saveSystem "WebServer/internal/server/handlers/forms/saving"

	newsPage "WebServer/internal/server/handlers/pages/feed"
	mediaPage "WebServer/internal/server/handlers/pages/mediafeed"
	rssFeed "WebServer/internal/server/handlers/pages/rss"

	fips "WebServer/internal/server/handlers/pages/fips"

	authPages "WebServer/internal/server/handlers/pages/auth"
	authMaster "WebServer/internal/services/auth"

	userPage "WebServer/internal/server/handlers/pages/user"

	"WebServer/internal/server/handlers/pages/admin/stats"
)

type App struct {
	Endpoint       endpoint.App
	login          string
	password       string
	tlsMode        bool
	idMaxLen       int
	logger         *slog.Logger
	dataBaseWorker *dbWorker.Worker
	auth           *auth.AuthentificationHandler
}

func New() *App {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	dataBaseWorker := dbWorker.New(logger)

	auth := authMaster.New(dataBaseWorker)
	return &App{
		Endpoint:       *endpoint.New(auth.AuthMiddleware("/login", 0, dataBaseWorker), auth.AuthMiddleware("/login", 1, dataBaseWorker)),
		idMaxLen:       35,
		logger:         logger,
		dataBaseWorker: dataBaseWorker,
		auth:           auth,
	}
}

func (a *App) init() {

	a.login = os.Getenv("LOGIN")
	a.password = os.Getenv("PASSWORD")
	a.tlsMode = os.Getenv("TLS_MODE") == "true"

	if a.login != "" && a.password != "" {
		a.Endpoint.SetBasicAuth(a.login, a.password)
	}

	// Основные страницы
	a.Endpoint.RegisterProtectedPage("/recognition", recognitionFromFilePage.New().GetPage)
	a.Endpoint.RegisterProtectedPage("/image", imageGenerationPage.New().GetPage)
	a.Endpoint.RegisterProtectedPage("/text", textProcessingPage.New().GetPage)
	a.Endpoint.RegisterProtectedPage("/imgprocess", upscalePage.New().GetPage)
	a.Endpoint.RegisterProtectedPage("/", userPage.New(a.dataBaseWorker, a.logger).GetPage)

	a.Endpoint.RegisterProtectedPage("/news", newsPage.New(a.logger).GetPage)
	a.Endpoint.RegisterProtectedPage("/media", mediaPage.New(a.logger).GetPage)
	a.Endpoint.RegisterProtectedPage("/fips", fips.New().GetPage)
	//----------------------------

	a.Endpoint.RegisterForm("/login", a.auth.HandleLogin)
	a.Endpoint.RegisterForm("/register", a.auth.HandleRegistration)

	authPage := authPages.New()
	a.Endpoint.RegisterPageNoCache("/", indexPage.New().GetPage)
	a.Endpoint.RegisterPageNoCache("/login", authPage.GetLoginPage)
	a.Endpoint.RegisterPageNoCache("/register", authPage.GetRegisterPage)
	a.Endpoint.RegisterPageNoCache("/logout", authPage.GetLogoutPage)

	a.Endpoint.RegisterPageNoCache("/rss", rssFeed.New(a.logger, "RSS_URL").GetPage)
	a.Endpoint.RegisterPageNoCache("/media", rssFeed.New(a.logger, "MEDIA_RSS_URL").GetPage)

	a.Endpoint.RegisterAdminPageNoCahce("/stats", stats.New().GetPage)

	a.Endpoint.RegisterAdminPageNoCahce("/operations", adminOperationList.New(a.dataBaseWorker).GetPage)
	a.Endpoint.RegisterAdminPageNoCahce("/images", imageOperationList.New(a.dataBaseWorker).GetPage)
	a.Endpoint.RegisterAdminPageNoCahce("/users/", adminUsersPage.New(a.dataBaseWorker).GetPage)
	a.Endpoint.RegisterAdminPageNoCahce("/users/:id", adminUserPage.New(a.dataBaseWorker, a.logger).GetPage)

	a.Endpoint.RegisterResultNoCache("/get", newID.New(a.dataBaseWorker).GetPage)
	a.Endpoint.RegisterProtectedPage("/operation/:id", result.New(notFoundOperationPage.New(), progressOperationPage.New(), a.dataBaseWorker, a.logger).GetPage)

	a.Endpoint.RegisterForm("/recognize", audioFormHandler.New(a.dataBaseWorker, a.logger).HandleForm)
	a.Endpoint.RegisterForm("/rewriteFromWeb", textFormHandler.New("{{ rewrite }}", a.dataBaseWorker, a.logger).HandleForm)
	a.Endpoint.RegisterForm("/processTextFromWeb", textFormHandler.New("", a.dataBaseWorker, a.logger).HandleForm)
	a.Endpoint.RegisterForm("/generateImage", imageFormHandler.New(a.dataBaseWorker, a.logger).HandleForm)
	a.Endpoint.RegisterForm("/upscaleImage", imageUpscalerFormHandler.New(a.logger).HandleForm)
	a.Endpoint.RegisterForm("/removeBackground", bgRemover.New(a.logger).HandleForm)
	a.Endpoint.RegisterForm("/photopea", photopea.New(a.logger).HandleForm)

	a.Endpoint.RegisterResultFormHandler("/saveOperation", saveSystem.New(a.dataBaseWorker, a.logger).HandleForm)
	a.Endpoint.RegisterResultFormHandler("/getVersion", saveSystem.NewVersionSystem(a.dataBaseWorker, a.logger).HandleForm)

	a.Endpoint.Register404Page(notFound.New().GetPage)
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
