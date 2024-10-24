package app

import (
	"crypto/tls"
	"log"
	"net/http"

	config "WebServer/internal/config/app"

	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type App struct {
	engine *gin.Engine
	result *gin.RouterGroup
	store  *persistence.InMemoryStore

	login    string
	password string
}

type PageHandler interface {
	GetPage(c *gin.Context)
}

type FormHandler interface {
	HandleForm(c *gin.Context)
}

func New() *App {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	r := router.Group("/operation")
	store := persistence.NewInMemoryStore(time.Minute)

	router.StaticFS("/web/", http.Dir("../../web"))

	router.LoadHTMLGlob("../../web/templates/*.html")

	return &App{
		engine: router,
		result: r,
		store:  store,
	}
}

func (a *App) SetBasicAuth(login, password string) {
	a.login = login
	a.password = password

	if a.login != "" && a.password != "" {
		log.Println("Using authorization")
		a.engine.Use(gin.BasicAuth(gin.Accounts{
			login: password,
		}))
	}
}

func (a *App) Register404Page(handler PageHandler) {
	log.Println("Registering 404 page")
	a.engine.NoRoute(cache.CachePage(a.store, time.Minute, handler.GetPage))
}

func (a *App) RegisterResult(pattern string, handler PageHandler) {
	log.Println("Registering result handler for pattern", pattern)
	a.result.GET(pattern, cache.CachePage(a.store, time.Minute, handler.GetPage))
}

func (a *App) RegisterPage(pattern string, handler PageHandler) {
	log.Println("Registering handler for pattern", pattern)
	a.engine.GET(pattern, handler.GetPage)
}

func (a *App) RegisterForm(pattern string, handler FormHandler) {
	log.Println("Registering form handler for pattern", pattern)
	a.engine.POST(pattern, handler.HandleForm)
}

func (a *App) StartLocal() {
	log.Println("Starting local server...")
	httpServer := &http.Server{
		Addr:    config.HTTP_PORT,
		Handler: a.engine,
	}
	log.Fatal(httpServer.ListenAndServe())
}

func (a *App) StartTLS() {
	log.Println("Starting TLS server...")

	m := &autocert.Manager{
		Cache:      autocert.DirCache("../../var/www/.cache"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.DOMAIN),
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
