package main

/**
TODO:
- detect liberties
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
	countLiberties(Coord) int
	getScores() (int, int)
	listSpacesForColor(color string) []Coord
	placeStone(coord Coord, color string) bool
	removeStone(coord Coord) bool
	getSpaceOwnership(coord Coord) string
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
func (board *Board) getSpaceOwnership(coord Coord) string {
	return board.Spaces[coord.X][coord.Y]
}

func (board *Board) canPlaceStone(coord Coord) bool {
	if board.getSpaceOwnership(coord) != FREE {
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
	if board.getSpaceOwnership(coord) == FREE {
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

func (board *Board) isOnBoard(coord Coord) bool {
	return coord.X >= 0 && coord.X <= board.Size && coord.Y >= 0 && coord.Y <= board.Size
}


func (board *Board) getNeighborCoords(coord Coord) []Coord {
	unverified := []Coord{Coord{X: coord.X - 1, Y: coord.Y},Coord{X: coord.X + 1, Y: coord.Y},Coord{X: coord.X, Y: coord.Y - 1},Coord{X: coord.X, Y: coord.Y + 1}}
	neighborCoords := []Coord{}
	for _, c := range unverified {
		if board.isOnBoard(c) {
			neighborCoords = append(neighborCoords, c)
		}
	}
	return neighborCoords
}

func (board *Board) countLiberties(coord Coord) int {
	neighborCoords := board.getNeighborCoords(coord)
	liberties := 0
	for _, neighborCoord := range neighborCoords {
		if board.getSpaceOwnership(neighborCoord) == FREE {
			liberties++
		}
	}
	return liberties
}

func (board *Board) getConnectedStones(coord Coord) ([]Coord, string) {
	color := board.getSpaceOwnership(coord)
	if color == FREE {
		return nil, FREE
	}

	neighborCoords := board.getNeighborCoords(coord)
	connected := []Coord{}
	for _, neighborCoord := range neighborCoords {
		if board.getSpaceOwnership(neighborCoord) == color {
			connected = append(connected, neighborCoord)
		}
	}

	return connected, FREE
}

func (board *Board) getAllConnectedStones(coord Coord, connected []Coord) ([]Coord, string) {
	color := board.getSpaceOwnership(coord)
	if color == FREE {
		return nil, FREE
	}

	connected = append(connected, coord)
	neighbors, _ := board.getConnectedStones(coord)
	if neighbors != nil {
		for _, n := range neighbors {
			// TODO: better duplicate checking (if c in list?)
			isNew := true
			for _, c := range connected {
				if c.X == n.X && c.Y == n.Y {
					isNew = false
				}
			}
			if isNew {
				connected, _ = board.getAllConnectedStones(n, connected)
			}
		}
	}

	return connected, color
}

func (board *Board) getScores() (int, int) {
	// TODO: tally points
	// whiteScore := 0
	// blackScore: := 0
	return 0, 0
}
