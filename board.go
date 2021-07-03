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

type StonePlacement struct {
	Color string
	Coord Coord
}

type Mutation struct {
	Add StonePlacement
	Remove []Coord
}

// Board contains the state of the game board
type Board struct {
	Size   int
	Mutations []Mutation
}

func coordsAreEqual(c1 Coord, c2 Coord) bool {
	return c1.X == c2.X && c1.Y == c2.Y
}

func coordIsInList(coord Coord, coords []Coord) bool {
	for _, compare := range coords {
		if coordsAreEqual(coord, compare) {
			return true
		}
	}
	return false
}

// BoardInterface defines methods a Board should implement
type BoardInterface interface {
	PlaceStone(coord Coord, color string) bool
	GetSpaces() [][]string
	GetScoreData() ScoreData
	GetAvailableSpaces(color string) []Coord
	ListSpacesForColor(spaces [][]string, color string) []Coord
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
		Mutations: []Mutation{},
	}
}

func (board *Board) getEmptySpaces() [][]string {
	spaces := make([][]string, board.Size)
	for x := 0; x < board.Size; x++ {
		spaces[x] = make([]string, board.Size)
		for y := 0; y < board.Size; y++ {
			spaces[x][y] = FREE
		}
	}
	return spaces
}

func (board *Board) spacesAreEqual(spaces1 [][]string, spaces2 [][]string) bool {
	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			if spaces1[x][y] != spaces2[x][y] {
				return false
			}
		}
	}
	return true
}

// add and remove stones for turn
func (board *Board) applyMutation(spaces [][]string, mutation Mutation) {
	spaces[mutation.Add.Coord.X][mutation.Add.Coord.Y] = mutation.Add.Color
	for _, toRemove := range mutation.Remove {
		spaces[toRemove.X][toRemove.Y] = FREE
	}
}

// Get state of board from previous turn
func (board *Board) getPreviousSpaces() [][]string {
	spaces := board.getEmptySpaces()
	for turn, mutation := range board.Mutations {
		if turn < len(board.Mutations) - 1 {
			board.applyMutation(spaces, mutation)
		}
	}
	return spaces
}

// get state of stones on board
func (board *Board) GetSpaces() [][]string {
	spaces := board.getEmptySpaces()
	for _, mutation := range board.Mutations {
		board.applyMutation(spaces, mutation)
	}
	return spaces
}

// returns the value of a space
func (board *Board) getSpaceOwnership(coord Coord) string {
	spaces := board.GetSpaces()
	return spaces[coord.X][coord.Y]
}

// determines which stones will be captured by a move
func (board *Board) getStonesToCapture(coord Coord, color string) []Coord {
	opponentColor := WHITE
	if color == WHITE {
		opponentColor = BLACK
	}

	// find all opponent stones bordering the proposed move
	neighboringOpponentStones := board.getNeighboringOpponentStones(coord, color)
	stonesToCapture := []Coord{}
	for _, stone := range neighboringOpponentStones {
		// Multiple neighboring opponent stones may belong to on group being captured,
		// but we could also be capturing multiple groups.
		alreadyCaptured := coordIsInList(stone, stonesToCapture)

		if !alreadyCaptured {
			// finds all stones attached to the neighbor
			opponentStoneGroup := board.getAllConnectedStones(stone, opponentColor, []Coord{})
			opponentGroupLiberties := 0
			for _, c := range opponentStoneGroup {
				opponentGroupLiberties += board.countLibertiesFuture(c, coord)
			}

			// if the group would have no liberties, add to the list of captured stones
			if opponentGroupLiberties == 0 {
				for _, toCapture := range opponentStoneGroup {
					stonesToCapture = append(stonesToCapture, toCapture)
				}
			}
		}
	}

	return stonesToCapture
}

