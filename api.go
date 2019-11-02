package main

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// API contains all info for the api
type API struct {
	listenURL string
	router    chi.Router
	UUID      uuid.UUID
	s         *Store
}

// NewServer ...
func NewServer(url string, uuid uuid.UUID) (a *API) {
	api := &API{}
	api.router = chi.NewRouter()
	api.listenURL = url
	api.UUID = uuid

	api.router.Get("/", a.IndexHandler)
	api.router.Get("/uuid", a.GetUUID)

	return api
}

// IndexHandler ...
func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

// GetUUID ...
func (a *API) GetUUID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(a.UUID.String()))
}

// Serve ...
func (a *API) Serve() {
	log.Info("Lisenting on ", a.listenURL)
	if err := http.ListenAndServe(a.listenURL, a.router); err != nil {
		log.WithError(err).Error("Unable to serve.")
	}
}

// func (a *API) GetMessageHandler() []byte {
// }
