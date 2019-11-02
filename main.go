package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Discover all services on the network (e.g. _workstation._tcp)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbname := os.Getenv("DBNAME")
	if dbname == "" {
		dbname = "disaster.db"
	}

	id := uuid.New()

	s, err := NewStore(dbname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("couldn't create store")
	}

	tmp := MyInfo{
		UserID: id,
		Name:   "me1",
		Email:  "me@example.com",
		Phone:  "(555) 555 5555",
		Lat:    -1,
		Long:   -1,
		Time:   time.Now(),
	}
	s.SetMyInfo(&tmp)
	a := NewServer("0.0.0.0", port, id, s)

	portInt, _ := strconv.Atoi(port)

	go RegisterService(id, portInt)
	peerChan := make(chan Peer)
	go FindPeers(peerChan)
	go a.Serve()
	var peerList []Peer
	for {
		select {
		case peer := <-peerChan:
			peerList = append(peerList, peer)
		case <-time.Tick(time.Second):
			var npl []Peer
			for _, p := range peerList {
				resp, err := http.Get("http://" + p.IP + ":" + strconv.Itoa(p.Port) + "/messages")

				if err != nil {
					log.Error("Could not connect to peer", p.IP)
					log.Error(err)
				} else {
					defer resp.Body.Close()

					if resp.StatusCode == http.StatusOK {
						bodyBytes, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							log.Fatal(err)
						}
						npl = append(npl, p)
						var encMessages []EncryptedMessage
						json.Unmarshal(bodyBytes, &encMessages)
						fmt.Println("Got ", len(encMessages), " Messages")

					}
				}
			}
			peerList = npl
		}
	}

}
