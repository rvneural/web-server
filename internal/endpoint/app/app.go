package app

import (
	"crypto/tls"
	"log"
	"net/http"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
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

func (a *App) StartLocal() {
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("../../web"))))
	log.Println("Starting local server...")
	httpServer := &http.Server{
		Addr: ":80",
	}
	log.Fatal(httpServer.ListenAndServe())
}

func (a *App) StartTLS() {
	log.Println("Starting TLS server...")
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("../../web"))))
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
	}

	go func() {
		if err := a.startRedirection(); err != nil {
			log.Fatalf("Redirect server error: %v", err)
		}
	}()
	log.Fatal(tlsServer.ListenAndServeTLS("", ""))
}
