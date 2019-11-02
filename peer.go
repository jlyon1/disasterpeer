package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
	log "github.com/sirupsen/logrus"
)

// Peer represents a peer to communicate with
type Peer struct {
	IP   string
	UUID string
	Port int
}

// RegisterService registers the peer on zeroconf
func RegisterService(myUuid uuid.UUID, port int) {
	for {
		server, err := zeroconf.Register(myUuid.String(), "_disasterpeer._tcp", "local.", port, nil, nil)
		if err != nil {
			panic(err)
		}
		defer server.Shutdown()
		log.Info("Registering Service on ", port)
		// Clean exit.
		<-time.After(time.Second * 50000)
		log.Println("Shutting down.")
	}
}

func FindPeers(peerChan chan Peer) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}
	entries := make(chan *zeroconf.ServiceEntry)

	ctx := context.Background()
	err = resolver.Browse(ctx, "_disasterpeer._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}
	for {
		select {
		case entry := <-entries:
			peerChan <- Peer{
				IP:   entry.HostName,
				UUID: entry.Instance,
				Port: entry.Port,
			}
		}
	}

	<-ctx.Done()
}
