package main

import (
	"sync"
)

// GameManager handles all requests and game states
type GameManager struct {
	M             sync.Mutex
	games         map[int]*Game
	gameIDPointer int
	wg            sync.WaitGroup
}

// GameManagerInterface defines methods a Game should implement
type GameManagerInterface interface {
	CreateGame(userID string, socketClient *SocketClient) int
	JoinGame(gameID int, userID string, socketClient *SocketClient) bool
	RejoinGame(gameID int, userID string, socketClient *SocketClient) bool
	PlaceStone(gameID int, userID string, coord Coord) bool
}

// assert that GameManager implements GameManagerInterface
var _ GameManagerInterface = (*GameManager)(nil)

// NewServer creates a GameManager instance
func NewGameManager() GameManager {
	return GameManager{
		games: make(map[int]*Game),
	}
}

func (gameManager *GameManager) CreateGame(userID string, socketClient *SocketClient) int {
	gameManager.M.Lock()
	defer gameManager.M.Unlock()
	defer func() { gameManager.gameIDPointer++ }()

	player := Player{
		UserID:       userID,
		SocketClient: socketClient,
	}

	players := make(map[string]*Player)

	players[userID] = &player

	gameManager.games[gameManager.gameIDPointer] = &Game{
		ID:      gameManager.gameIDPointer,
		Players: players,
		Turn:    0,
		Board:   NewBoard(9),
	}

	return gameManager.gameIDPointer
}

func (gameManager *GameManager) JoinGame(gameID int, userID string, socketClient *SocketClient) bool {
	game := gameManager.games[gameID]
	if game == nil {
		return false
	}
	game.M.Lock()
	defer game.M.Unlock()

	player := Player{
		UserID:       userID,
		SocketClient: socketClient,
	}

	game.Players[userID] = &player
	return true
}

func (gameManager *GameManager) RejoinGame(gameID int, userID string, socketClient *SocketClient) bool {
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

func (gameManager *GameManager) PlaceStone(gameID int, userID string, coord Coord) bool {
	game := gameManager.games[gameID]
	if game == nil {
		return false
	}

	game.M.Lock()
	defer game.M.Unlock()

	placed := game.PlaceStone(userID, coord)
	return placed
}