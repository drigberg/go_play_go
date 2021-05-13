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

// NewServer creates a GameManager instance
func NewGameManager() GameManager {
	return GameManager{
		games: make(map[int]*Game),
	}
}

func (gameManager *GameManager) createGame(userID string) int {
	gameManager.M.Lock()
	defer gameManager.M.Unlock()
	defer func() { gameManager.gameIDPointer++ }()

	player := Player{
		UserID: userID,
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

func (gameManager *GameManager) OnPlaceStone(gameID int, userID string, coord Coord) bool {
	game := gameManager.games[gameID]
	if game != nil {
		game.M.Lock()
		defer game.M.Unlock()
	}

	placed := game.PlaceStone(userID, coord)
	return placed
}
