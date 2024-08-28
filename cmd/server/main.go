package main

import (
	"log"

	"github.com/maliByatzes/polling/internal/server"
)

const port = "8080"

func main() {
	srv := server.NewWebsocketServer(":" + port)
	log.Printf("Server running on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
