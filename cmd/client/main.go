package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	var server string
	flag.StringVar(&server, "server", "", "Server address to connect to.")
	flag.Parse()

	if server == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(server); err != nil {
		log.Fatal(err)
	}
}

func run(server string) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	URL := url.URL{Scheme: "ws", Host: server}

	conn, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() error {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			log.Printf("%s", message)
		}
	}()

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			log.Println("Caught interrupt signal - quitting")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return err
			}
			select {
			case <-done:
			case <-time.After(2 * time.Second):
			}
			return nil
		default:
			msg := "Hello, Server!"
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				return err
			}
			time.Sleep(10 * time.Second)
		}
	}
}
