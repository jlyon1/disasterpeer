package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello server")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	a := NewServer("0.0.0.0", port)
	go a.Serve()

	for {

	}

}
