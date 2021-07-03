package main

import (
	"testing"
)

func TestBoardNew(t *testing.T) {
	board := NewBoard(9)
	if board.Size != 9 {
		t.Errorf("Expected board size 9, got %d", board.Size)
	}
}

func TestBoardPlaceStone(t *testing.T) {
	board := NewBoard(9)

	// Placing stones on empty spaces
	placements := make([]bool, 3)
	placements[0] = board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	placements[1] = board.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	placements[2] = board.PlaceStone(Coord{X: 2, Y: 0}, BLACK)

	spaces := board.GetSpaces()
	whiteSpaces := board.ListSpacesForColor(spaces, WHITE)
	blackSpaces := board.ListSpacesForColor(spaces, BLACK)

	for i := range placements {
		if !placements[i] {
			t.Errorf("Should have been able to place piece on empty space")
		}
	}

	if len(whiteSpaces) != 1 {
		t.Errorf("Expected 1 white space, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 2 {
		t.Errorf("Expected 2 black spaces, got %d", len(blackSpaces))
	}

	// Placing stones on occupied spaces
	placedAgain := board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	if placedAgain {
		t.Errorf("Should not have been able to play on same spot twice")
	}
}

func TestBoardPlaceStoneInEyes(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 0, Y: 1}, BLACK)
	board.PlaceStone(Coord{X: 1, Y: 2}, BLACK)
	board.PlaceStone(Coord{X: 0, Y: 3}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 1}, BLACK)

	// Placing stones on occupied spaces
	placedInEye := board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	if placedInEye {
		t.Errorf("White should not be able to place stone in black corner eye")
	}

	placedInEye = board.PlaceStone(Coord{X: 0, Y: 2}, WHITE)
	if placedInEye {
		t.Errorf("White should not be able to place stone in black side eye")
	}

	placedInEye = board.PlaceStone(Coord{X: 1, Y: 1}, WHITE)
	if placedInEye {
		t.Errorf("White should not be able to place stone in black center eye")
	}

	placedInEye = board.PlaceStone(Coord{X: 0, Y: 2}, BLACK)
	if !placedInEye {
		t.Errorf("Black should be able to play in its own eyes")
	}
}

func TestBoardGetAllConnectedStonesSingle(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	connectedStones := board.getAllConnectedStones(Coord{X: 0, Y: 0}, WHITE, []Coord{})

	if len(connectedStones) != 1 {
		t.Errorf("Expected 1 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesMultiple(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 1, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 1}, WHITE)

	connectedStones := board.getAllConnectedStones(Coord{X: 0, Y: 0}, WHITE, []Coord{})

	if len(connectedStones) != 4 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesBroken(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 1, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 1}, WHITE)
	board.PlaceStone(Coord{X: 3, Y: 2}, WHITE)

	connectedStones := board.getAllConnectedStones(Coord{X: 0, Y: 0}, WHITE, []Coord{})

	if len(connectedStones) != 4 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesMixed(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 1, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 1, Y: 1}, WHITE)
	board.PlaceStone(Coord{X: 1, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 1}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 0, Y: 2}, BLACK)
	board.PlaceStone(Coord{X: 5, Y: 5}, WHITE)

	connectedStones := board.getAllConnectedStones(Coord{X: 0, Y: 0}, WHITE, []Coord{})

	if len(connectedStones) != 6 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesBlack(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 0, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 1, Y: 1}, BLACK)
	board.PlaceStone(Coord{X: 1, Y: 2}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 2}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 1}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 0, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 5, Y: 5}, BLACK)

	connectedStones := board.getAllConnectedStones(Coord{X: 0, Y: 0}, BLACK, []Coord{})

	if len(connectedStones) != 6 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetNeighboringOpponentStone(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 3, Y: 3}, BLACK)
	board.PlaceStone(Coord{X: 3, Y: 4}, WHITE)
	board.PlaceStone(Coord{X: 3, Y: 5}, BLACK)

	opponentStones := board.getNeighboringOpponentStones(Coord{X: 3, Y: 3}, BLACK)
	if len(opponentStones) != 1 {
		t.Errorf("Expected 1 neighboring opponent stone, got %d", len(opponentStones))
	}

	opponentStones = board.getNeighboringOpponentStones(Coord{X: 3, Y: 4}, WHITE)
	if len(opponentStones) != 2 {
		t.Errorf("Expected 1 neighboring opponent stone, got %d", len(opponentStones))
	}

	opponentStones = board.getNeighboringOpponentStones(Coord{X: 3, Y: 5}, BLACK)
	if len(opponentStones) != 1 {
		t.Errorf("Expected 1 neighboring opponent stone, got %d", len(opponentStones))
	}
}

