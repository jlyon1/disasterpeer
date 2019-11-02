package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Discover all services on the network (e.g. _workstation._tcp)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	id := uuid.New()
	a := NewServer("0.0.0.0", port, id)

	portInt, _ := strconv.Atoi(port)
	fmt.Println("port", portInt)
	go RegisterService(id, portInt)
	peerChan := make(chan Peer)
	go FindPeers(peerChan)
	go a.Serve()
	for {
		select {
		case peer := <-peerChan:
			log.Info(peer)
		}
	}

}
