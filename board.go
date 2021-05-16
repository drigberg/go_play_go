package main

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

func coordIsInList(coord Coord, coords []Coord) bool {
	for _, compare := range coords {
		if coord.X == compare.X && coord.Y == compare.Y {
			return true
		}
	}
	return false
}

// BoardInterface defines methods a Board should implement
type BoardInterface interface {
	PlaceStone(coord Coord, color string) bool
	GetScoreData() ScoreData
	GetAvailableSpaces(color string) []Coord
	ListSpacesForColor(color string) []Coord
}

// assert that Board implements BoardInterface
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

	neighboringOpponentStones := board.getNeighboringOpponentStones(coord, color)
	stonesToCapture := []Coord{}
	for _, stone := range neighboringOpponentStones {
		// Multiple neighboring opponent stones may belong to on group being captured,
		// but we could also be capturing multiple groups.
		alreadyCaptured := coordIsInList(stone, stonesToCapture)

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
	return coordIsInList(coord, available)
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
		if board.getSpaceOwnership(neighborCoord) == FREE && !(neighborCoord.X == proposed.X && neighborCoord.Y == proposed.Y) {
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
			isNew := !coordIsInList(n, connected)
			if isNew {
				connected = board.getAllConnectedStones(n, color, connected)
			}
		}
	}

	return connected
}

func (board *Board) getFreeSpaces(spaces [][]string) []Coord {
	coords := []Coord{}
	for x := 0; x < len(spaces); x++ {
		for y := 0; y < len(spaces[x]); y++ {
			if spaces[x][y] == FREE {
				coords = append(coords, Coord{X: x, Y: y})
			}
		}
	}

	return coords
}

func (board *Board) getGroupedFreeSpaces() [][]Coord {
	coords := board.getFreeSpaces(board.Spaces)
	grouped := []Coord{}
	groups := [][]Coord{}

	for _, coord := range coords {
		if !coordIsInList(coord, grouped) {
			group := board.getAllConnectedStones(coord, FREE, []Coord{})
			groups = append(groups, group)
			for _, c := range group {
				grouped = append(grouped, c)
			}
		}
	}

	return groups
}

type Territories struct {
	BLACK [][]Coord
	WHITE [][]Coord
}

// only returns BLACK or WHITE: skips FREE
func (board *Board) getAllNeighborColors(coord Coord) []string {
	neighborCoords := board.getNeighborCoords(coord)
	colors := []string{}

	for _, neighborCoord := range neighborCoords {
		color := board.getSpaceOwnership(neighborCoord)
		if color == BLACK || color == WHITE {
			colorInList := false
			for _, c := range colors {
				if c == color {
					colorInList = true
				}
			}
			if !colorInList {
				colors = append(colors, color)
			}
		}
	}

	return colors
}

// only returns BLACK or WHITE: skips FREE
func (board *Board) getAllNeighborColorsForGroup(coords []Coord) []string {
	colors := []string{}

	for _, coord := range coords {
		neighborColors := board.getAllNeighborColors(coord)
		for _, neighborColor := range neighborColors {
			colorInList := false
			for _, c := range colors {
				if c == neighborColor {
					colorInList = true
				}
			}
			if !colorInList {
				colors = append(colors, neighborColor)
			}
		}
	}

	return colors
}

func (board *Board) getTerritories() Territories {
	groups := board.getGroupedFreeSpaces()

	black := [][]Coord{}
	white := [][]Coord{}
	for _, group := range groups {
		neighborColors := board.getAllNeighborColorsForGroup(group)
		if len(neighborColors) == 1 {
			if neighborColors[0] == BLACK {
				black = append(black, group)
			} else {
				white = append(white, group)
			}
		}
	}

	return Territories{
		BLACK: black,
		WHITE: white,
	}
}

func placeSingleKomiInGroup(group []Coord) ([]Coord, Coord) {
	komi := group[0]
	group = group[1:]
	return group, komi
}

