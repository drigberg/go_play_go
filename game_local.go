package main

import (
	"sync"
)

// State can be one of:
// - PLAYING
// - GAME_OVER

type GameLocal struct {
	M            sync.Mutex
	Game         Game
	ID           string
	UserID       string
	SocketClient *SocketClient
	State        string
}

// GameLocalInterface defines methods a GameLocal must implement
type GameLocalInterface interface {
	RejoinGame(userID string, socketClient *SocketClient) bool
	GetInfo() GameInfoLocal
	CurrentTurnColor() string
	Pass()
	PlaceStone(coord Coord) bool
}

// assert that GameLocal implements GameLocalInterface
var _ GameLocalInterface = (*GameLocal)(nil)

// New creates an empty board
func NewGameLocal(gameID string, userID string, size int, socketClient *SocketClient) GameLocal {
	player := Player{
		UserID:       userID,
		SocketClient: socketClient,
	}
	players := make(map[string]*Player)
	players[userID] = &player

	return GameLocal{
		ID:           gameID,
		UserID:       userID,
		State:        "PLAYING",
		SocketClient: socketClient,
		Game:         NewGame(size),
	}
}

type GameInfoLocal struct {
	Size             int
	Turn             int
	ScoreData        ScoreData
	State            string
	CurrentTurnColor string
	AvailableSpaces  []Coord
	Spaces           Spaces
}

func (gameLocal *GameLocal) RejoinGame(userID string, socketClient *SocketClient) bool {
	// return false if player is not part of game
	if gameLocal.UserID != userID {
		return false
	}

	gameLocal.M.Lock()
	defer gameLocal.M.Unlock()
	gameLocal.SocketClient = socketClient
	return true
}

func (gameLocal *GameLocal) CurrentTurnColor() string {
	if gameLocal.Game.Turn%2 == 1 {
		return "BLACK"
	}
	return "WHITE"
}

func (gameLocal *GameLocal) PlaceStone(coord Coord) bool {
	color := gameLocal.CurrentTurnColor()
	placed := gameLocal.Game.PlaceStone(color, coord)
	return placed
}

func (gameLocal *GameLocal) Pass() {
	gameLocal.M.Lock()
	defer gameLocal.M.Unlock()

	// If both players pass, the game is over
	gameOver := gameLocal.Game.Pass()
	if gameOver {
		gameLocal.State = "GAME_OVER"
	}
}

// Returns all the information that the client needs for the game state
func (gameLocal *GameLocal) GetInfo() GameInfoLocal {
	color := gameLocal.CurrentTurnColor()
	spaces := Spaces{
		BLACK: gameLocal.Game.Board.ListSpacesForColor(gameLocal.Game.Board.GetSpaces(), BLACK),
		WHITE: gameLocal.Game.Board.ListSpacesForColor(gameLocal.Game.Board.GetSpaces(), WHITE),
	}

	return GameInfoLocal{
		Size:             gameLocal.Game.Board.Size,
		CurrentTurnColor: color,
		State:            gameLocal.State,
		ScoreData:        gameLocal.Game.Board.GetScoreData(),
		AvailableSpaces:  gameLocal.Game.Board.GetAvailableSpaces(color),
		Spaces:           spaces,
		Turn:             gameLocal.Game.Turn,
	}
}
