package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

var broadcast chan []byte
var clients []chan []byte

func broadcaster() {
	for {
		message:=  <-broadcast
		fmt.Println(string(message))
		for index, channel := range clients {
			select {
			case channel <- message:
			default:
				// remove the bad client
				clients = append(clients[:index], clients[index+1:]...)
				close(channel)
			}
		}
	}
}

func client_read(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		broadcast <- message
	}
	conn.Close()
}

func client_write(conn *websocket.Conn, messages chan []byte) {
	for message := range messages {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	conn.Close()
}

func wsHandler(write http.ResponseWriter, read *http.Request) {
	conn, err := upgrader.Upgrade(write, read, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := make(chan []byte)
	clients = append(clients, client)
	go client_read(conn)
	client_write(conn, client)
}

func main() {
	addr := "localhost:5678"

	http.HandleFunc("/", wsHandler)

	fmt.Println("Starting server at ", addr)

	broadcast = make(chan []byte)
	go broadcaster()

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
