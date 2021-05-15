package main

import (
	"errors"
	"math/rand"
	"sync"
)

// GameManager handles all requests and game states
type GameManager struct {
	M     sync.Mutex
	games map[string]*Game
}

const idChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

// GameManagerInterface defines methods a Game should implement
type GameManagerInterface interface {
	CreateGame(userID string, size int, socketClient *SocketClient) *Game
	GetGameInfo(gameID string, userID string) (GameInfo, error)
	GetOtherPlayer(gameID string, userID string) (*Player, error)
	LeaveGame(gameID string, userID string) bool
	JoinGame(gameID string, userID string, socketClient *SocketClient) bool
	RejoinGame(gameID string, userID string, socketClient *SocketClient) bool
	Pass(gameID string, userID string) bool
	PlaceStone(gameID string, userID string, coord Coord) bool
}

// assert that GameManager implements GameManagerInterface
var _ GameManagerInterface = (*GameManager)(nil)

// NewServer creates a GameManager instance
func NewGameManager() GameManager {
	return GameManager{
		games: make(map[string]*Game),
	}
}

func (gameManager *GameManager) createGameId() string {
	letters := []rune(idChars)
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (gameManager *GameManager) CreateGame(userID string, size int, socketClient *SocketClient) *Game {
	gameManager.M.Lock()
	defer gameManager.M.Unlock()

	player := Player{
		UserID:       userID,
		SocketClient: socketClient,
	}
	players := make(map[string]*Player)
	players[userID] = &player

	gameID := gameManager.createGameId()
	game := Game{
		ID:            gameID,
		State:         "WAITING_FOR_OPPONENT",
		FirstPlayerID: userID,
		Players:       players,
		Turn:          1,
		Board:         NewBoard(size),
	}

	gameManager.games[gameID] = &game
	return &game
}

func (gameManager *GameManager) JoinGame(gameID string, userID string, socketClient *SocketClient) bool {
	game := gameManager.games[gameID]
	if game == nil {
		return false
	}

	if len(game.Players) >= 2 {
		return false
	}

	game.M.Lock()
	defer game.M.Unlock()

	player := Player{
		UserID:       userID,
		SocketClient: socketClient,
	}

	game.Players[userID] = &player
	game.State = "PLAYING"
	return true
}

func (gameManager *GameManager) LeaveGame(gameID string, userID string) bool {
	game := gameManager.games[gameID]
	// return false if no game or player is not part of game
	if game == nil || game.Players[userID] == nil {
		return false
	}
	game.M.Lock()
	defer game.M.Unlock()

	if game.State != "GAME_OVER_PASSED" {
		game.State = "GAME_OVER_FORFEIT"
	}

	game.Players[userID].SocketClient = nil
	return true
}

func (gameManager *GameManager) RejoinGame(gameID string, userID string, socketClient *SocketClient) bool {
	game := gameManager.games[gameID]
	// return false if no game or player is not part of game
	if game == nil || game.Players[userID] == nil {
		return false
	}
	game.M.Lock()
	defer game.M.Unlock()
	game.Players[userID].SocketClient = socketClient
	return true
}

func (gameManager *GameManager) GetGameInfo(gameID string, userID string) (GameInfo, error) {
	game := gameManager.games[gameID]
	// return false if no game or player is not part of game
	if game == nil || game.Players[userID] == nil {
		return GameInfo{}, errors.New("Cannot get game info")
	}
	return game.GetInfo(userID), nil
}

func (gameManager *GameManager) PlaceStone(gameID string, userID string, coord Coord) bool {
	game := gameManager.games[gameID]
	if game == nil {
		return false
	}

	placed := game.PlaceStone(userID, coord)
	return placed
}

func (gameManager *GameManager) Pass(gameID string, userID string) bool {
	game := gameManager.games[gameID]
	if game == nil {
		return false
	}

	game.Pass(userID)
	return true
}

func (gameManager *GameManager) GetOtherPlayer(gameID string, userID string) (*Player, error) {
	game := gameManager.games[gameID]
	if game == nil {
		return &Player{}, errors.New("No other player")
	}

	return game.GetOtherPlayer(userID)
}