func (board *Board) placeKomi(territories Territories) (Territories, []Coord) {
	komi := []Coord{}
	for len(territories.BLACK) > 0 {
		for len(territories.BLACK[0]) > 0 {
			komi = append(komi, territories.BLACK[0][0])
			territories.BLACK[0] = territories.BLACK[0][1:]
			if len(territories.BLACK[0]) == 0 {
				// remove group if empty
				territories.BLACK = territories.BLACK[1:]
			}
			if len(komi) == 4 {
				return territories, komi
			}
		}
	}
	return territories, komi
}

type StoneCounts struct {
	BLACK int
	WHITE int
}

func (board *Board) countStones(komi []Coord) StoneCounts {
	numBlackStones := 0
	numWhiteStones := len(komi)

	for x := 0; x < len(board.Spaces); x++ {
		for y := 0; y < len(board.Spaces[x]); y++ {
			if board.Spaces[x][y] == BLACK {
				numBlackStones++
			} else if board.Spaces[x][y] == WHITE {
				numWhiteStones++
			}
		}
	}

	return StoneCounts{
		BLACK: numBlackStones,
		WHITE: numWhiteStones,
	}
}

func (board *Board) copySpaces() [][]string {
	spacesCopy := make([][]string, board.Size)
	for x := 0; x < board.Size; x++ {
		spacesCopy[x] = make([]string, board.Size)
		copy(spacesCopy[x], board.Spaces[x])
	}
	return spacesCopy
}

func (board *Board) fillBoard(territories Territories, komi []Coord) ([][]string, StoneCounts) {
	totalsStonesPerPlayer := (board.Size*board.Size - 1) / 2
	stoneCounts := board.countStones(komi)
	remaining := StoneCounts{
		BLACK: totalsStonesPerPlayer - stoneCounts.BLACK,
		WHITE: totalsStonesPerPlayer - stoneCounts.WHITE,
	}
	spacesCopy := board.copySpaces()

	for _, c := range komi {
		spacesCopy[c.X][c.Y] = WHITE
	}

	for _, group := range territories.BLACK {
		for _, c := range group {
			if remaining.BLACK > 0 {
				spacesCopy[c.X][c.Y] = BLACK
				remaining.BLACK--
			}
		}
	}

	for _, group := range territories.WHITE {
		for _, c := range group {
			if remaining.WHITE > 0 {
				spacesCopy[c.X][c.Y] = WHITE
				remaining.WHITE--
			}
		}
	}

	return spacesCopy, remaining
}

type ScoreData struct {
	Winner          string
	PointDifference float32
}

func (board *Board) getPointDifference(winner string, numFreeSpaces int, remaining StoneCounts) float32 {
	if winner == BLACK {
		return float32(numFreeSpaces-remaining.BLACK+remaining.WHITE) - 0.5
	} else {
		return float32(numFreeSpaces+remaining.BLACK-remaining.WHITE) + 0.5
	}
}
func (board *Board) countPoints(spaces [][]string, remaining StoneCounts) ScoreData {
	freeSpaces := board.getFreeSpaces(spaces)
	var winner string

	colors := board.getAllNeighborColorsForGroup(freeSpaces)
	if len(colors) == 1 {
		// the winner is whoever claims the final free space
		winner = colors[0]
	} else {
		// the game wasn't complete, so we'll do our best?
		if remaining.BLACK < remaining.WHITE {
			winner = BLACK
		} else if remaining.BLACK > remaining.WHITE {
			winner = WHITE
		}
	}

	pointDifference := board.getPointDifference(winner, len(freeSpaces), remaining)
	return ScoreData{
		Winner:          winner,
		PointDifference: pointDifference,
	}
}

// GetScore() tallies points using the Ing method, with 4 komi placed in black territory
func (board *Board) GetScoreData() ScoreData {
	// First we find all the free spaces surrounded by each placer
	territories := board.getTerritories()
	// Then we place four white stones in black territory
	territories, komi := board.placeKomi(territories)
	// Using the remaining stones belonging to each player, we fill the claimed territory
	spaces, remaining := board.fillBoard(territories, komi)
	// The winner is whoever holds the final free space
	scoreData := board.countPoints(spaces, remaining)
	return scoreData
}
