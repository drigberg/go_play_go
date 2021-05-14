// websocket logic courtesy of https://github.com/snassr/blog-goreactsockets

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var gameManager GameManager

func onMessage(c *SocketClient, data []byte) {
	log.Printf("message: %v\n", data)

	// set and write response message
	c.send = Message{Name: "message", Data: "Message received!"}
	c.Write()
}

type GameIdData struct {
	GameID int
}

type CreateGameRequest struct {
	UserID string
}

type ErrorData struct {
	Message string
}

func onCreateGame(c *SocketClient, data []byte) {
	log.Println("Request: createGame")

	// parse and validate request
	var req CreateGameRequest
	json.Unmarshal(data, &req)
	userID := req.UserID

	if userID == "" {
		log.Println("Invalid request format")
		c.send = Message{Name: "error", Data: ErrorData{Message: "Invalid request format"}}
		c.Write()
		return
	}

	// Create game
	gameID := gameManager.CreateGame(userID, c)

	// set and write response message
	log.Println("Player " + userID + " created game " + strconv.Itoa(gameID))
	c.send = Message{Name: "gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()
}

type JoinGameRequest struct {
	UserID string
	GameID int
}

func onJoinGame(c *SocketClient, data []byte) {
	log.Println("Request: joinGame")

	// parse and validate request
	var req JoinGameRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID <= 0 {
		log.Println("Invalid request format")
		c.send = Message{Name: "error", Data: ErrorData{Message: "Invalid request format"}}
		c.Write()
		return
	}

	// Rejoin game if already registered, or register as part of game
	joined := gameManager.RejoinGame(gameID, userID, c)
	if !joined {
		joined = gameManager.JoinGame(gameID, userID, c)
	}

	if !joined {
		log.Println("Player " + userID + " could not join game " + strconv.Itoa(gameID))
		c.send = Message{Name: "error", Data: ErrorData{Message: "cannot join game " + strconv.Itoa(gameID)}}
		c.Write()
		return
	}

	// set and write response message
	log.Println("Player " + userID + " joined game " + strconv.Itoa(gameID))
	c.send = Message{Name: "gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()
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
