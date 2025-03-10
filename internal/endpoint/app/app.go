package app

import (
	"log/slog"

	"os"

	sloggin "github.com/samber/slog-gin"

	"crypto/tls"
	"log"
	"net/http"
	"text/template"

	config "WebServer/internal/config/app"

	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-contrib/gzip"

	"github.com/gin-contrib/rollbar"
	roll "github.com/rollbar/rollbar-go"

	stats "github.com/semihalev/gin-stats"

	"github.com/dvwright/xss-mw"
)

type App struct {
	engine    *gin.Engine
	result    *gin.RouterGroup
	admin     *gin.RouterGroup
	protected *gin.RouterGroup
	store     *persistence.InMemoryStore
	logger    *slog.Logger

	login    string
	password string
}

type PageHandler interface {
	GetPage(c *gin.Context)
}

type FormHandler interface {
	HandleForm(c *gin.Context)
}

func New(protectedAuth gin.HandlerFunc, adminProtectedAuth gin.HandlerFunc) *App {
	config := config.Init()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Use XSS protection
	var xssMdlwr xss.XssMw
	router.Use(xssMdlwr.RemoveXss())

	// Use Recovery
	router.Use(gin.Recovery())

	// Use session
	if config.SESSION_SECRET != "" {
		cookieStore := cookie.NewStore([]byte(config.SESSION_SECRET))
		router.Use(sessions.Sessions("neuron-nexus-session", cookieStore))
	}

	// Use GZIP compression
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Use stats
	router.Use(stats.RequestStats())

	// Create cache
	store := persistence.NewInMemoryStore(1 * time.Hour)

	// Add URL validation function
	router.SetFuncMap(template.FuncMap{
		"isValidURL": isValidURL,
	})

	// Use rollbar
	roll.SetToken(config.ROLLBAR_TOKEN)
	roll.SetEnvironment("development")
	roll.SetCodeVersion("v2")
	roll.SetServerHost("web.1")
	router.Use(rollbar.Recovery(false))

	r := router.Group("/operation")

	a := router.Group("/admin")
	if adminProtectedAuth != nil {
		a.Use(adminProtectedAuth)
	}

	protected := router.Group("/protected")
	if protectedAuth != nil {
		protected.Use(protectedAuth)
	}

	router.StaticFS("/web/", http.Dir("../../web"))
	router.StaticFile("/robots.txt", "../../web/static/robots.txt")
	router.LoadHTMLGlob("../../web/templates/*.html")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Use logger
	router.Use(sloggin.New(logger))

	return &App{
		engine:    router,
		result:    r,
		admin:     a,
		store:     store,
		protected: protected,
		logger:    logger,
	}
}

func (a *App) GetLogger() *slog.Logger {
	return a.logger
}

func (a *App) SetBasicAuth(login, password string) {
	a.login = login
	a.password = password
	config := config.Init()
	if a.login != "" && a.password != "" {
		a.logger.Info("Admin", "login", config.ADMIN_LOGIN, "password", config.ADMIN_PASSWORD)
		a.logger.Info("User", "login", login, "password", password)
		a.logger.Info("Using authorization")
		a.engine.Use(gin.BasicAuth(gin.Accounts{
			config.ADMIN_LOGIN: config.ADMIN_PASSWORD,
			login:              password,
		}))
	}
}

func (a *App) Register404Page(handler gin.HandlerFunc) {
	a.logger.Info("Registering 404 handler")
	a.engine.NoRoute(cache.CachePage(a.store, 5*time.Minute, handler))
}

func (a *App) RegisterResultWithCache(pattern string, handler gin.HandlerFunc) {
	a.logger.Info("Registering result handler for", "pattern", pattern)
	a.result.GET(pattern, cache.CachePage(a.store, 5*time.Minute, handler))
}

func (a *App) RegisterProtectedPage(pattern string, handler gin.HandlerFunc) {
	a.logger.Info("Registering protected handler for", "pattern", pattern)
	a.protected.GET(pattern, handler)
}

func (a *App) RegisterProtectedPageWithCache(pattern string, handler gin.HandlerFunc, t time.Duration) {
	a.logger.Info("Registering protected handler for", "pattern", pattern)
	a.protected.GET(pattern, cache.CachePage(a.store, t, handler))
}

func (a *App) RegisterResultNoCache(pattern string, handler gin.HandlerFunc) {
	a.logger.Info("Registering result handler no cache for", "pattern", pattern)
	a.result.GET(pattern, handler)
}

func (a *App) RegisterResultFormHandler(pattern string, handler gin.HandlerFunc) {
	a.logger.Info("Registering result form handler for", "pattern", pattern)
	a.result.POST(pattern, handler)
}

func (a *App) RegisterPageWithCache(pattern string, handler gin.HandlerFunc) {
	a.logger.Info("Registering page handler with cache for", "pattern", pattern)
	a.engine.GET(pattern, cache.CachePage(a.store, 5*time.Minute, handler))
}

func (a *App) RegisterPageNoCache(pattern string, handler gin.HandlerFunc) {
	a.logger.Info("Registering page handler no cache for", "pattern", pattern)
	a.engine.GET(pattern, handler)
}

func (a *App) RegisterForm(pattern string, handler gin.HandlerFunc) {
	a.logger.Info("Registering form handler for", "pattern", pattern)
	a.engine.POST(pattern, handler)
}

func (a *App) RegisterAdminPageNoCahce(pattern string, handler gin.HandlerFunc) {
	a.logger.Info("Registering admin page handler no cache for", "pattern", pattern)
	a.admin.GET(pattern, handler)
}

func (a *App) StartLocal() {
	config := config.Init()
	a.logger.Info("Starting local server...")
	httpServer := &http.Server{
		Addr:    config.HTTP_PORT,
		Handler: a.engine,
	}
	log.Fatal(httpServer.ListenAndServe())
}

func (a *App) StartTLS() {
	config := config.Init()
	a.logger.Info("Starting TLS server...")

	m := &autocert.Manager{
		Cache:      autocert.DirCache("../../var/www/.cache"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.DOMAIN, "www."+config.DOMAIN, "doc."+config.DOMAIN, "docs."+config.DOMAIN),
		Email:      "info@realnoevremya.ru",
	}
	tlsServer := &http.Server{
		Addr: config.HTTPS_PORT,
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
		Handler: a.engine,
	}

	go func() {
		if err := a.startRedirection(); err != nil {
			log.Fatalf("Redirect server error: %v", err)
		}
	}()
	log.Fatal(tlsServer.ListenAndServeTLS("", ""))
}

func (a *App) GetEngine() *gin.Engine {
	return a.engine
}
