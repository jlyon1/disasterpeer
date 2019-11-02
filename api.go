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
	port      string
	router    chi.Router
	UUID      uuid.UUID
	s         *Store
}

// NewServer ...
func NewServer(url string, port string, uuid uuid.UUID) *API {
	api := API{}
	api.router = chi.NewRouter()
	api.listenURL = url
	api.port = port
	api.UUID = uuid

	api.router.Get("/", api.IndexHandler)
	api.router.Get("/app.js", api.ScriptHandler)
	api.router.Get("/uuid", api.GetUUID)
	api.router.Get("/info", api.GetMyInfo)
	api.router.Post("/info", api.PostMyInfo)

	return &api
}

// IndexHandler ...
func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

// ScriptHandler ...
func (a *API) ScriptHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/app.js")
}

// GetUUID ...
func (a *API) GetUUID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(a.UUID.String()))
}

func (a *API) GetMyInfo(w http.ResponseWriter, r *http.Request) {

}

func (a *API) PostMyInfo(w http.ResponseWriter, r *http.Request) {

}

// Serve ...
func (a *API) Serve() {
	log.Info("Lisenting on ", a.listenURL+":"+a.port)
	if err := http.ListenAndServe(a.listenURL+":"+a.port, a.router); err != nil {
		log.WithError(err).Error("Unable to serve.")
	}
}

// func (a *API) GetMessageHandler() []byte {
// }