func TestBoardGetLiberties(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 0, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 5, Y: 5}, WHITE)
	board.PlaceStone(Coord{X: 7, Y: 7}, WHITE)
	board.PlaceStone(Coord{X: 7, Y: 8}, WHITE)

	l := board.countLiberties(Coord{X: 0, Y: 0})
	if l != 2 {
		t.Errorf("Expected 2 liberties, got %d", l)
	}

	l = board.countLiberties(Coord{X: 5, Y: 5})
	if l != 4 {
		t.Errorf("Expected 4 liberties, got %d", l)
	}

	l = board.countLiberties(Coord{X: 7, Y: 7})
	if l != 3 {
		t.Errorf("Expected 3 liberties, got %d", l)
	}
}

func TestBoardCaptureSingleCorner(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 0, Y: 1}, BLACK)
	board.PlaceStone(Coord{X: 1, Y: 0}, BLACK)

	spaces := board.GetSpaces()
	whiteSpaces := board.ListSpacesForColor(spaces, WHITE)
	blackSpaces := board.ListSpacesForColor(spaces, BLACK)

	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 2 {
		t.Errorf("Expected 2 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardCaptureGroupCorner(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 0, Y: 1}, WHITE)
	board.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 1, Y: 1}, BLACK)
	board.PlaceStone(Coord{X: 0, Y: 2}, BLACK)

	spaces := board.GetSpaces()
	whiteSpaces := board.ListSpacesForColor(spaces, WHITE)
	blackSpaces := board.ListSpacesForColor(spaces, BLACK)

	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 3 {
		t.Errorf("Expected 3 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardCaptureGroupCenter(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 1, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 1}, WHITE)

	board.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 3, Y: 1}, BLACK)
	board.PlaceStone(Coord{X: 3, Y: 2}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 3}, BLACK)
	board.PlaceStone(Coord{X: 1, Y: 3}, BLACK)

	// place second-to-last stone
	board.PlaceStone(Coord{X: 0, Y: 2}, BLACK)

	spaces := board.GetSpaces()
	whiteSpaces := board.ListSpacesForColor(spaces, WHITE)
	blackSpaces := board.ListSpacesForColor(spaces, BLACK)
	if len(whiteSpaces) != 3 {
		t.Errorf("Expected 3 white spaces, got %d", len(whiteSpaces))
	}
	if len(blackSpaces) != 6 {
		t.Errorf("Expected 6 black spaces, got %d", len(blackSpaces))
	}

	// place final stone
	board.PlaceStone(Coord{X: 1, Y: 1}, BLACK)

	spaces = board.GetSpaces()
	whiteSpaces = board.ListSpacesForColor(spaces, WHITE)
	blackSpaces = board.ListSpacesForColor(spaces, BLACK)
	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}
	if len(blackSpaces) != 7 {
		t.Errorf("Expected 7 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardCaptureMultipleGroups(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 1, Y: 1}, WHITE)
	board.PlaceStone(Coord{X: 3, Y: 1}, WHITE)

	board.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 3, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 4, Y: 1}, BLACK)
	board.PlaceStone(Coord{X: 3, Y: 2}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 2}, BLACK)
	board.PlaceStone(Coord{X: 1, Y: 2}, BLACK)

	// place second-to-last stone
	board.PlaceStone(Coord{X: 0, Y: 1}, BLACK)

	spaces := board.GetSpaces()
	whiteSpaces := board.ListSpacesForColor(spaces, WHITE)
	blackSpaces := board.ListSpacesForColor(spaces, BLACK)

	if len(whiteSpaces) != 2 {
		t.Errorf("Expected 2 white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 8 {
		t.Errorf("Expected 8 black spaces, got %d", len(blackSpaces))
	}

	// place final stone
	board.PlaceStone(Coord{X: 2, Y: 1}, BLACK)

	spaces = board.GetSpaces()
	whiteSpaces = board.ListSpacesForColor(spaces, WHITE)
	blackSpaces = board.ListSpacesForColor(spaces, BLACK)

	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 9 {
		t.Errorf("Expected 9 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardCaptureDonut(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 2, Y: 4}, WHITE)
	board.PlaceStone(Coord{X: 3, Y: 3}, WHITE)
	board.PlaceStone(Coord{X: 4, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 5, Y: 3}, WHITE)
	board.PlaceStone(Coord{X: 6, Y: 4}, WHITE)
	board.PlaceStone(Coord{X: 5, Y: 5}, WHITE)
	board.PlaceStone(Coord{X: 4, Y: 6}, WHITE)
	board.PlaceStone(Coord{X: 3, Y: 5}, WHITE)

	board.PlaceStone(Coord{X: 3, Y: 4}, BLACK)
	board.PlaceStone(Coord{X: 4, Y: 3}, BLACK)
	board.PlaceStone(Coord{X: 5, Y: 4}, BLACK)
	board.PlaceStone(Coord{X: 4, Y: 5}, BLACK)

	// place second-to-last stone
	spaces := board.GetSpaces()
	whiteSpaces := board.ListSpacesForColor(spaces, WHITE)
	blackSpaces := board.ListSpacesForColor(spaces, BLACK)

	if len(whiteSpaces) != 8 {
		t.Errorf("Expected 8 white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 4 {
		t.Errorf("Expected 4 black spaces, got %d", len(blackSpaces))
	}

	// place final stone
	board.PlaceStone(Coord{X: 4, Y: 4}, WHITE)
	spaces = board.GetSpaces()
	whiteSpaces = board.ListSpacesForColor(spaces, WHITE)
	blackSpaces = board.ListSpacesForColor(spaces, BLACK)

	if len(whiteSpaces) != 9 {
		t.Errorf("Expected 9 white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 0 {
		t.Errorf("Expected 0 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardGroupFreeSpaces(t *testing.T) {
	board := NewBoard(9)
	groups := board.getGroupedFreeSpaces()

	if len(groups) != 1 {
		t.Errorf("Expected 1 group of free spaces, got %d", len(groups))
	}

	if len(groups[0]) != 81 {
		t.Errorf("Expected 81 free spaces in group, got %d", len(groups[0]))
	}

	// Seal off a white corner
	board.PlaceStone(Coord{X: 2, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 1}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 1, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 0, Y: 2}, WHITE)

	groups = board.getGroupedFreeSpaces()

	if len(groups) == 2 {
		if len(groups[0]) != 4 {
			t.Errorf("Expected 4 free spaces in group 0, got %d", len(groups[0]))
		}

		if len(groups[1]) != 72 {
			t.Errorf("Expected 72 free spaces in group 1, got %d", len(groups[1]))
		}
	} else {
		t.Errorf("Expected 2 groups of free spaces, got %d", len(groups))
	}
}

func TestBoardGetTerritories(t *testing.T) {
	board := NewBoard(9)
	board.PlaceStone(Coord{X: 2, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 5, Y: 5}, BLACK)

	territories := board.getTerritories()

	if len(territories.BLACK) != 0 {
		t.Errorf("Expected 0 black territories, got %d", len(territories.BLACK))
	}

	if len(territories.WHITE) != 0 {
		t.Errorf("Expected 0 white territories, got %d", len(territories.WHITE))
	}

	// Seal off a white corner
	board.PlaceStone(Coord{X: 2, Y: 0}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 1}, WHITE)
	board.PlaceStone(Coord{X: 2, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 1, Y: 2}, WHITE)
	board.PlaceStone(Coord{X: 0, Y: 2}, WHITE)

	territories = board.getTerritories()

	if len(territories.BLACK) != 0 {
		t.Errorf("Expected 0 black territories, got %d", len(territories.BLACK))
	}

	if len(territories.WHITE) == 1 {
		if len(territories.WHITE[0]) != 4 {
			t.Errorf("Expected 4 free spaces in white territory 0, got %d", len(territories.WHITE[0]))
		}
	} else {
		t.Errorf("Expected 1 white territory, got %d", len(territories.WHITE))
	}
}

func TestBoardPlaceKomi(t *testing.T) {
	board := NewBoard(9)
	board.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 5, Y: 5}, WHITE)

	territories := board.getTerritories()
	territories, komi := board.placeKomi(territories)

	if len(territories.BLACK) != 0 {
		t.Errorf("Expected 0 black territories, got %d", len(territories.BLACK))
	}

	if len(territories.WHITE) != 0 {
		t.Errorf("Expected 0 white territories, got %d", len(territories.WHITE))
	}

	if len(komi) != 0 {
		t.Errorf("Expected 0 komi, got %d", len(komi))
	}

	// seal off a black corner
	board.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 1}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 2}, BLACK)
	board.PlaceStone(Coord{X: 2, Y: 3}, BLACK)
	board.PlaceStone(Coord{X: 1, Y: 3}, BLACK)
	board.PlaceStone(Coord{X: 0, Y: 3}, BLACK)

	territories = board.getTerritories()
	territories, komi = board.placeKomi(territories)

	if len(territories.WHITE) != 0 {
		t.Errorf("Expected 0 white territories, got %d", len(territories.BLACK))
	}

	if len(territories.BLACK) == 1 {
		if len(territories.BLACK[0]) != 2 {
			t.Errorf("Expected 2 free spaces in black territory 0, got %d", len(territories.BLACK[0]))
		}
	} else {
		t.Errorf("Expected 1 black territory, got %d", len(territories.WHITE))
	}

	if len(komi) != 4 {
		t.Errorf("Expected 4 komi, got %d", len(komi))
	}
}

