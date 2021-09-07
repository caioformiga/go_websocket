package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func setRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func homePage(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Home Page of WebSocket Hello Wolrd!")
}

func wsEndpoint(writer http.ResponseWriter, req *http.Request) {

	// CheckOrigin returns true if the request Origin header is acceptable. If
	// CheckOrigin is nil, then a safe default is used: return false if the
	// Origin request header is present and the origin host is not equal to
	// request Host header.
	// This implementation always returns true,
	upgrader.CheckOrigin = func(req *http.Request) bool { return true }
	log.Print("WebSocket Endpoint of WebSocket Hello Wolrd!")

	wsConn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		log.Print("[upgrader.Upgrade]", err)
		return
	}
	log.Print("Successfully connected to WebSocket...")

	startReader(wsConn)
}

func startReader(wsConn *websocket.Conn) {
	// defer statement places wsConn.Close() before this function ends
	defer wsConn.Close()

	for {
		messageType, p, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("[wsConn.ReadMessage]:", err)
			break
		}

		log.Printf("recv: %s", string(p))

		err = wsConn.WriteMessage(messageType, p)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func starServer() {
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Go WebSocket using github.com/gorilla/websocket is read to handle connections...")
	setRoutes()
	starServer()
}
