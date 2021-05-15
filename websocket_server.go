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

type ErrorDataJoinGame struct {
	Type string "joinGame"
}

type ErrorDataGetGameInfo struct {
	Type string "getGameINfo"
}

type ErrorData400 struct {
	Type    string "400"
	Message string
}

func create400Error(message string) Message {
	return Message{
		Name: "error",
		Data: ErrorData400{
			Type:    "400",
			Message: message,
		},
	}
}

func createJoinGameError() Message {
	return Message{
		Name: "error",
		Data: ErrorDataJoinGame{
			Type: "joinGame",
		},
	}
}

func createGetGameInfoError() Message {
	return Message{
		Name: "error",
		Data: ErrorDataGetGameInfo{
			Type: "getGameInfo",
		},
	}
}

func onCreateGame(c *SocketClient, data []byte) {
	log.Println("Request: createGame")

	// parse and validate request
	var req CreateGameRequest
	json.Unmarshal(data, &req)
	userID := req.UserID

	if userID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("Invalid request format")
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
		c.send = create400Error("Invalid request format")
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
		c.send = createJoinGameError()
		c.Write()
		return
	}

	// set and write response message
	log.Println("Player " + userID + " joined game " + strconv.Itoa(gameID))
	c.send = Message{Name: "gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()

	otherPlayer, err := gameManager.GetOtherPlayer(gameID, userID)
	if err == nil {
		otherPlayer.SocketClient.send = Message{Name: "update", Data: nil}
		otherPlayer.SocketClient.Write()
	} else {
		log.Println("No other player found!")
	}
}

type GetGameInfoRequest struct {
	UserID string
	GameID int
}

func onGetGameInfo(c *SocketClient, data []byte) {
	log.Println("Request: getGameInfo")

	// parse and validate request
	var req GetGameInfoRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID <= 0 {
		log.Println("Invalid request format")
		c.send = create400Error("Invalid request format")
		c.Write()
		return
	}

	gameInfo, err := gameManager.GetGameInfo(gameID, userID)

	if err != nil {
		log.Println("Unable to fetch game info")
		c.send = createGetGameInfoError()
		c.Write()
		return
	}

	// set and write response message
	log.Println("Sending game info to player " + userID)
	c.send = Message{Name: "gameInfo", Data: gameInfo}
	c.Write()
}

// TODO: onPlaceStone
// TODO: onPass
// TODO: onMessage
// TODO: onGetInfo (check if player is in room first)

func RunServer() {
	router := NewRouter()
	gameManager = NewGameManager()
	router.Handle("message", onMessage)
	router.Handle("createGame", onCreateGame)
	router.Handle("joinGame", onJoinGame)
	router.Handle("getGameInfo", onGetGameInfo)

	// handle all requests to /, upgrade to WebSocket via our router handler.
	http.Handle("/", router)

	// start server.
	log.Println("Listening on port 3001")
	http.ListenAndServe(":3001", nil)
}