// Returns all valid placements for a player, where stone is on the board and:
// 1) stone will have liberties, or
// 2) capture opponent stones
func (board *Board) GetAvailableSpaces(color string) []Coord {
	available := []Coord{}
	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			coord := Coord{X: x, Y: y}
			isAvailable := false
			if board.getSpaceOwnership(coord) == FREE {
				// if the position has liberties, then it is a valid move
				if board.countLiberties(coord) > 0 {
					isAvailable = true
				} else {
					// if no liberties, assert that we are capturing stones
					stonesToCapture := board.getStonesToCapture(coord, color)
					if len(stonesToCapture) > 0 {
						// ko rule: when capturing, new state cannot equal state from last turn
						previousSpaces := board.getPreviousSpaces()
						spaces := board.GetSpaces()
						mutation := Mutation{
							Add: StonePlacement{
								Coord: coord,
								Color: color,
							},
							Remove: stonesToCapture,
						}
						board.applyMutation(spaces, mutation)
						if !board.spacesAreEqual(spaces, previousSpaces) {
							isAvailable = true
						}
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

// Returns true if the coord is on the board and the stone will either:
// 1) have liberties, or
// 2) capture opponent stones
func (board *Board) canPlaceStone(coord Coord, color string) bool {
	available := board.GetAvailableSpaces(color)
	return coordIsInList(coord, available)
}

// Places a stone on the board, if possible
func (board *Board) PlaceStone(coord Coord, color string) bool {
	if !board.canPlaceStone(coord, color) {
		return false
	}

	stonesToCapture := board.getStonesToCapture(coord, color)
	mutation := Mutation{
		Add: StonePlacement{
			Coord: coord,
			Color: color,
		},
		Remove: stonesToCapture,
	}
	board.Mutations = append(board.Mutations, mutation)

	return true
}

// Lists all spaces belonging to a color (BLACK, WHITE, or FREE)
func (board *Board) ListSpacesForColor(spaces [][]string, color string) []Coord {
	spacesForColor := []Coord{}
	for x := range spaces {
		for y := range spaces[x] {
			if spaces[x][y] == color {
				spacesForColor = append(spacesForColor, Coord{X: x, Y: y})
			}
		}
	}
	return spacesForColor
}

// Returns true if the coord is valid for the board size
func (board *Board) isOnBoard(coord Coord) bool {
	return coord.X >= 0 && coord.X < board.Size && coord.Y >= 0 && coord.Y < board.Size
}

// Returns all valid positions bordering a coordinate
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

// Returns the number of liberties for a stone
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

// Returns the number of liberties a stone will have after the proposed move
func (board *Board) countLibertiesFuture(coord Coord, proposed Coord) int {
	neighborCoords := board.getNeighborCoords(coord)
	liberties := 0
	for _, neighborCoord := range neighborCoords {
		if board.getSpaceOwnership(neighborCoord) == FREE && !(coordsAreEqual(neighborCoord, proposed)) {
			liberties++
		}
	}
	return liberties
}

// Returns a connected group (BLACK, WHITE, or FREE)
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

// Returns all opponent stones bordering a stone
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

// Returns all stones connected to a stone (existing or proposed)
func (board *Board) getAllConnectedStones(coord Coord, color string, connected []Coord) []Coord {
	connected = append(connected, coord)
	// get connected stones for the color of the existing or proposed move
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

// Groups free spaces into chains
func (board *Board) getGroupedFreeSpaces() [][]Coord {
	coords := board.ListSpacesForColor(board.GetSpaces(), FREE)
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

// Returns all colors bordering a coord, ignoring FREE
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

// Returns all colors bordering a group, ignoring FREE
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

// Return groups of free spaces surrounded by each color
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
		} else {
			// no-op: if free spaces touch both colors, then they are unclaimed (dame)
		}
	}

	return Territories{
		BLACK: black,
		WHITE: white,
	}
}

// Place up to 4 white stones in black territory
func (board *Board) placeKomi(territories Territories) (Territories, []Coord) {
	komi := []Coord{}
	for len(territories.BLACK) > 0 && len(territories.BLACK[0]) > 0 {
		// place white stones in four black spaces
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
	return territories, komi
}

type StoneCounts struct {
	BLACK int
	WHITE int
}

func (board *Board) countStones(komi []Coord) StoneCounts {
	numBlackStones := 0
	numWhiteStones := len(komi)

	spaces := board.GetSpaces()
	for x := range spaces {
		for y := range spaces[x] {
			if spaces[x][y] == BLACK {
				numBlackStones++
			} else if spaces[x][y] == WHITE {
				numWhiteStones++
			}
		}
	}

	return StoneCounts{
		BLACK: numBlackStones,
		WHITE: numWhiteStones,
	}
}

// Adds all remaining stones to the board for point-counting
func (board *Board) fillBoard(territories Territories, komi []Coord) ([][]string, StoneCounts) {
	totalsStonesPerPlayer := (board.Size*board.Size - 1) / 2
	stoneCounts := board.countStones(komi)
	remaining := StoneCounts{
		BLACK: totalsStonesPerPlayer - stoneCounts.BLACK,
		WHITE: totalsStonesPerPlayer - stoneCounts.WHITE,
	}
	spaces := board.GetSpaces()

	for _, c := range komi {
		spaces[c.X][c.Y] = WHITE
	}

	for _, group := range territories.BLACK {
		for _, c := range group {
			if remaining.BLACK > 0 {
				spaces[c.X][c.Y] = BLACK
				remaining.BLACK--
			}
		}
	}

	for _, group := range territories.WHITE {
		for _, c := range group {
			if remaining.WHITE > 0 {
				spaces[c.X][c.Y] = WHITE
				remaining.WHITE--
			}
		}
	}

	return spaces, remaining
}

type ScoreData struct {
	Winner          string
	PointDifference float32
}

// Returns the winner's point lead
func (board *Board) getPointDifference(winner string, numFreeSpaces int, remaining StoneCounts) float32 {
	if winner == BLACK {
		return float32(numFreeSpaces-remaining.BLACK+remaining.WHITE) - 0.5
	} else {
		return float32(numFreeSpaces+remaining.BLACK-remaining.WHITE) + 0.5
	}
}

// Determines who is winning / has won the game
func (board *Board) countPoints(spaces [][]string, remaining StoneCounts) ScoreData {
	freeSpaces := board.ListSpacesForColor(spaces, FREE)

	// if we can't determine a winner, fall back to WHITE
	winner := WHITE

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
