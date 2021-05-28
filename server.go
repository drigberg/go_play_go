// websocket logic courtesy of https://github.com/snassr/blog-goreactsockets

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var gameManager GameManager

type GameIdData struct {
	GameID string
}

type CreateGameLocalRequest struct {
	UserID string
	Size   int
}

type CreateGameRemoteRequest struct {
	UserID string
	Size   int
}

type JoinGameRemoteRequest struct {
	UserID string
	GameID string
}

type LeaveGameRemoteRequest struct {
	UserID string
	GameID string
}

type RejoinGameRemoteRequest struct {
	UserID string
	GameID string
}

type RejoinGameLocalRequest struct {
	UserID string
	GameID string
}

type GetGameInfoLocalRequest struct {
	UserID string
	GameID string
}

type GetGameInfoRemoteRequest struct {
	UserID string
	GameID string
}

type PlaceStoneLocalRequest struct {
	UserID string
	GameID string
	Coord  Coord
}

type PlaceStoneRemoteRequest struct {
	UserID string
	GameID string
	Coord  Coord
}

type PassLocalRequest struct {
	UserID string
	GameID string
}
type PassRemoteRequest struct {
	UserID string
	GameID string
}

type ErrorDataJoinGameRemote struct {
	Type string "remote/joinGame"
}

type ErrorDataRejoinGameLocal struct {
	Type string "local/rejoinGame"
}

type ErrorDataRejoinGameRemote struct {
	Type string "remote/rejoinGame"
}

type ErrorDataGetGameInfoRemote struct {
	Type string "remote/getGameInfo"
}

