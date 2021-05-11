package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type Ping struct {
	Ok bool
}

func health(writer http.ResponseWriter, request *http.Request) {
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer connection.Close()
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// TODO: switch-case for message.operation
		// TODO: get response message

		ping := &Ping{Ok: true}
		b, err := json.Marshal(ping)

		if err != nil {
				log.Println("marshal:",err)
				return
		}

		err = connection.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func RunServer(host string, port string) {
	http.HandleFunc("/health", health)

	var addr = host + ":" + port
	log.Println("Listening at " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}