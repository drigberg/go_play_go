package main

import (
	"testing"
)

func TestBoardNew(t *testing.T) {
	gameBoard := NewBoard(9)
	if gameBoard.Size != 9 {
		t.Errorf("Expected board size 9, got %d", gameBoard.Size)
	}
}

func TestBoardGetScoreEmpty(t *testing.T) {
	gameBoard := NewBoard(9)
	scores := gameBoard.GetScores()
	if scores.WHITE != 0 {
		t.Errorf("Expected white to have score 0, got %d", scores.WHITE)
	}
	if scores.BLACK != 0 {
		t.Errorf("Expected black to have score 0, got %d", scores.BLACK)
	}

	whiteSpaces := gameBoard.ListSpacesForColor(WHITE)
	blackSpaces := gameBoard.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 0 {
		t.Errorf("Expected no black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardPlaceStone(t *testing.T) {
	gameBoard := NewBoard(9)

	// Placing stones on empty spaces
	placements := make([]bool, 3)
	placements[0] = gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	placements[1] = gameBoard.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	placements[2] = gameBoard.PlaceStone(Coord{X: 2, Y: 0}, BLACK)

	whiteSpaces := gameBoard.ListSpacesForColor(WHITE)
	blackSpaces := gameBoard.ListSpacesForColor(BLACK)

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
	placedAgain := gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	if placedAgain {
		t.Errorf("Should not have been able to play on same spot twice")
	}
}

func TestBoardPlaceStoneInEyes(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 0, Y: 1}, BLACK)
	gameBoard.PlaceStone(Coord{X: 1, Y: 2}, BLACK)
	gameBoard.PlaceStone(Coord{X: 0, Y: 3}, BLACK)
	gameBoard.PlaceStone(Coord{X: 2, Y: 1}, BLACK)

	// Placing stones on occupied spaces
	placedInEye := gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	if placedInEye {
		t.Errorf("White should not be able to place stone in black corner eye")
	}

	placedInEye = gameBoard.PlaceStone(Coord{X: 0, Y: 2}, WHITE)
	if placedInEye {
		t.Errorf("White should not be able to place stone in black side eye")
	}

	placedInEye = gameBoard.PlaceStone(Coord{X: 1, Y: 1}, WHITE)
	if placedInEye {
		t.Errorf("White should not be able to place stone in black center eye")
	}

	placedInEye = gameBoard.PlaceStone(Coord{X: 0, Y: 2}, BLACK)
	if !placedInEye {
		t.Errorf("Black should be able to play in its own eyes")
	}
}

func TestBoardRemoveStone(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.removeStone(Coord{X: 0, Y: 0})
	placedAgain := gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)

	if !placedAgain {
		t.Errorf("Should have been able to place stone after removing from that space")
	}
}