type ErrorDataGetGameInfoLocal struct {
	Type string "local/getGameInfo"
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

func onCreateGameLocal(c *SocketClient, data []byte) {
	log.Println("Request: createGameLocal")

	// parse and validate request
	var req CreateGameLocalRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	size := req.Size

	if userID == "" || (size != 9 && size != 13 && size != 19) {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	// Create game
	gameID := gameManager.CreateGameLocal(userID, size, c)

	// set and write response message
	log.Println("Player " + userID + " created game " + gameID)
	c.send = Message{Name: "local/gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()
}

func onCreateGameRemote(c *SocketClient, data []byte) {
	log.Println("Request: createGameRemote")

	// parse and validate request
	var req CreateGameRemoteRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	size := req.Size

	if userID == "" || (size != 9 && size != 13 && size != 19) {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	// Create game
	gameID := gameManager.CreateGameRemote(userID, size, c)

	// set and write response message
	log.Println("Player " + userID + " created game " + gameID)
	c.send = Message{Name: "remote/gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()
}

func sendOtherPlayerUpdate(gameID string, userID string) {
	otherPlayer, err := gameManager.GetOtherPlayerRemote(gameID, userID)
	if err == nil {
		log.Println("Telling player " + otherPlayer.UserID + " to refresh")
		if otherPlayer.SocketClient != nil {
			otherPlayer.SocketClient.send = Message{Name: "remote/update", Data: nil}
			otherPlayer.SocketClient.Write()
		}
	} else {
		log.Println("No other player found!")
	}
}

func onRejoinGameRemote(c *SocketClient, data []byte) {
	log.Println("Request: remote/rejoinGame")

	// parse and validate request
	var req RejoinGameRemoteRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	// Rejoin game if already registered
	joined := gameManager.RejoinGameRemote(gameID, userID, c)

	if !joined {
		log.Println("Player " + userID + " could not rejoin game " + gameID)
		c.send = Message{
			Name: "error",
			Data: ErrorDataRejoinGameRemote{
				Type: "remote/rejoinGame",
			},
		}
		c.Write()
		return
	}

	// set and write response message
	log.Println("Player " + userID + " rejoined game " + gameID)
	c.send = Message{Name: "remote/gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()

	sendOtherPlayerUpdate(gameID, userID)
}

func onRejoinGameLocal(c *SocketClient, data []byte) {
	log.Println("Request: local/rejoinGame")

	// parse and validate request
	var req RejoinGameLocalRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	// Rejoin game if already registered
	joined := gameManager.RejoinGameLocal(gameID, userID, c)

	if !joined {
		log.Println("Player " + userID + " could not rejoin game " + gameID)
		c.send = Message{
			Name: "error",
			Data: ErrorDataRejoinGameLocal{
				Type: "local/rejoinGame",
			},
		}
		c.Write()
		return
	}

	// set and write response message
	log.Println("Player " + userID + " rejoined game " + gameID)
	c.send = Message{Name: "local/gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()
}

func onJoinGameRemote(c *SocketClient, data []byte) {
	log.Println("Request: joinGameRemote")

	// parse and validate request
	var req JoinGameRemoteRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	// Register as second player in existing remote game
	joined := gameManager.JoinGameRemote(gameID, userID, c)

	if !joined {
		log.Println("Player " + userID + " could not join game " + gameID)
		c.send = Message{
			Name: "error",
			Data: ErrorDataJoinGameRemote{
				Type: "remote/joinGame",
			},
		}
		c.Write()
		return
	}

	// set and write response message
	log.Println("Player " + userID + " joined game " + gameID)
	c.send = Message{Name: "remote/gameJoined", Data: GameIdData{GameID: gameID}}
	c.Write()

	sendOtherPlayerUpdate(gameID, userID)
}

func onLeaveGameRemote(c *SocketClient, data []byte) {
	log.Println("Request: leaveGameRemote")

	// parse and validate request
	var req LeaveGameRemoteRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
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
	c.send = Message{Name: "remote/gameLeft", Data: nil}
	c.Write()

	sendOtherPlayerUpdate(gameID, userID)
}

func onGetGameInfoRemote(c *SocketClient, data []byte) {
	log.Println("Request: remote/getGameInfo")

	// parse and validate request
	var req GetGameInfoRemoteRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	gameInfo, err := gameManager.GetGameInfoRemote(gameID, userID)

	if err != nil {
		log.Println("Unable to fetch game info")
		c.send = Message{
			Name: "error",
			Data: ErrorDataGetGameInfoRemote{
				Type: "remote/getGameInfo",
			},
		}
		c.Write()
		return
	}

	// set and write response message
	log.Println("Sending game info to player " + userID)
	c.send = Message{Name: "remote/gameInfo", Data: gameInfo}
	c.Write()
}

func onGetGameInfoLocal(c *SocketClient, data []byte) {
	log.Println("Request: local/getGameInfo")

	// parse and validate request
	var req GetGameInfoLocalRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	gameInfo, err := gameManager.GetGameInfoLocal(gameID, userID)

	if err != nil {
		log.Println("Unable to fetch game info")
		c.send = Message{
			Name: "error",
			Data: ErrorDataGetGameInfoLocal{
				Type: "local/getGameInfo",
			},
		}
		c.Write()
		return
	}

	// set and write response message
	log.Println("Sending game info to player " + userID)
	c.send = Message{Name: "local/gameInfo", Data: gameInfo}
	c.Write()
}

func onPlaceStoneRemote(c *SocketClient, data []byte) {
	log.Println("Request: remote/placeStone")

	// parse and validate request
	var req PlaceStoneRemoteRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID
	coord := req.Coord

	if userID == "" || gameID == "" || coord.X < 0 || coord.Y < 0 {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	placed := gameManager.PlaceStoneRemote(gameID, userID, coord)
	if !placed {
		log.Println("Unable to play move")
		c.send = create400Error("Unable to play move")
		c.Write()
		return
	}

	c.send = Message{Name: "remote/update", Data: nil}
	c.Write()
	sendOtherPlayerUpdate(gameID, userID)
}

func onPlaceStoneLocal(c *SocketClient, data []byte) {
	log.Println("Request: local/placeStone")

	// parse and validate request
	var req PlaceStoneLocalRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID
	coord := req.Coord

	if userID == "" || gameID == "" || coord.X < 0 || coord.Y < 0 {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	placed := gameManager.PlaceStoneLocal(gameID, userID, coord)
	if !placed {
		log.Println("Unable to play move")
		c.send = create400Error("Unable to play move")
		c.Write()
		return
	}

	c.send = Message{Name: "local/update", Data: nil}
	c.Write()
}

func onPassRemote(c *SocketClient, data []byte) {
	log.Println("Request: remote/pass")

	// parse and validate request
	var req PassRemoteRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	passed := gameManager.PassRemote(gameID, userID)
	if !passed {
		log.Println("Unable to pass turn")
		c.send = create400Error("Unable to pass turn")
		c.Write()
		return
	}

	c.send = Message{Name: "remote/update", Data: nil}
	c.Write()
	sendOtherPlayerUpdate(gameID, userID)
}

func onPassLocal(c *SocketClient, data []byte) {
	log.Println("Request: local/pass")

	// parse and validate request
	var req PassLocalRequest
	json.Unmarshal(data, &req)
	userID := req.UserID
	gameID := req.GameID

	if userID == "" || gameID == "" {
		log.Println("Invalid request format")
		c.send = create400Error("invalid request format")
		c.Write()
		return
	}

	passed := gameManager.PassLocal(gameID, userID)
	if !passed {
		log.Println("Unable to pass turn")
		c.send = create400Error("Unable to pass turn")
		c.Write()
		return
	}

	c.send = Message{Name: "local/update", Data: nil}
	c.Write()
}

func onChatRemote(c *SocketClient, data []byte) {
	log.Println("Request: remote/chat")
	log.Printf("Data: %v\n", data)

	// TODO: handle chat

	c.send = Message{Name: "remote/chat", Data: "Chat received!"}
	c.Write()
}

func RunServer(port string) {
	router := NewRouter(port)
	gameManager = NewGameManager()

	// shared actions
	router.Handle("local/createGame", onCreateGameLocal)
	router.Handle("remote/createGame", onCreateGameRemote)
	router.Handle("local/rejoinGame", onRejoinGameLocal)
	router.Handle("remote/rejoinGame", onRejoinGameRemote)
	router.Handle("local/getGameInfo", onGetGameInfoLocal)
	router.Handle("remote/getGameInfo", onGetGameInfoRemote)
	router.Handle("local/placeStone", onPlaceStoneLocal)
	router.Handle("remote/placeStone", onPlaceStoneRemote)
	router.Handle("local/pass", onPassLocal)
	router.Handle("remote/pass", onPassRemote)

	// remote-only actions
	router.Handle("remote/chat", onChatRemote)
	router.Handle("remote/joinGame", onJoinGameRemote)
	router.Handle("remote/leaveGame", onLeaveGameRemote)

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
