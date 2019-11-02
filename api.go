package main

import (
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// API contains all info for the api
type API struct {
	listenURL string
	router    chi.Router
}

func NewServer(url string) (a *API) {
	api := &API{}
	api.router = chi.NewRouter()
	api.listenURL = url
	api.router.Get("/", a.IndexHandler)
	return api
}

func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func (a *API) Serve() {
	log.Info("Lisenting on ", a.listenURL)
	if err := http.ListenAndServe(a.listenURL, a.router); err != nil {
		log.WithError(err).Error("Unable to serve.")
	}
}
