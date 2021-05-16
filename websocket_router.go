// websocket logic courtesy of https://github.com/snassr/blog-goreactsockets

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Handler is a type representing functions which resolve requests.
type Handler func(*SocketClient, []byte)

// Event is a type representing request names.
type Event string

// Router is a message routing object mapping events to function handlers.
type Router struct {
	Port  string
	rules map[Event]Handler
}

// NewRouter returns an initialized Router.
func NewRouter(port string) *Router {
	return &Router{
		Port:  port,
		rules: make(map[Event]Handler),
	}
}

// ServeHTTP creates the socket connection and begins the read routine.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// configure upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// only allow requests from app port
			// origin := r.Header.Get("Origin")
			// TODO: clean up
			// log.Println("Origin:", r.Header.Get("Origin"))
			// return origin == "http://localhost:"+rt.AppPort || origin == "http://127.0.0.1:"+rt.AppPort || origin == "http://0.0.0.0:"+rt.AppPort
			return true
		},
	}

	// upgrade connection to socket
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("socket server configuration error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := NewClient(socket, rt.FindHandler)

	// running method for reading from sockets, in main routine
	client.Read()
}

func (rt *Router) FindHandler(event Event) (Handler, bool) {
	handler, found := rt.rules[event]
	return handler, found
}

// Handle is a function to add handlers to the router.
func (rt *Router) Handle(event Event, handler Handler) {
	rt.rules[event] = handler
}
