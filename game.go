package main

import (
	"sync"
)

type Game struct {
	M                sync.Mutex
	ID               int
	Players          map[string]*Player
	Turn             int
	Board            Board
	FirstPlayerID    string
	IsOver           bool
	LastPlayerPassed bool
}

type Player struct {
	UserID       string
	SocketClient *SocketClient
}

// GameInterface defines methods a Game should implement
type GameInterface interface {
	IsTurn(userID string) bool
}

// assert that Game implements GameInterface
var _ GameInterface = (*Game)(nil)

type Spaces struct {
	BLACK []Coord
	WHITE []Coord
}

type GameInfo struct {
	Turn            int
	Scores          Scores
	IsOver          bool
	PlayerColor     string
	PlayerTurn      bool
	OpponentID      string
	AvailableSpaces []Coord
	Spaces          Spaces
}

func (game *Game) IsTurn(userID string) bool {
	if userID == game.FirstPlayerID {
		return game.Turn%2 == 1
	}
	return game.Turn%2 == 0
}

func (game *Game) GetPlayerColor(userID string) string {
	if userID == game.FirstPlayerID {
		return BLACK
	}
	return WHITE
}

func (game *Game) PlaceStone(userID string, coord Coord) bool {
	color := game.GetPlayerColor(userID)
	placed := game.Board.PlaceStone(coord, color)
	if placed {
		game.LastPlayerPassed = false
	}
	return placed
}

func (game *Game) Pass(userID string) {
	// If both players pass, the game is over
	if game.LastPlayerPassed {
		game.IsOver = true
	} else {
		game.LastPlayerPassed = true
	}
	game.Turn += 1
}

// Returns all the information that the client needs for the game state
func (game *Game) GetInfo(userID string) GameInfo {
	color := game.GetPlayerColor(userID)
	Spaces := Spaces{
		BLACK: game.Board.ListSpacesForColor(BLACK),
		WHITE: game.Board.ListSpacesForColor(WHITE),
	}
	var opponentId string
	for _, player := range game.Players {
		if player.UserID != userID {
			opponentId = player.UserID
		}
	}

	playerTurn := game.IsTurn(userID)

	// TODO: end game if no possible moves left

	return GameInfo{
		OpponentID:      opponentId,
		PlayerColor:     color,
		PlayerTurn:      playerTurn,
		IsOver:          game.IsOver,
		Scores:          game.Board.GetScores(),
		AvailableSpaces: game.Board.GetAvailableSpaces(color),
		Spaces:          Spaces,
		Turn:            game.Turn,
	}
}

func (game *Game) Message(userID string, message string) bool {
	// TODO: send message
	return true
}
