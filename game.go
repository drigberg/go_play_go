package main

import (
	"errors"
	"sync"
)

// State can be one of:
// - WAITING_FOR_OPPONENT
// - PLAYING
// - GAME_OVER_PASSED
// - GAME_OVER_FORFEIT

type Game struct {
	M                sync.Mutex
	ID               string
	Players          map[string]*Player
	Turn             int
	Board            Board
	FirstPlayerID    string
	State            string
	LastPlayerPassed bool
}

type Player struct {
	UserID       string
	SocketClient *SocketClient
}

// GameInterface defines methods a Game should implement
type GameInterface interface {
	GetInfo(userID string) GameInfo
	GetOtherPlayer(userID string) (*Player, error)
	GetPlayerColor(userID string) string
	IsTurn(userID string) bool
	Pass(userID string)
	PlaceStone(userID string, coord Coord) bool
}

// assert that Game implements GameInterface
var _ GameInterface = (*Game)(nil)

type Spaces struct {
	BLACK []Coord
	WHITE []Coord
}

type GameInfo struct {
	Size            int
	Turn            int
	ScoreData       ScoreData
	State           string
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
	game.M.Lock()
	defer game.M.Unlock()

	color := game.GetPlayerColor(userID)
	placed := game.Board.PlaceStone(coord, color)
	if placed {
		game.LastPlayerPassed = false
		game.Turn++
	}
	return placed
}

func (game *Game) Pass(userID string) {
	game.M.Lock()
	defer game.M.Unlock()

	// If both players pass, the game is over
	if game.LastPlayerPassed {
		game.State = "GAME_OVER_PASSED"
	} else {
		game.LastPlayerPassed = true
	}
	game.Turn++
}

func (game *Game) GetOtherPlayer(userID string) (*Player, error) {
	for _, player := range game.Players {
		if player.UserID != userID {
			return player, nil
		}
	}
	return &Player{}, errors.New("No other player")
}

// Returns all the information that the client needs for the game state
func (game *Game) GetInfo(userID string) GameInfo {
	color := game.GetPlayerColor(userID)
	Spaces := Spaces{
		BLACK: game.Board.ListSpacesForColor(game.Board.Spaces, BLACK),
		WHITE: game.Board.ListSpacesForColor(game.Board.Spaces, WHITE),
	}
	opponentId := "NONE"
	for _, player := range game.Players {
		if player.UserID != userID {
			opponentId = player.UserID
		}
	}

	playerTurn := game.IsTurn(userID)

	return GameInfo{
		Size:            game.Board.Size,
		OpponentID:      opponentId,
		PlayerColor:     color,
		PlayerTurn:      playerTurn,
		State:           game.State,
		ScoreData:       game.Board.GetScoreData(),
		AvailableSpaces: game.Board.GetAvailableSpaces(color),
		Spaces:          Spaces,
		Turn:            game.Turn,
	}
}
