package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index Page")
}

// Reader is an event listener that reads messages on the server to broadcast
func reader(conn *websocket.Conn) {
	for {
		messageType, byteData, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println(string(byteData))

		if err := conn.WriteMessage(messageType, byteData); err != nil {
			log.Fatal(err)
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Successfully connected...")

	reader(ws)
}

func routeHandlers() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	routeHandlers()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
