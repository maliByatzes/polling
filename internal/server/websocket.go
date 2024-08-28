package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWebsocketServer(addr string) *http.Server {
	wssrv := newWsServer()
	r := mux.NewRouter()
	r.HandleFunc("/", wssrv.handleListen)
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type wsServer struct {
	Conn *websocket.Conn
}

func newWsServer() *wsServer {
	return &wsServer{
		Conn: &websocket.Conn{},
	}
}

func (s *wsServer) handleListen(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Conn = conn

	go handleConnection(s.Conn, w)
}

func handleConnection(conn *websocket.Conn, w http.ResponseWriter) {
	log.Printf("Client connected: %s\n", conn.RemoteAddr())
	for {
		msgType, msgContent, err := conn.ReadMessage()
		timeReceived := time.Now()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msgRsp := fmt.Sprintf("Received message: %s. Time received: %v.\n", string(msgContent), timeReceived.Format(time.UnixDate))
		fmt.Print(msgRsp)

		if err := conn.WriteMessage(msgType, []byte(msgRsp)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
