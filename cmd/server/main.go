package main

import (
	"log"

	"github.com/maliByatzes/polling/internal/server"
)

func main() {
	srv := server.NewWebsocketServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
