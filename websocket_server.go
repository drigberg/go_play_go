// websocket logic courtesy of https://github.com/snassr/blog-goreactsockets

package main

import (
	"log"
	"net/http"
)

func onMessage(c *Client, data interface{}) {
	log.Printf("message: %v\n", data)

	// set and write response message
	c.send = Message{Name: "message", Data: "Message received!"}
	c.Write()
}

func RunServer() {
	router := NewRouter()

	router.Handle("message", onMessage)

	// handle all requests to /, upgrade to WebSocket via our router handler.
	http.Handle("/", router)

	// start server.
	log.Println("Listening on port 3001")
	http.ListenAndServe(":3001", nil)
}
