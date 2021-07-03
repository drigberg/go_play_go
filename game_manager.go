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

// GameManagerInterface defines methods a Game must implement
type GameManagerInterface interface {
	CreateGameLocal(userID string, size int, socketClient *SocketClient) string
	CreateGameRemote(userID string, size int, socketClient *SocketClient) string
	GetGameInfoLocal(gameID string, userID string) (GameInfoLocal, error)
	GetGameInfoRemote(gameID string, userID string) (GameInfoRemote, error)
	RejoinGameLocal(gameID string, userID string, socketClient *SocketClient) bool
	RejoinGameRemote(gameID string, userID string, socketClient *SocketClient) bool
	PassLocal(gameID string, userID string) bool
	PassRemote(gameID string, userID string) bool
	PlaceStoneLocal(gameID string, userID string, coord Coord) bool
	PlaceStoneRemote(gameID string, userID string, coord Coord) bool
	// remote-only methods
	LeaveGameRemote(gameID string, userID string) bool
	GetOtherPlayerRemote(gameID string, userID string) (*Player, error)
	JoinGameRemote(gameID string, userID string, socketClient *SocketClient) bool
}

// assert that GameManager implements GameManagerInterface
var _ GameManagerInterface = (*GameManager)(nil)

// NewServer creates a GameManager instance
func NewGameManager() GameManager {
	return GameManager{
		localGames:  make(map[string]*GameLocal),
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

func (gameManager *GameManager) RejoinGameLocal(gameID string, userID string, socketClient *SocketClient) bool {
	game := gameManager.localGames[gameID]
	if game == nil || game.UserID != userID {
		return false
	}
	rejoined := game.RejoinGame(userID, socketClient)
	return rejoined
}

func (gameManager *GameManager) RejoinGameRemote(gameID string, userID string, socketClient *SocketClient) bool {
	game := gameManager.remoteGames[gameID]
	if game == nil || game.Players[userID] == nil {
		return false
	}
	rejoined := game.RejoinGame(userID, socketClient)
	return rejoined
}

func (gameManager *GameManager) GetGameInfoLocal(gameID string, userID string) (GameInfoLocal, error) {
	game := gameManager.localGames[gameID]
	if game == nil || game.UserID != userID {
		return GameInfoLocal{}, errors.New("Cannot get game info")
	}

	gameInfo := game.GetInfo()
	return gameInfo, nil
}

func (gameManager *GameManager) GetGameInfoRemote(gameID string, userID string) (GameInfoRemote, error) {
	game := gameManager.remoteGames[gameID]
	if game == nil || game.Players[userID] == nil {
		return GameInfoRemote{}, errors.New("Cannot get game info")
	}

	gameInfo, err := game.GetInfo(userID)
	if err != nil {
		return GameInfoRemote{}, err
	}
	return gameInfo, nil
}

func (gameManager *GameManager) PlaceStoneLocal(gameID string, userID string, coord Coord) bool {
	game := gameManager.localGames[gameID]
	if game == nil || game.UserID != userID {
		return false
	}

	placed := game.PlaceStone(coord)
	return placed
}

func (gameManager *GameManager) PlaceStoneRemote(gameID string, userID string, coord Coord) bool {
	game := gameManager.remoteGames[gameID]
	if game == nil {
		return false
	}

	placed := game.PlaceStone(userID, coord)
	return placed
}

func (gameManager *GameManager) PassLocal(gameID string, userID string) bool {
	game := gameManager.localGames[gameID]
	if game == nil || game.UserID != userID {
		return false
	}

	game.Pass()
	return true
}

func (gameManager *GameManager) PassRemote(gameID string, userID string) bool {
	game := gameManager.remoteGames[gameID]
	if game == nil || game.Players[userID] == nil {
		return false
	}

	game.Pass()
	return true
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
	if game == nil || game.Players[userID] == nil {
		return false
	}

	left := game.LeaveGame(userID)
	return left
}

func (gameManager *GameManager) GetOtherPlayerRemote(gameID string, userID string) (*Player, error) {
	game := gameManager.remoteGames[gameID]
	if game == nil || game.Players[userID] == nil {
		return &Player{}, errors.New("Game not found")
	}

	otherPlayer, err := game.GetOtherPlayer(userID)
	if err != nil {
		return &Player{}, err
	}
	return otherPlayer, nil
}
