// websocket logic courtesy of https://github.com/snassr/blog-goreactsockets

package main

import (
	"log"
	"net/http"
	"strconv"
)

var gameManager GameManager

func onMessage(c *SocketClient, data interface{}) {
	log.Printf("message: %v\n", data)

	// set and write response message
	c.send = Message{Name: "message", Data: "Message received!"}
	c.Write()
}

type GameIdData struct {
	GameID int
}

type CreateGameRequest struct {
	UserId string
}

func onCreateGame(c *SocketClient, data interface{}) {
	log.Println("Request: createGame")
	userID := "some-user-id"
	gameID := gameManager.CreateGame("some-user-id", c)

	// set and write response message
	log.Println("Player " + userID + " created game " + strconv.Itoa(gameID))
	c.send = Message{Name: "gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()
}

func onJoinGame(c *SocketClient, data interface{}) {
	log.Println("Request: joinGame")

	userID := "something"
	gameID := 1

	// Rejoin game if already registered, or register as part of game
	joined := gameManager.RejoinGame(gameID, userID, c)
	if !joined {
		joined = gameManager.JoinGame(gameID, userID, c)
	}
	if joined {
		// set and write response message
		log.Println("Player " + userID + " joined game " + strconv.Itoa(gameID))
		c.send = Message{Name: "gameJoined", Data: GameIdData{GameID: gameID}}
		c.Write()
	}
}

func RunServer() {
	router := NewRouter()
	gameManager = NewGameManager()
	router.Handle("message", onMessage)
	router.Handle("createGame", onCreateGame)
	router.Handle("joinGame", onJoinGame)

	// handle all requests to /, upgrade to WebSocket via our router handler.
	http.Handle("/", router)

	// start server.
	log.Println("Listening on port 3001")
	http.ListenAndServe(":3001", nil)
}