func TestBoardGetAllConnectedStonesSingle(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, WHITE, []Coord{})

	if len(connectedStones) != 1 {
		t.Errorf("Expected 1 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesMultiple(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 1, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 2, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 2, Y: 1}, WHITE)

	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, WHITE, []Coord{})

	if len(connectedStones) != 4 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesBroken(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 1, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 2, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 2, Y: 1}, WHITE)
	gameBoard.PlaceStone(Coord{X: 3, Y: 2}, WHITE)

	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, WHITE, []Coord{})

	if len(connectedStones) != 4 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesMixed(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 1, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 1, Y: 1}, WHITE)
	gameBoard.PlaceStone(Coord{X: 1, Y: 2}, WHITE)
	gameBoard.PlaceStone(Coord{X: 2, Y: 2}, WHITE)
	gameBoard.PlaceStone(Coord{X: 2, Y: 1}, WHITE)
	gameBoard.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 0, Y: 2}, BLACK)
	gameBoard.PlaceStone(Coord{X: 5, Y: 5}, WHITE)

	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, WHITE, []Coord{})

	if len(connectedStones) != 6 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesBlack(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 0, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 1, Y: 1}, BLACK)
	gameBoard.PlaceStone(Coord{X: 1, Y: 2}, BLACK)
	gameBoard.PlaceStone(Coord{X: 2, Y: 2}, BLACK)
	gameBoard.PlaceStone(Coord{X: 2, Y: 1}, BLACK)
	gameBoard.PlaceStone(Coord{X: 2, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 0, Y: 2}, WHITE)
	gameBoard.PlaceStone(Coord{X: 5, Y: 5}, BLACK)

	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, BLACK, []Coord{})

	if len(connectedStones) != 6 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetNeighboringOpponentStone(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 3, Y: 3}, BLACK)
	gameBoard.PlaceStone(Coord{X: 3, Y: 4}, WHITE)
	gameBoard.PlaceStone(Coord{X: 3, Y: 5}, BLACK)

	opponentStones := gameBoard.getNeighboringOpponentStones(Coord{X: 3, Y: 3}, BLACK)
	if len(opponentStones) != 1 {
		t.Errorf("Expected 1 neighboring opponent stone, got %d", len(opponentStones))
	}

	opponentStones = gameBoard.getNeighboringOpponentStones(Coord{X: 3, Y: 4}, WHITE)
	if len(opponentStones) != 2 {
		t.Errorf("Expected 1 neighboring opponent stone, got %d", len(opponentStones))
	}

	opponentStones = gameBoard.getNeighboringOpponentStones(Coord{X: 3, Y: 5}, BLACK)
	if len(opponentStones) != 1 {
		t.Errorf("Expected 1 neighboring opponent stone, got %d", len(opponentStones))
	}
}

func TestBoardGetLiberties(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 0, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 5, Y: 5}, WHITE)
	gameBoard.PlaceStone(Coord{X: 7, Y: 7}, WHITE)
	gameBoard.PlaceStone(Coord{X: 7, Y: 8}, WHITE)

	l := gameBoard.countLiberties(Coord{X: 0, Y: 0})
	if l != 2 {
		t.Errorf("Expected 2 liberties, got %d", l)
	}

	l = gameBoard.countLiberties(Coord{X: 5, Y: 5})
	if l != 4 {
		t.Errorf("Expected 4 liberties, got %d", l)
	}

	l = gameBoard.countLiberties(Coord{X: 7, Y: 7})
	if l != 3 {
		t.Errorf("Expected 3 liberties, got %d", l)
	}
}

