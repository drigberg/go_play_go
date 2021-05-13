package main

import (
	"testing"
)

func TestBoardNew(t *testing.T) {
	gameBoard := NewBoard()
	if gameBoard.Spaces == nil {
		t.Error("Expected board.Spaces to not be nil")
	}
	keys := make([]string, 0, len(gameBoard.Spaces))
	for k := range gameBoard.Spaces {
		keys = append(keys, k)
	}
	if len(keys) != 2 {
		t.Errorf("Expected board.Spaces to have 2 keys, found %d", len(keys))
	}

	if gameBoard.Spaces[WHITE] == nil {
		t.Error("Expected board.Spaces['white'] to not be nil")
	}

	if gameBoard.Spaces[BLACK] == nil {
		t.Error("Expected board.Spaces['black'] to not be nil")
	}
}

func TestBoardGetScoreEmpty(t *testing.T) {
	gameBoard := NewBoard()
	whiteScore, blackScore := gameBoard.getScores()
	if whiteScore != 0 {
		t.Errorf("Expected white to have score 0, got %d", whiteScore)
	}
	if blackScore != 0 {
		t.Errorf("Expected black to have score 0, got %d", blackScore)
	}

	whiteSpaces := gameBoard.listSpaces(WHITE)
	blackSpaces := gameBoard.listSpaces(BLACK)

	if len(whiteSpaces) != 0 {
		t.Errorf("Expected no white spaces, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 0 {
		t.Errorf("Expected no black spaces, got %d", len(blackSpaces))
	}
}

func TestBoardPlaceStone(t *testing.T) {
	gameBoard := NewBoard()

	gameBoard.placeStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 1, Y: 0}, BLACK)
	gameBoard.placeStone(Coord{X: 2, Y: 0}, BLACK)

	whiteSpaces := gameBoard.listSpaces(WHITE)
	blackSpaces := gameBoard.listSpaces(BLACK)

	if len(whiteSpaces) != 1 {
		t.Errorf("Expected 1 white space, got %d", len(whiteSpaces))
	}

	if len(blackSpaces) != 2 {
		t.Errorf("Expected 2 black spaces, got %d", len(blackSpaces))
	}
}
