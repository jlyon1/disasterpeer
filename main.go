package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Discover all services on the network (e.g. _workstation._tcp)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	id := uuid.New()

	myInfo := MyInfo{
		ID:     0,
		UserID: id,
		Name:   "Grace asdfRoller",
		Email:  "gracearoller@gmail.com",
		Phone:  "7247993419",
		Lat:    75.5,
		Long:   75.5,
		Time:   time.Now(),
		// Meta: ,
	}
	s, err := NewStore()
	if err != nil {
		fmt.Println(err)
		fmt.Println("couldn't create store")
	}
	s.SetMyInfo(&myInfo)

	a := NewServer("0.0.0.0", port, id, s)
	// s.SetMyInfo(&myInfo)
	// s.UpdateLocation(id, 100, 100)
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
