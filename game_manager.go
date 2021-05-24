package main

import (
	"errors"
	"math/rand"
	"sync"
)

// GameManager handles all requests and game states
type GameManager struct {
	M           sync.Mutex
	remoteGames map[string]*GameRemote
	localGames  map[string]*GameLocal
}

const idChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

// GameManagerInterface defines methods a Game should implement
type GameManagerInterface interface {
	CreateGameLocal(userID string, size int, socketClient *SocketClient) string
	CreateGameRemote(userID string, size int, socketClient *SocketClient) string
	GetGameType(gameID string) string
	GetGameInfo(gameID string, userID string) (interface{}, error)
	GetOtherPlayer(gameID string, userID string) (*Player, error)
	LeaveGameRemote(gameID string, userID string) bool
	JoinGameRemote(gameID string, userID string, socketClient *SocketClient) bool
	RejoinGame(gameID string, userID string, socketClient *SocketClient) bool
	Pass(gameID string, userID string) bool
	PlaceStone(gameID string, userID string, coord Coord) bool
}

// assert that GameManager implements GameManagerInterface
var _ GameManagerInterface = (*GameManager)(nil)

// NewServer creates a GameManager instance
func NewGameManager() GameManager {
	return GameManager{
		remoteGames: make(map[string]*GameRemote),
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

func (gameManager *GameManager) GetGameType(gameID string) string {
	if gameManager.remoteGames[gameID] != nil {
		return "REMOTE"
	} else if gameManager.localGames[gameID] != nil {
		return "LOCAL"
	}
	return "NONE"
}

func (gameManager *GameManager) CreateGameLocal(userID string, size int, socketClient *SocketClient) string {
	gameManager.M.Lock()
	defer gameManager.M.Unlock()

	gameID := gameManager.createGameId()
	game := NewGameLocal(gameID, userID, size, socketClient)
	gameManager.localGames[gameID] = &game

	return gameID
}

func (gameManager *GameManager) CreateGameRemote(userID string, size int, socketClient *SocketClient) string {
	gameManager.M.Lock()
	defer gameManager.M.Unlock()

	gameID := gameManager.createGameId()
	game := NewGameRemote(gameID, userID, size, socketClient)
	gameManager.remoteGames[gameID] = &game

	return gameID
}

func (gameManager *GameManager) JoinGameRemote(gameID string, userID string, socketClient *SocketClient) bool {
	game := gameManager.remoteGames[gameID]
	if game == nil {
		return false
	}

	joined := game.JoinGame(userID, socketClient)
	return joined
}

func (gameManager *GameManager) LeaveGameRemote(gameID string, userID string) bool {
	game := gameManager.remoteGames[gameID]
	if game == nil {
		return false
	}

	left := game.LeaveGame(userID)
	return left
}

func (gameManager *GameManager) RejoinGame(gameID string, userID string, socketClient *SocketClient) bool {
	gameType := gameManager.GetGameType(gameID)

	if gameType == "REMOTE" {
		game := gameManager.remoteGames[gameID]
		if game == nil {
			return false
		}
		rejoined := game.RejoinGame(userID, socketClient)
		return rejoined
	} else if gameType == "LOCAL" {
		game := gameManager.localGames[gameID]
		if game == nil {
			return false
		}
		rejoined := game.RejoinGame(userID, socketClient)
		return rejoined
	}

	return false
}

func (gameManager *GameManager) GetGameInfo(gameID string, userID string) (interface{}, error) {
	gameType := gameManager.GetGameType(gameID)

	if gameType == "REMOTE" {
		game := gameManager.remoteGames[gameID]
		if game == nil {
			return GameInfoRemote{}, errors.New("Cannot get game info")
		}

		gameInfo, err := game.GetInfo(userID)
		if err != nil {
			return GameInfoRemote{}, err
		}
		return gameInfo, nil
	} else if gameType == "LOCAL" {
		game := gameManager.localGames[gameID]
		if game == nil {
			return GameInfoLocal{}, errors.New("Cannot get game info")
		}

		gameInfo := game.GetInfo()
		return gameInfo, nil
	}

	return GameInfoLocal{}, errors.New("Cannot get game info")
}

func (gameManager *GameManager) PlaceStone(gameID string, userID string, coord Coord) bool {
	game := gameManager.remoteGames[gameID]
	if game == nil {
		return false
	}

	placed := game.PlaceStone(userID, coord)
	return placed
}

func (gameManager *GameManager) Pass(gameID string, userID string) bool {
	game := gameManager.remoteGames[gameID]
	if game == nil {
		return false
	}

	game.Pass()
	return true
}

func (gameManager *GameManager) GetOtherPlayer(gameID string, userID string) (*Player, error) {
	game := gameManager.remoteGames[gameID]
	if game == nil {
		return &Player{}, errors.New("Game not found")
	}

	otherPlayer, err := game.GetOtherPlayer(userID)
	if err != nil {
		return &Player{}, err
	}
	return otherPlayer, nil
}
