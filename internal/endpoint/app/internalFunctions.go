package app

import (
	config "WebServer/internal/config/app"
	"log"
	"net/http"
	"net/url"
)

func (a *App) startRedirection() error {
	log.Println("starting redirect")
	httpServer := &http.Server{
		Addr: ":http",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := config.DOMAIN + ":" + config.HTTPS_PORT
			http.Redirect(w, r, "https://"+host, http.StatusPermanentRedirect)
		}),
	}
	return httpServer.ListenAndServe()
}

func isValidURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}
