package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/grandcat/zeroconf"
	log "github.com/sirupsen/logrus"
)

// Peer represents a peer to communicate with
type Peer struct {
	IP   string
	UUID string
}

// RegisterService registers the peer on zeroconf
func RegisterService() {
	server, err := zeroconf.Register("GoZeroconf", "_disasterpeer._tcp", "local.", 42424, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down.")
}