func TestBoardCaptureSingleCorner(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 0, Y: 1}, BLACK)
	gameBoard.PlaceStone(Coord{X: 1, Y: 0}, BLACK)

	whiteSpaces := gameBoard.ListSpacesForColor(WHITE)
	blackSpaces := gameBoard.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 2 {
		t.Errorf("Expected 2 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardCaptureGroupCorner(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.PlaceStone(Coord{X: 0, Y: 1}, WHITE)
	gameBoard.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 1, Y: 1}, BLACK)
	gameBoard.PlaceStone(Coord{X: 0, Y: 2}, BLACK)

	whiteSpaces := gameBoard.ListSpacesForColor(WHITE)
	blackSpaces := gameBoard.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 3 {
		t.Errorf("Expected 3 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardCaptureGroupCenter(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 1, Y: 2}, WHITE)
	gameBoard.PlaceStone(Coord{X: 2, Y: 2}, WHITE)
	gameBoard.PlaceStone(Coord{X: 2, Y: 1}, WHITE)

	gameBoard.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 3, Y: 1}, BLACK)
	gameBoard.PlaceStone(Coord{X: 3, Y: 2}, BLACK)
	gameBoard.PlaceStone(Coord{X: 2, Y: 3}, BLACK)
	gameBoard.PlaceStone(Coord{X: 1, Y: 3}, BLACK)

	// place second-to-last stone
	gameBoard.PlaceStone(Coord{X: 0, Y: 2}, BLACK)
	whiteSpaces := gameBoard.ListSpacesForColor(WHITE)
	blackSpaces := gameBoard.ListSpacesForColor(BLACK)
	if len(whiteSpaces) != 3 {
		t.Errorf("Expected 3 white spaces, got %d", len(whiteSpaces))
	}
	if len(blackSpaces) != 6 {
		t.Errorf("Expected 6 black spaces, got %d", len(blackSpaces))
	}

	// place final stone
	gameBoard.PlaceStone(Coord{X: 1, Y: 1}, BLACK)

	whiteSpaces = gameBoard.ListSpacesForColor(WHITE)
	blackSpaces = gameBoard.ListSpacesForColor(BLACK)
	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}
	if len(blackSpaces) != 7 {
		t.Errorf("Expected 7 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardCaptureMultipleGroups(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 1, Y: 1}, WHITE)
	gameBoard.PlaceStone(Coord{X: 3, Y: 1}, WHITE)

	gameBoard.PlaceStone(Coord{X: 1, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 2, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 3, Y: 0}, BLACK)
	gameBoard.PlaceStone(Coord{X: 4, Y: 1}, BLACK)
	gameBoard.PlaceStone(Coord{X: 3, Y: 2}, BLACK)
	gameBoard.PlaceStone(Coord{X: 2, Y: 2}, BLACK)
	gameBoard.PlaceStone(Coord{X: 1, Y: 2}, BLACK)

	// place second-to-last stone
	gameBoard.PlaceStone(Coord{X: 0, Y: 1}, BLACK)
	whiteSpaces := gameBoard.ListSpacesForColor(WHITE)
	blackSpaces := gameBoard.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 2 {
		t.Errorf("Expected 2 white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 8 {
		t.Errorf("Expected 8 black spaces, got %d", len(blackSpaces))
	}

	// place final stone
	gameBoard.PlaceStone(Coord{X: 2, Y: 1}, BLACK)
	whiteSpaces = gameBoard.ListSpacesForColor(WHITE)
	blackSpaces = gameBoard.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 9 {
		t.Errorf("Expected 9 black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardCaptureDonut(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.PlaceStone(Coord{X: 2, Y: 4}, WHITE)
	gameBoard.PlaceStone(Coord{X: 3, Y: 3}, WHITE)
	gameBoard.PlaceStone(Coord{X: 4, Y: 2}, WHITE)
	gameBoard.PlaceStone(Coord{X: 5, Y: 3}, WHITE)
	gameBoard.PlaceStone(Coord{X: 6, Y: 4}, WHITE)
	gameBoard.PlaceStone(Coord{X: 5, Y: 5}, WHITE)
	gameBoard.PlaceStone(Coord{X: 4, Y: 6}, WHITE)
	gameBoard.PlaceStone(Coord{X: 3, Y: 5}, WHITE)

	gameBoard.PlaceStone(Coord{X: 3, Y: 4}, BLACK)
	gameBoard.PlaceStone(Coord{X: 4, Y: 3}, BLACK)
	gameBoard.PlaceStone(Coord{X: 5, Y: 4}, BLACK)
	gameBoard.PlaceStone(Coord{X: 4, Y: 5}, BLACK)

	// place second-to-last stone
	whiteSpaces := gameBoard.ListSpacesForColor(WHITE)
	blackSpaces := gameBoard.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 8 {
		t.Errorf("Expected 8 white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 4 {
		t.Errorf("Expected 4 black spaces, got %d", len(blackSpaces))
	}

	// place final stone
	gameBoard.PlaceStone(Coord{X: 4, Y: 4}, WHITE)
	whiteSpaces = gameBoard.ListSpacesForColor(WHITE)
	blackSpaces = gameBoard.ListSpacesForColor(BLACK)

	if len(whiteSpaces) != 9 {
		t.Errorf("Expected 9 white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 0 {
		t.Errorf("Expected 0 black spaces, got %d", len(blackSpaces))
	}
}
