package app

import (
	"log"
	"net/http"
)

func (a *App) startRedirection() error {
	log.Println("starting redirect")
	httpServer := &http.Server{
		Addr: ":80",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := "neuron-nexus.ru"
			http.Redirect(w, r, "https://"+host, http.StatusPermanentRedirect)
		}),
	}
	return httpServer.ListenAndServe()
}
