package main

import (
	"sort"
)

var (
	axes = [4][2]int{[2]int{-1, -1}, [2]int{-1, 0}, [2]int{-1, 1}, [2]int{0, -1}}
)

const (
	FREE = "FREE"
)

// Board contains the state of the game board
type Board struct {
	Spaces map[string]map[string]bool
}

// BoardInterface defines methods a Board should implement
type BoardInterface interface {
	getScores() (int, int)
	checkOwnership(int, string, Coord) bool
	isTakenBy(Coord) string
}

// assert that Board implements Interface
var _ BoardInterface = (*Board)(nil)

// New creates an empty board
func NewBoard() Board {
	spaces := make(map[string]map[string]bool)
	spaces["white"] = make(map[string]bool)
	spaces["black"] = make(map[string]bool)
	return Board{
		Spaces: spaces,
	}
}

func (board *Board) listSpaces(color string) []string {
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

func (board *Board) isTakenBy(move Coord) string {
	spotStr := move.String()

	for color := range board.Spaces {
		if board.Spaces[color][spotStr] {
			return color
		}
	}

	return FREE
}

func (board *Board) checkOwnership(gameID int, userID string, move Coord) (bool) {
	return board.isTakenBy(move) != FREE
}


