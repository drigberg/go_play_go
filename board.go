package main

/**
TODO:
- determine territory
- count points
*/

const (
	FREE  = "FREE"
	WHITE = "WHITE"
	BLACK = "BLACK"
)

type Scores struct {
	BLACK int
	WHITE int
}

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
	PlaceStone(coord Coord, color string) bool
	GetScores() Scores
	GetAvailableSpaces(color string) []Coord
	ListSpacesForColor(color string) []Coord
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

func (board *Board) getStonesToCapture(coord Coord, color string) []Coord {
	opponentColor := WHITE
	if color == WHITE {
		opponentColor = BLACK
	}

	// if no liberties, assert that we are capturing stones
	neighboringOpponentStones := board.getNeighboringOpponentStones(coord, color)
	stonesToCapture := []Coord{}
	for _, stone := range neighboringOpponentStones {
		// Multiple neighboring opponent stones may belong to on group being captured,
		// but we could also be capturing multiple groups.
		alreadyCaptured := false
		for _, captured := range stonesToCapture {
			if stone.X == captured.X && stone.Y == captured.Y {
				alreadyCaptured = true
			}
		}

		if !alreadyCaptured {
			opponentStoneGroup := board.getAllConnectedStones(stone, opponentColor, []Coord{})
			opponentGroupLiberties := 0
			for _, c := range opponentStoneGroup {
				opponentGroupLiberties += board.countLibertiesFuture(c, coord)
			}

			if opponentGroupLiberties == 0 {
				for _, toCapture := range opponentStoneGroup {
					stonesToCapture = append(stonesToCapture, toCapture)
				}
			}
		}
	}

	return stonesToCapture
}

func (board *Board) GetAvailableSpaces(color string) []Coord {
	available := []Coord{}
	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			coord := Coord{X: x, Y: y}
			isAvailable := false
			if board.getSpaceOwnership(coord) == FREE {
				if board.countLiberties(coord) > 0 {
					isAvailable = true
				} else {
					// if no liberties, assert that we are capturing stones
					stonesToCapture := board.getStonesToCapture(coord, color)

					if len(stonesToCapture) > 0 {
						isAvailable = true
					} else {
						// if no liberties and not capturing, assert that connected stones will have at
						// least one remaining liberty
						allConnectedStones := board.getAllConnectedStones(coord, color, []Coord{})
						groupLiberties := 0
						for _, c := range allConnectedStones {
							groupLiberties += board.countLibertiesFuture(c, coord)
						}

						if groupLiberties > 0 {
							isAvailable = true
						}
					}
				}
			}
			if isAvailable {
				available = append(available, coord)
			}
		}
	}
	return available
}

func (board *Board) canPlaceStone(coord Coord, color string) bool {
	available := board.GetAvailableSpaces(color)
	for _, c := range available {
		if c.X == coord.X && c.Y == coord.Y {
			return true
		}
	}

	return false
}

func (board *Board) PlaceStone(coord Coord, color string) bool {
	if !board.canPlaceStone(coord, color) {
		return false
	}

	board.Spaces[coord.X][coord.Y] = color

	stonesToCapture := board.getStonesToCapture(coord, color)
	for _, toCapture := range stonesToCapture {
		board.removeStone(toCapture)
	}
	return true
}

func (board *Board) removeStone(coord Coord) bool {
	if board.getSpaceOwnership(coord) == FREE {
		return false
	}

	board.Spaces[coord.X][coord.Y] = FREE
	return true
}

func (board *Board) ListSpacesForColor(color string) []Coord {
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
	return coord.X >= 0 && coord.X < board.Size && coord.Y >= 0 && coord.Y < board.Size
}

func (board *Board) getNeighborCoords(coord Coord) []Coord {
	unverified := []Coord{Coord{X: coord.X - 1, Y: coord.Y}, Coord{X: coord.X + 1, Y: coord.Y}, Coord{X: coord.X, Y: coord.Y - 1}, Coord{X: coord.X, Y: coord.Y + 1}}
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

func (board *Board) countLibertiesFuture(coord Coord, proposed Coord) int {
	neighborCoords := board.getNeighborCoords(coord)
	liberties := 0
	for _, neighborCoord := range neighborCoords {
		if board.getSpaceOwnership(neighborCoord) == FREE && !(coord.X == proposed.X && coord.Y == proposed.Y) {
			liberties++
		}
	}
	return liberties
}

func (board *Board) getConnectedStones(coord Coord, color string) []Coord {
	neighborCoords := board.getNeighborCoords(coord)
	connected := []Coord{}
	for _, neighborCoord := range neighborCoords {
		if board.getSpaceOwnership(neighborCoord) == color {
			connected = append(connected, neighborCoord)
		}
	}

	return connected
}

func (board *Board) getNeighboringOpponentStones(coord Coord, color string) []Coord {
	neighborCoords := board.getNeighborCoords(coord)
	opponentStones := []Coord{}
	for _, neighborCoord := range neighborCoords {
		neighborColor := board.getSpaceOwnership(neighborCoord)
		if neighborColor != FREE && neighborColor != color {
			opponentStones = append(opponentStones, neighborCoord)
		}
	}

	return opponentStones
}

// we require a color argument so that we can see connected stones for proposed placements: if
// we were just checking color by existing ownership, we'd get no results
func (board *Board) getAllConnectedStones(coord Coord, color string, connected []Coord) []Coord {
	connected = append(connected, coord)
	neighbors := board.getConnectedStones(coord, color)
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
				connected = board.getAllConnectedStones(n, color, connected)
			}
		}
	}

	return connected
}

func (board *Board) GetScores() Scores {
	// TODO: tally points
	// whiteScore := 0
	// blackScore: := 0
	return Scores{
		BLACK: 0,
		WHITE: 0,
	}
}
