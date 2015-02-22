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

type Client struct {
	receive chan []byte
	broadcast chan []byte
	conn *websocket.Conn
}

type Server struct {
	broadcast chan []byte
	clients []Client
}


func (server Server) broadcaster() {
	for {
		message :=  <-server.broadcast
		fmt.Println("< ", string(message))
		for index, client := range server.clients {
			select {
			case client.receive <- message:
			default:
				// remove the bad client
				server.clients = append(server.clients[:index], server.clients[index+1:]...)
				close(client.receive)
			}
		}
	}
}

func (server Server) Handler(write http.ResponseWriter, read *http.Request) {
	conn, err := upgrader.Upgrade(write, read, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := Client{make(chan []byte), server.broadcast, conn}

	server.clients = append(server.clients, client)

	go client.read()
	client.write()
}

func (client Client) read() {
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		client.broadcast <- message
	}
	client.conn.Close()
}

func (client Client) write() {
	for message := range client.receive {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	client.conn.Close()
}

func main() {
	addr := "localhost:5678"


	fmt.Println("Starting server at ", addr)

	server := Server{make(chan []byte), make([]Client, 1)}

	http.HandleFunc("/", server.Handler)
	go server.broadcaster()

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
