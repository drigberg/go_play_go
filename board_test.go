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

func TestBoardGetAllConnectedStonesSingle(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.placeStone(Coord{X: 0, Y: 0}, WHITE)
	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, []Coord{})

	if len(connectedStones) != 1 {
		t.Errorf("Expected 1 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesMultiple(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.placeStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 1, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 2, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 2, Y: 1}, WHITE)

	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, []Coord{})

	if len(connectedStones) != 4 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesBroken(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.placeStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 1, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 2, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 2, Y: 1}, WHITE)
	gameBoard.placeStone(Coord{X: 3, Y: 2}, WHITE)

	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, []Coord{})

	if len(connectedStones) != 4 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesMixed(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.placeStone(Coord{X: 0, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 1, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 1, Y: 1}, WHITE)
	gameBoard.placeStone(Coord{X: 1, Y: 2}, WHITE)
	gameBoard.placeStone(Coord{X: 2, Y: 2}, WHITE)
	gameBoard.placeStone(Coord{X: 2, Y: 1}, WHITE)
	gameBoard.placeStone(Coord{X: 2, Y: 0}, BLACK)
	gameBoard.placeStone(Coord{X: 0, Y: 2}, BLACK)
	gameBoard.placeStone(Coord{X: 5, Y: 5}, WHITE)

	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, []Coord{})

	if len(connectedStones) != 6 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetAllConnectedStonesBlack(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.placeStone(Coord{X: 0, Y: 0}, BLACK)
	gameBoard.placeStone(Coord{X: 1, Y: 0}, BLACK)
	gameBoard.placeStone(Coord{X: 1, Y: 1}, BLACK)
	gameBoard.placeStone(Coord{X: 1, Y: 2}, BLACK)
	gameBoard.placeStone(Coord{X: 2, Y: 2}, BLACK)
	gameBoard.placeStone(Coord{X: 2, Y: 1}, BLACK)
	gameBoard.placeStone(Coord{X: 2, Y: 0}, WHITE)
	gameBoard.placeStone(Coord{X: 0, Y: 2}, WHITE)
	gameBoard.placeStone(Coord{X: 5, Y: 5}, BLACK)

	connectedStones := gameBoard.getAllConnectedStones(Coord{X: 0, Y: 0}, []Coord{})

	if len(connectedStones) != 6 {
		t.Errorf("Expected 4 connected stone, got %d", len(connectedStones))
	}
}

func TestBoardGetNeighboringOpponentStone(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.placeStone(Coord{X: 3, Y: 3}, BLACK)
	gameBoard.placeStone(Coord{X: 3, Y: 4}, WHITE)
	gameBoard.placeStone(Coord{X: 3, Y: 5}, BLACK)

	opponentStones := gameBoard.getNeighboringOpponentStones(Coord{X: 3, Y: 3})
	if len(opponentStones) != 1 {
		t.Errorf("Expected 1 neighboring opponent stone, got %d", len(opponentStones))
	}

	opponentStones = gameBoard.getNeighboringOpponentStones(Coord{X: 3, Y: 4})
	if len(opponentStones) != 2 {
		t.Errorf("Expected 1 neighboring opponent stone, got %d", len(opponentStones))
	}

	opponentStones = gameBoard.getNeighboringOpponentStones(Coord{X: 3, Y: 5})
	if len(opponentStones) != 1 {
		t.Errorf("Expected 1 neighboring opponent stone, got %d", len(opponentStones))
	}
}

func TestBoardGetLiberties(t *testing.T) {
	gameBoard := NewBoard(9)

	gameBoard.placeStone(Coord{X: 0, Y: 0}, BLACK)
	gameBoard.placeStone(Coord{X: 5, Y: 5}, WHITE)
	gameBoard.placeStone(Coord{X: 7, Y: 7}, WHITE)
	gameBoard.placeStone(Coord{X: 7, Y: 8}, WHITE)

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
