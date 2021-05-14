// websocket logic courtesy of https://github.com/snassr/blog-goreactsockets

package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Message is an object used to pass data on sockets.
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// FindHandler is a type that defines handler finding functions.
type FindHandler func(Event) (Handler, bool)

// SocketClient is a type that reads and writes on sockets.
type SocketClient struct {
	send        Message
	socket      *websocket.Conn
	findHandler FindHandler
}

// NewClient accepts a socket and returns an initialized SocketClient.
func NewClient(socket *websocket.Conn, findHandler FindHandler) *SocketClient {
	return &SocketClient{
		socket:      socket,
		findHandler: findHandler,
	}
}

// Write receives messages from the channel and writes to the socket.
func (c *SocketClient) Write() {
	msg := c.send
	err := c.socket.WriteJSON(msg)
	if err != nil {
		log.Printf("socket write error: %v\n", err)
	}
}

// Read intercepts messages on the socket and assigns them to a handler function.
func (c *SocketClient) Read() {
	var msg Message
	for {
		// read incoming message from socket
		if err := c.socket.ReadJSON(&msg); err != nil {
			log.Printf("socket read error: %v\n", err)
			break
		}
		// assign message to a function handler
		if handler, found := c.findHandler(Event(msg.Name)); found {
			dataJsonString, err := json.Marshal(msg.Data)
			if err != nil {
				log.Println(err)
				return
			}

			handler(c, dataJsonString)
		}
	}
	log.Println("exiting read loop")

	// close interrupted socket connection
	c.socket.Close()
}
