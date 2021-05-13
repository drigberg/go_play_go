package main

import (
	"sort"
)

var (
	axes = [4][2]int{[2]int{-1, -1}, [2]int{-1, 0}, [2]int{-1, 1}, [2]int{0, -1}}
)

const (
	FREE  = "FREE"
	WHITE = "WHITE"
	BLACK = "BLACK"
)

// Board contains the state of the game board
type Board struct {
	Spaces map[string]map[string]bool
}

// BoardInterface defines methods a Board should implement
type BoardInterface interface {
	canPlaceStone(Coord) bool
	getScores() (int, int)
	listSpacesForColor(color string) []string
	placeStone(move Coord, color string) bool
}

// assert that Board implements Interface
var _ BoardInterface = (*Board)(nil)

// New creates an empty board
func NewBoard() Board {
	spaces := make(map[string]map[string]bool)
	spaces[WHITE] = make(map[string]bool)
	spaces[BLACK] = make(map[string]bool)
	return Board{
		Spaces: spaces,
	}
}

func (board *Board) canPlaceStone(move Coord) bool {
	spotStr := move.String()

	for color := range board.Spaces {
		if board.Spaces[color][spotStr] {
			return false
		}
	}

	// TODO: cannot place in eyes unless taking capturing

	return true
}

func (board *Board) placeStone(move Coord, color string) bool {
	if !board.canPlaceStone(move) {
		return false
	}

	spotStr := move.String()
	board.Spaces[color][spotStr] = true
	return true
}

func (board *Board) listSpacesForColor(color string) []string {
	spaces := []string{}
	for space, _ := range board.Spaces[color] {
		spaces = append(spaces, space)
	}
	sort.Strings(spaces)
	return spaces
}

func (board *Board) getScores() (int, int) {
	// TODO: tally points
	// whiteScore := 0
	// blackScore: := 0
	return 0, 0
}


