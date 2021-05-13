package main

/**
TODO:
- capture stones
- detect eyes
- determine territory
- count points
*/

const (
	FREE  = "FREE"
	WHITE = "WHITE"
	BLACK = "BLACK"
)

type Coord struct {
	X int
	Y int
}

// Board contains the state of the game board
type Board struct {
	Size   int
	Spaces [][]string
}

// BoardInterface defines methods a Board should implement
type BoardInterface interface {
	canPlaceStone(Coord) bool
	getScores() (int, int)
	listSpacesForColor(color string) []Coord
	placeStone(coord Coord, color string) bool
	removeStone(coord Coord) bool
	spaceIsFree(coord Coord) bool
}

// assert that Board implements Interface
var _ BoardInterface = (*Board)(nil)

// New creates an empty board
func NewBoard(size int) Board {
	spaces := make([][]string, size)
	for x := 0; x < size; x++ {
		spaces[x] = make([]string, size)
		for y := 0; y < size; y++ {
			spaces[x][y] = FREE
		}
	}

	return Board{
		Size:   size,
		Spaces: spaces,
	}
}

// returns true if either player has claimed a space
func (board *Board) spaceIsFree(coord Coord) bool {
	return board.Spaces[coord.X][coord.Y] == FREE
}

func (board *Board) canPlaceStone(coord Coord) bool {
	if !board.spaceIsFree(coord) {
		return false
	}

	// TODO: cannot place in eyes unless taking capturing

	return true
}

func (board *Board) placeStone(coord Coord, color string) bool {
	if !board.canPlaceStone(coord) {
		return false
	}

	board.Spaces[coord.X][coord.Y] = color
	return true
}

func (board *Board) removeStone(coord Coord) bool {
	if board.spaceIsFree(coord) {
		return false
	}

	board.Spaces[coord.X][coord.Y] = FREE
	return true
}

func (board *Board) listSpacesForColor(color string) []Coord {
	spaces := []Coord{}
	for x := 0; x < len(board.Spaces); x++ {
		for y := 0; y < len(board.Spaces[x]); y++ {
			if board.Spaces[x][y] == color {
				spaces = append(spaces, Coord{X: x, Y: y})
			}
		}
	}
	return spaces
}

func (board *Board) getScores() (int, int) {
	// TODO: tally points
	// whiteScore := 0
	// blackScore: := 0
	return 0, 0
}