func TestBoardCountStones(t *testing.T) {
	board := NewBoard(9)
	board.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 5, Y: 5}, WHITE)
	komi := []Coord{}
	stoneCounts := board.countStones(komi)

	if stoneCounts.BLACK != 1 {
		t.Errorf("Expected 1 black stone, got %d", stoneCounts.BLACK)
	}

	if stoneCounts.WHITE != 1 {
		t.Errorf("Expected 1 white stone, got %d", stoneCounts.WHITE)
	}

	board.PlaceStone(Coord{X: 3, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 6, Y: 5}, WHITE)
	komi = []Coord{Coord{X: 0, Y: 0}, Coord{X: 1, Y: 0}}
	stoneCounts = board.countStones(komi)

	if stoneCounts.BLACK != 2 {
		t.Errorf("Expected 2 black stones, got %d", stoneCounts.BLACK)
	}

	if stoneCounts.WHITE != 4 {
		t.Errorf("Expected 4 white stones, got %d", stoneCounts.WHITE)
	}
}

func TestBoardFillBoard(t *testing.T) {
	board := NewBoard(9)
	for x := 0; x < 9; x++ {
		board.PlaceStone(Coord{X: x, Y: 3}, BLACK)
		board.PlaceStone(Coord{X: x, Y: 4}, WHITE)
	}

	territories := board.getTerritories()

	if len(territories.WHITE) != 1 {
		t.Errorf("Expected 1 white territory, got %d", len(territories.WHITE))
	}

	if len(territories.BLACK) != 1 {
		t.Errorf("Expected 1 black territory, got %d", len(territories.BLACK))
	}

	territories, komi := board.placeKomi(territories)

	if len(territories.WHITE) != 1 {
		t.Errorf("Expected 1 white territory, got %d", len(territories.WHITE))
	}

	if len(territories.BLACK) != 1 {
		t.Errorf("Expected 1 black territory, got %d", len(territories.BLACK))
	}

	if len(komi) != 4 {
		t.Errorf("Expected 4 komi, got %d", len(komi))
	}

	spaces, remaining := board.fillBoard(territories, komi)
	freeSpaces := board.ListSpacesForColor(spaces, FREE)

	if len(freeSpaces) != 9 {
		t.Errorf("Expected 9 free spaces, got %d", len(freeSpaces))
	}

	if remaining.BLACK != 8 {
		t.Errorf("Expected 8 black stones, got %d", remaining.BLACK)
	}

	if remaining.WHITE != 0 {
		t.Errorf("Expected 0 white stones, got %d", remaining.WHITE)
	}
}

