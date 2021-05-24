// websocket logic courtesy of https://github.com/snassr/blog-goreactsockets

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var gameManager GameManager

func onMessage(c *SocketClient, data []byte) {
	log.Printf("message: %v\n", data)

	// set and write response message
	c.send = Message{Name: "message", Data: "Message received!"}
	c.Write()
}

type GameIdData struct {
	GameID string
}

type CreateGameRequest struct {
	UserID string
	Size   int
}

type ErrorDataJoinGame struct {
	Type string "joinGame"
}

type ErrorDataGetGameInfo struct {
	Type string "getGameInfo"
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

func onCreateGameLocal(c *SocketClient, data []byte) {
	log.Println("Request: createGameLocal")

	// parse and validate request
	var req CreateGameRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	size := req.Size

	if userID == "" || (size != 9 && size != 13 && size != 19) {
		log.Println("Invalid request format")
		c.send = create400Error("Invalid request format")
		c.Write()
		return
	}

	// Create game
	gameID := gameManager.CreateGameLocal(userID, size, c)

	// set and write response message
	log.Println("Player " + userID + " created game " + gameID)
	c.send = Message{Name: "gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()
}

func onCreateGameRemote(c *SocketClient, data []byte) {
	log.Println("Request: createGameRemote")

	// parse and validate request
	var req CreateGameRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	size := req.Size

	if userID == "" || (size != 9 && size != 13 && size != 19) {
		log.Println("Invalid request format")
		c.send = create400Error("Invalid request format")
		c.Write()
		return
	}

	// Create game
	gameID := gameManager.CreateGameRemote(userID, size, c)

	// set and write response message
	log.Println("Player " + userID + " created game " + gameID)
	c.send = Message{Name: "gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()
}

type JoinGameRequest struct {
	UserID string
	GameID string
}

func sendOtherPlayerUpdate(gameID string, userID string) {
	otherPlayer, err := gameManager.GetOtherPlayer(gameID, userID)
	if err == nil {
		log.Println("Telling player " + otherPlayer.UserID + " to refresh")
		if otherPlayer.SocketClient != nil {
			otherPlayer.SocketClient.send = Message{Name: "update", Data: nil}
			otherPlayer.SocketClient.Write()
		}
	} else {
		log.Println("No other player found!")
	}
}

func onJoinGameRemote(c *SocketClient, data []byte) {
	log.Println("Request: joinGameRemote")

	// parse and validate request
	var req JoinGameRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("Invalid request format")
		c.Write()
		return
	}

	// Rejoin game if already registered, or register as part of game
	joined := gameManager.RejoinGame(gameID, userID, c)
	if !joined {
		joined = gameManager.JoinGameRemote(gameID, userID, c)
	}

	if !joined {
		log.Println("Player " + userID + " could not join game " + gameID)
		c.send = createJoinGameError()
		c.Write()
		return
	}

	// set and write response message
	log.Println("Player " + userID + " joined game " + gameID)
	c.send = Message{Name: "gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()

	sendOtherPlayerUpdate(gameID, userID)
}

type LeaveGameRequest struct {
	UserID string
	GameID string
}

func onLeaveGameRemote(c *SocketClient, data []byte) {
	log.Println("Request: leaveGameRemote")

	// parse and validate request
	var req LeaveGameRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("Invalid request format")
		c.Write()
		return
	}

	// Mark game as over
	left := gameManager.LeaveGameRemote(gameID, userID)
	if !left {
		log.Println("Player " + userID + " could not leave game " + gameID)
		c.send = create400Error("Unable to leave game")
		c.Write()
		return
	}

	// set and write response message
	log.Println("Player " + userID + " left game " + gameID)
	c.send = Message{Name: "gameLeft", Data: nil}
	c.Write()

	sendOtherPlayerUpdate(gameID, userID)
}

type GetGameInfoRequest struct {
	UserID string
	GameID string
}

func onGetGameInfo(c *SocketClient, data []byte) {
	log.Println("Request: getGameInfo")

	// parse and validate request
	var req GetGameInfoRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
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

type PlaceStoneRequest struct {
	UserID string
	GameID string
	Coord  Coord
}

func onPlaceStone(c *SocketClient, data []byte) {
	log.Println("Request: placeStone")

	// parse and validate request
	var req PlaceStoneRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID
	coord := req.Coord

	if userID == "" || gameID == "" || coord.X < 0 || coord.Y < 0 {
		log.Println("Invalid request format")
		c.send = create400Error("Invalid request format")
		c.Write()
		return
	}

	placed := gameManager.PlaceStone(gameID, userID, coord)
	if !placed {
		log.Println("Unable to play move")
		c.send = create400Error("Unable to play move")
		c.Write()
		return
	}

	c.send = Message{Name: "update", Data: nil}
	c.Write()
	sendOtherPlayerUpdate(gameID, userID)
}

type PassRequest struct {
	UserID string
	GameID string
}

func onPass(c *SocketClient, data []byte) {
	log.Println("Request: pass")

	// parse and validate request
	var req PassRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("Invalid request format")
		c.Write()
		return
	}

	passed := gameManager.Pass(gameID, userID)
	if !passed {
		log.Println("Unable to pass turn")
		c.send = create400Error("Unable to pass turn")
		c.Write()
		return
	}

	c.send = Message{Name: "update", Data: nil}
	c.Write()
	sendOtherPlayerUpdate(gameID, userID)
}

func RunServer(port string) {
	router := NewRouter(port)
	gameManager = NewGameManager()
	router.Handle("message", onMessage)
	router.Handle("createGameLocal", onCreateGameLocal)
	router.Handle("createGameRemote", onCreateGameRemote)
	router.Handle("joinGameRemote", onJoinGameRemote)
	router.Handle("getGameInfo", onGetGameInfo)
	router.Handle("placeStone", onPlaceStone)
	router.Handle("pass", onPass)
	router.Handle("leaveGameRemote", onLeaveGameRemote)

	// handle all requests to /, upgrade to WebSocket via our router handler.
	http.Handle("/socket", router)

	if os.Getenv("ENV") == "PRODUCTION" {
		r := http.NewServeMux()
		buildHandler := http.FileServer(http.Dir("app/build"))
		r.Handle("/", buildHandler)
		http.Handle("/", r)
		log.Println("Production: serving client app")
	}

	// start server.
	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
