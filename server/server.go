package main

import (
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type EncryptedMessage struct {
	ID   int
	Sent time.Time
	Body []byte
}

type API struct {
	listenURL string
	port      string
	router    http.Handler
}

type Info struct {
	// ID     int
	UserID uuid.UUID // user ID
	Name   string
	Email  string
	Phone  string
	Lat    float64
	Long   float64
	Time   time.Time
	Status string
	// Meta  MetaImages
}

type InfoSet map[string]bool

var infoLookup InfoSet
var allInfo []Info

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Error(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Error(err)
	}
	return key
}

func (a *API) UpdateMessages(w http.ResponseWriter, r *http.Request) {
	var messages []EncryptedMessage
	err := json.NewDecoder(r.Body).Decode(&messages)

	if err != nil {
		w.WriteHeader(400)
		log.WithError(err).Error("Bad Request")
		w.Write([]byte("Malformed request"))
		return
	}

	priv, err := ioutil.ReadFile("mykey.pem")
	if err != nil {
		log.WithError(err)
		panic(err)
	}
	key := BytesToPrivateKey(priv)
	hash := sha512.New()

	for _, m := range messages {
		bytes, err := rsa.DecryptOAEP(hash, nil, key, m.Body, nil)
		if err != nil {
			log.Error(err)
		}
		// fmt.Println(string(bytes))

		var info Info
		err = json.Unmarshal(bytes, &info)
		if err != nil {
			fmt.Println(err)
		}
		infoKey := info.UserID.String() + info.Time.String()
		if !infoLookup[infoKey] {
			allInfo = append(allInfo, info)
			infoLookup[infoKey] = true
		}
	}

	file, err := json.MarshalIndent(allInfo, "", "")

	if err != nil {
		fmt.Println(err)
	}
	_ = ioutil.WriteFile("data.json", file, 0644)

	w.WriteHeader(200)
}

func NewServer(l string, p string) *API {
	infoLookup = make(InfoSet)

	a := &API{
		listenURL: l,
		port:      p,
	}
	r := chi.NewRouter()
	r.Post("/update", a.UpdateMessages)
	a.router = r

	return a
}

// Serve ...
func (a *API) Serve() {
	log.Info("Lisenting on ", a.listenURL+":"+a.port)
	if err := http.ListenAndServe(a.listenURL+":"+a.port, a.router); err != nil {
		log.WithError(err).Error("Unable to serve.")
	}
}
