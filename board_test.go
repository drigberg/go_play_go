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
	whiteScore, blackScore := gameBoard.getScores()
	if whiteScore != 0 {
		t.Errorf("Expected white to have score 0, got %d", whiteScore)
	}
	if blackScore != 0 {
		t.Errorf("Expected black to have score 0, got %d", blackScore)
	}

	whiteSpaces := gameBoard.listSpacesForColor(WHITE)
	blackSpaces := gameBoard.listSpacesForColor(BLACK)

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
	placements[0] = gameBoard.placeStone(Coord{X: 0, Y: 0}, WHITE)
	placements[1] = gameBoard.placeStone(Coord{X: 1, Y: 0}, BLACK)
	placements[2] = gameBoard.placeStone(Coord{X: 2, Y: 0}, BLACK)

	whiteSpaces := gameBoard.listSpacesForColor(WHITE)
	blackSpaces := gameBoard.listSpacesForColor(BLACK)

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
	placedAgain := gameBoard.placeStone(Coord{X: 0, Y: 0}, WHITE)
	if placedAgain {
		t.Errorf("Should not have been able to play on same spot twice")
	}
}

func TestBoardRemoveStone(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.placeStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.removeStone(Coord{X: 0, Y: 0})
	placedAgain := gameBoard.placeStone(Coord{X: 0, Y: 0}, WHITE)

	if !placedAgain {
		t.Errorf("Should have been able to place stone after removing from that space")
	}
}
