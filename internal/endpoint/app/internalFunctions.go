package app

import (
	config "WebServer/internal/config/app"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func (a *App) redirectToTls(w http.ResponseWriter, r *http.Request) {
	log.Println("redirect to TLS server for ", r.RemoteAddr)
	url := fmt.Sprintf("https://%s%s", config.IP, config.HTTPS_PORT)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (a *App) startRedirection() error {
	log.Println("starting redirect")
	return http.ListenAndServe(config.HTTP_PORT, http.HandlerFunc(a.redirectToTls))
}

func (a *App) startTLSServer() error {

	s, _ := filepath.Abs("./")
	log.Printf("Serving from: %s\n", s)

	certFile := "../../internal/config/ssl/domain.crt"
	keyFile := "../../internal/config/ssl/domain.key"

	log.Println("starting HTTPS server")

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("../../web"))))

	return http.ListenAndServeTLS(config.HTTPS_PORT, certFile, keyFile, nil)
}
