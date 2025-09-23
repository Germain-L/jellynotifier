package main

import (
	"log"

	"jellynotifier/server"
)

func main() {
	// Create and start the server
	srv := server.New()
	log.Fatal(srv.Start())
}
