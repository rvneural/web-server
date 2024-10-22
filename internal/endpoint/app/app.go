package app

import (
	"crypto/tls"
	"log"
	"net/http"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

type PageHandler interface {
	GetPage(c *gin.Context)
}

type FormHandler interface {
	HandleForm(c *gin.Context)
}

// [x] New App
func New() *App {
	router := gin.Default()
	router.StaticFS("/web/", http.Dir("../../web"))
	router.LoadHTMLGlob("../../web/templates/*.html")
	return &App{engine: router}
}

// [x] Page Registration
func (a *App) RegisterPage(pattern string, handler PageHandler) {
	log.Println("Registering handler for pattern", pattern)
	a.engine.GET(pattern, handler.GetPage)
}

// [x] Form Registration
func (a *App) RegisterForm(pattern string, handler FormHandler) {
	log.Println("Registering form handler for pattern", pattern)
	a.engine.POST(pattern, handler.HandleForm)
}

// [x] Handler
func (a *App) StartLocal() {
	log.Println("Starting local server...")
	httpServer := &http.Server{
		Addr:    ":80",
		Handler: a.engine,
	}
	log.Fatal(httpServer.ListenAndServe())
}

// [x] TLS Handler
func (a *App) StartTLS() {
	log.Println("Starting TLS server...")

	m := &autocert.Manager{
		Cache:      autocert.DirCache("../../var/www/.cache"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("neuron-nexus.ru"),
	}
	tlsServer := &http.Server{
		Addr: ":443",
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
