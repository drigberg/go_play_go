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

	if gameBoard.Spaces["white"] == nil {
		t.Error("Expected board.Spaces['white'] to not be nil")
	}

	if gameBoard.Spaces["black"] == nil {
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
}

