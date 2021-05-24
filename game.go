package main

import (
	"sync"
)

type Game struct {
	M                sync.Mutex
	Turn             int
	Board            Board
	LastPlayerPassed bool
}

type Spaces struct {
	BLACK []Coord
	WHITE []Coord
}

// GameInterface defines methods a Game should implement
type GameInterface interface {
	Pass() bool
	PlaceStone(color string, coord Coord) bool
}

// assert that Game implements GameInterface
var _ GameInterface = (*Game)(nil)

// New creates an empty board
func NewGame(size int) Game {
	return Game{
		Turn:  1,
		Board: NewBoard(size),
	}
}

func (game *Game) PlaceStone(color string, coord Coord) bool {
	game.M.Lock()
	defer game.M.Unlock()

	placed := game.Board.PlaceStone(coord, color)
	if placed {
		game.LastPlayerPassed = false
		game.Turn++
	}
	return placed
}

// Returns true if game is over
func (game *Game) Pass() bool {
	game.M.Lock()
	defer game.M.Unlock()

	game.Turn++

	// If both players pass, the game is over
	if game.LastPlayerPassed {
		return true
	} else {
		game.LastPlayerPassed = true
		return false
	}
}
