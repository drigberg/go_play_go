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

	whiteSpaces := board.ListSpacesForColor(WHITE)
	blackSpaces := board.ListSpacesForColor(BLACK)

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

func TestBoardRemoveStone(t *testing.T) {
	board := NewBoard(9)

	board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	board.removeStone(Coord{X: 0, Y: 0})
	placedAgain := board.PlaceStone(Coord{X: 0, Y: 0}, WHITE)

	if !placedAgain {
		t.Errorf("Should have been able to place stone after removing from that space")
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

	whiteSpaces := board.ListSpacesForColor(WHITE)
	blackSpaces := board.ListSpacesForColor(BLACK)

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

	whiteSpaces := board.ListSpacesForColor(WHITE)
	blackSpaces := board.ListSpacesForColor(BLACK)

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
	whiteSpaces := board.ListSpacesForColor(WHITE)
	blackSpaces := board.ListSpacesForColor(BLACK)
	if len(whiteSpaces) != 3 {
		t.Errorf("Expected 3 white spaces, got %d", len(whiteSpaces))
	}
	if len(blackSpaces) != 6 {
		t.Errorf("Expected 6 black spaces, got %d", len(blackSpaces))
	}

	// place final stone
	board.PlaceStone(Coord{X: 1, Y: 1}, BLACK)

	whiteSpaces = board.ListSpacesForColor(WHITE)
	blackSpaces = board.ListSpacesForColor(BLACK)
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
	whiteSpaces := board.ListSpacesForColor(WHITE)
	blackSpaces := board.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 2 {
		t.Errorf("Expected 2 white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 8 {
		t.Errorf("Expected 8 black spaces, got %d", len(blackSpaces))
	}

	// place final stone
	board.PlaceStone(Coord{X: 2, Y: 1}, BLACK)
	whiteSpaces = board.ListSpacesForColor(WHITE)
	blackSpaces = board.ListSpacesForColor(BLACK)

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
	whiteSpaces := board.ListSpacesForColor(WHITE)
	blackSpaces := board.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 8 {
		t.Errorf("Expected 8 white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 4 {
		t.Errorf("Expected 4 black spaces, got %d", len(blackSpaces))
	}

	// place final stone
	board.PlaceStone(Coord{X: 4, Y: 4}, WHITE)
	whiteSpaces = board.ListSpacesForColor(WHITE)
	blackSpaces = board.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 9 {
		t.Errorf("Expected 9 white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 0 {
		t.Errorf("Expected 0 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardGetFreeSpaces(t *testing.T) {
	board := NewBoard(9)
	freeSpaces := board.getFreeSpaces(board.Spaces)
	if len(freeSpaces) != 81 {
		t.Errorf("Expected 81 free spaces, got %d", len(freeSpaces))
	}

	// place stones and check again
	board.PlaceStone(Coord{X: 2, Y: 4}, WHITE)
	board.PlaceStone(Coord{X: 3, Y: 3}, WHITE)
	board.PlaceStone(Coord{X: 3, Y: 4}, BLACK)
	board.PlaceStone(Coord{X: 4, Y: 3}, BLACK)

	freeSpaces = board.getFreeSpaces(board.Spaces)
	if len(freeSpaces) != 77 {
		t.Errorf("Expected 77 free spaces, got %d", len(freeSpaces))
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

func TestBoardCopySpaces(t *testing.T) {
	board := NewBoard(9)
	board.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	board.PlaceStone(Coord{X: 5, Y: 5}, WHITE)

	spacesCopy := board.copySpaces()
	copyIsEqual := true
	for x := 0; x < len(board.Spaces); x++ {
		for y := 0; y < len(board.Spaces[x]); y++ {
			if board.Spaces[x][y] != spacesCopy[x][y] {
				copyIsEqual = false
			}
		}
	}

	if !copyIsEqual {
		t.Errorf("Expected spacesCopy to equal board.Spaces")
	}

	spacesCopy[0][0] = WHITE

	if board.Spaces[0][0] == WHITE {
		t.Errorf("board.Spaces should not have been mutated when spacesCopy was changed")
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
	freeSpaces := board.getFreeSpaces(spaces)

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

	if scoreData.PointDifference != 17 {
		t.Errorf("Expected white to win by 17 points, got %d", scoreData.PointDifference)
	}
}

func TestBoardGetScoreDataEyes(t *testing.T) {
	board := NewBoard(9)
	for x := 0; x < 9; x++ {
		board.PlaceStone(Coord{X: x, Y: 3}, BLACK)
		board.PlaceStone(Coord{X: x, Y: 4}, WHITE)
	}

	scoreData := board.GetScoreData()

	if scoreData.Winner != WHITE {
		t.Errorf("Expected white to win")
	}

	if scoreData.PointDifference != 17 {
		t.Errorf("Expected white to win by 17 points, got %d", scoreData.PointDifference)
	}
}
