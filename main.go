package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Discover all services on the network (e.g. _workstation._tcp)
	// resolver, err := zeroconf.NewResolver(nil)
	// if err != nil {
	// 	log.Fatalln("Failed to initialize resolver:", err.Error())
	// }

	// entries := make(chan *zeroconf.ServiceEntry)
	// go func(results <-chan *zeroconf.ServiceEntry) {
	// 	for entry := range results {
	// 		log.Println(entry.HostName)
	// 	}
	// 	log.Println("No more entries.")
	// }(entries)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	// defer cancel()
	// err = resolver.Browse(ctx, "_workstation._tcp", "local.", entries)
	// if err != nil {
	// 	log.Fatalln("Failed to browse:", err.Error())
	// }

	// <-ctx.Done()

	id := uuid.New()
	// a := NewServer("0.0.0.0:8080", id)
	// a.Serve()

	s, err := NewStore()
	if err != nil {
		fmt.Println(err)
		fmt.Println("couldn't create store")
	}

	myInfo := MyInfo{
		ID:    id,
		Name:  "Grace Roller",
		Email: "gracearoller@gmail.com",
		Phone: "7247993419",
		Lat:   75.5,
		Long:  75.5,
		Time:  time.Now(),
		// Meta: ,
	}
	// fmt.Println(myInfo)

	//
	s.SetMyInfo(&myInfo)
	fmt.Println("Get messages: ")
	fmt.Println(string(s.GetAllMessages(id)))

	s.UpdateLocation(id, 100, 100)
	fmt.Println("Get messages after update: ")
	fmt.Println(string(s.GetAllMessages(id)))

	fmt.Println("Get messages after save messages: ")
	s.UpdateLocation(id, 200, 200)
	s.SaveMessages(s.GetAllMessages(id))
	fmt.Println(string(s.GetAllMessages(id)))

}