func TestBoardGetScoreDataBasic(t *testing.T) {
	board := NewBoard(9)
	for x := 0; x < 9; x++ {
		board.PlaceStone(Coord{X: x, Y: 3}, BLACK)
		board.PlaceStone(Coord{X: x, Y: 4}, WHITE)
	}

	scoreData := board.GetScoreData()

	if scoreData.Winner != WHITE {
		t.Errorf("Expected white to win")
	}

	if scoreData.PointDifference != 17.5 {
		t.Errorf("Expected white to win by 17.5 points, got %f", scoreData.PointDifference)
	}
}

func TestBoardGetScoreDataEyes(t *testing.T) {
	board := NewBoard(9)
	blackCoords := []Coord{
		Coord{0, 0},
		Coord{0, 1},
		Coord{0, 2},
		Coord{0, 3},
		Coord{0, 4},
		Coord{0, 5},
		Coord{0, 6},
		Coord{1, 1},
		Coord{1, 3},
		Coord{1, 6},
		Coord{1, 7},
		Coord{1, 8},
		Coord{2, 0},
		Coord{2, 7},
		Coord{3, 0},
		Coord{3, 7},
		Coord{4, 0},
		Coord{4, 1},
		Coord{4, 7},
		Coord{5, 1},
		Coord{5, 7},
		Coord{6, 0},
		Coord{6, 1},
		Coord{6, 2},
		Coord{6, 7},
		Coord{7, 1},
		Coord{7, 2},
		Coord{7, 7},
		Coord{8, 1},
		Coord{8, 2},
		Coord{8, 3},
		Coord{8, 7},
	}

	whiteCoords := []Coord{
		Coord{1, 2},
		Coord{1, 4},
		Coord{1, 5},
		Coord{2, 1},
		Coord{2, 2},
		Coord{2, 3},
		Coord{2, 4},
		Coord{2, 5},
		Coord{2, 6},
		Coord{3, 1},
		Coord{3, 2},
		Coord{3, 4},
		Coord{3, 5},
		Coord{3, 6},
		Coord{4, 2},
		Coord{4, 3},
		Coord{4, 6},
		Coord{5, 2},
		Coord{5, 3},
		Coord{5, 4},
		Coord{5, 5},
		Coord{5, 6},
		Coord{6, 3},
		Coord{6, 6},
		Coord{7, 3},
		Coord{7, 4},
		Coord{7, 5},
		Coord{7, 6},
		Coord{8, 4},
		Coord{8, 6},
	}

	for _, c := range blackCoords {
		board.PlaceStone(c, BLACK)
	}

	for _, c := range whiteCoords {
		board.PlaceStone(c, WHITE)
	}

	scoreData := board.GetScoreData()

	if scoreData.Winner != BLACK {
		t.Errorf("Expected black to win")
	}

	if scoreData.PointDifference != 0.5 {
		t.Errorf("Expected black to win by 0.5 points, got %f", scoreData.PointDifference)
	}
}
