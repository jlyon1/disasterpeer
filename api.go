package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
func NewServer(url string, port string, uuid uuid.UUID, s *Store) *API {
	api := API{}
	api.router = chi.NewRouter()
	api.listenURL = url
	api.port = port
	api.UUID = uuid
	api.s = s
	api.router.Get("/", api.IndexHandler)
	api.router.Get("/app.js", api.ScriptHandler)
	api.router.Get("/vue.js", api.VueHandler)
	api.router.Get("/bulma.css", api.CSSHandler)
	api.router.Get("/uuid", api.GetUUID)
	api.router.Get("/info", api.GetMyInfo)
	api.router.Get("/messages", api.GetAllMessages)
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

// ScriptHandler ...
func (a *API) VueHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/vue.js")
}

// ScriptHandler ...
func (a *API) CSSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/bulma.css")
}

func (a *API) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	x := a.s.GetAllMessages(a.UUID)

	WriteJSON(w, x)
}

// GetUUID ...
func (a *API) GetUUID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(a.UUID.String()))
}

func (a *API) GetMyInfo(w http.ResponseWriter, r *http.Request) {
	x, _ := a.s.GetMyInfo(a.UUID)
	WriteJSON(w, x)
}

func (a *API) PostMyInfo(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var tmp MyInfo
	err := dec.Decode(&tmp)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	err = a.s.SetMyInfo(&MyInfo{
		UserID: a.UUID,
		Name:   tmp.Name,
		Email:  tmp.Email,
		Phone:  tmp.Phone,
		Lat:    tmp.Lat,
		Long:   tmp.Long,
		Time:   time.Now(),
		Status: tmp.Status,
	})
	if err != nil {
		fmt.Println(err)
	}
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
// WriteJSON writes the data as JSON.
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}
