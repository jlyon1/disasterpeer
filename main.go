package main

import (
	"fmt"
	"os"
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
	// a := NewServer("0.0.0.0", port, id)

	// portInt, _ := strconv.Atoi(port)
	// fmt.Println("port", portInt)
	// go RegisterService(id, portInt)
	// peerChan := make(chan Peer)
	// go FindPeers(peerChan)
	// go a.Serve()
	// for {
	// 	select {
	// 	case peer := <-peerChan:
	// 		log.Info(peer)
	// 	}
	// }

	s, err := NewStore()
	if err != nil {
		fmt.Println(err)
		fmt.Println("couldn't create store")
	}

	myInfo := MyInfo{
		UserID: id,
		Name:   "hello",
		Email:  "111@gmail.com",
		Phone:  "11111111",
		Lat:    75.5,
		Long:   75.5,
		Time:   time.Now(),
		// Meta: ,
	}
	// fmt.Println(myInfo)

	s.SetMyInfo(&myInfo)
	fmt.Println("Get messages: ")
	fmt.Println(string(s.GetAllMessages(id)))

	s.UpdateLocation(id, 100, 100)
	fmt.Println("Get messages after update: ")
	fmt.Println(string(s.GetAllMessages(id)))

	// fmt.Println("Get messages after save messages: ")
	// s.UpdateLocation(id, 200, 200)
	// s.SaveMessages(s.GetAllMessages(id))
	// fmt.Println(string(s.GetAllMessages(id)))

}
