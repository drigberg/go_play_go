package main

import (
	"sort"
)

const (
	FREE  = "FREE"
	WHITE = "WHITE"
	BLACK = "BLACK"
)

type Spaces struct {
	BLACK map[string]bool
	WHITE map[string]bool
}

// Board contains the state of the game board
type Board struct {
	Size   int
	Spaces Spaces
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
func NewBoard(size int) Board {
	spaces := Spaces{BLACK: make(map[string]bool), WHITE: make(map[string]bool)}
	return Board{
		Size:   size,
		Spaces: spaces,
	}
}

func (board *Board) canPlaceStone(move Coord) bool {
	spotStr := move.String()

	if board.Spaces.WHITE[spotStr] {
		return false
	}

	if board.Spaces.BLACK[spotStr] {
		return false
	}

	// TODO: cannot place in eyes unless taking capturing

	return true
}

func (board *Board) placeStone(move Coord, color string) bool {
	if !board.canPlaceStone(move) {
		return false
	}

	spotStr := move.String()

	if color == BLACK {
		board.Spaces.BLACK[spotStr] = true
	} else {
		board.Spaces.WHITE[spotStr] = true
	}
	return true
}

func (board *Board) listSpacesForColor(color string) []string {
	unsorted := board.Spaces.BLACK
	if color == WHITE {
		unsorted = board.Spaces.WHITE
	}

	spaces := []string{}
	for space, _ := range unsorted {
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
