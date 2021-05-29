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

type Player struct {
	UserID       string
	SocketClient *SocketClient
}

type GameRemote struct {
	M             sync.Mutex
	Game          Game
	ID            string
	Players       map[string]*Player
	FirstPlayerID string
	State         string
}

// GameRemoteInterface defines methods a GameRemote should implement
type GameRemoteInterface interface {
	JoinGame(userID string, socketClient *SocketClient) bool
	RejoinGame(userID string, socketClient *SocketClient) bool
	LeaveGame(userID string) bool
	GetInfo(userID string) (GameInfoRemote, error)
	GetOtherPlayer(userID string) (*Player, error)
	GetPlayerColor(userID string) string
	IsTurn(userID string) bool
	Pass()
	PlaceStone(userID string, coord Coord) bool
}

// assert that GameRemote implements GameRemoteInterface
var _ GameRemoteInterface = (*GameRemote)(nil)

// New creates an empty board
func NewGameRemote(gameID string, userID string, size int, socketClient *SocketClient) GameRemote {
	player := Player{
		UserID:       userID,
		SocketClient: socketClient,
	}
	players := make(map[string]*Player)
	players[userID] = &player

	return GameRemote{
		ID:            gameID,
		State:         "WAITING_FOR_OPPONENT",
		FirstPlayerID: userID,
		Players:       players,
		Game:          NewGame(size),
	}
}

type GameInfoRemote struct {
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

func (gameRemote *GameRemote) IsTurn(userID string) bool {
	if userID == gameRemote.FirstPlayerID {
		return gameRemote.Game.Turn%2 == 1
	}
	return gameRemote.Game.Turn%2 == 0
}

func (gameRemote *GameRemote) GetPlayerColor(userID string) string {
	if userID == gameRemote.FirstPlayerID {
		return BLACK
	}
	return WHITE
}

func (gameRemote *GameRemote) GetOtherPlayer(userID string) (*Player, error) {
	for _, player := range gameRemote.Players {
		if player.UserID != userID {
			return player, nil
		}
	}
	return &Player{}, errors.New("No other player")
}

func (gameRemote *GameRemote) PlaceStone(userID string, coord Coord) bool {
	if !gameRemote.IsTurn(userID) {
		return false
	}
	color := gameRemote.GetPlayerColor(userID)
	placed := gameRemote.Game.PlaceStone(color, coord)
	return placed
}

func (gameRemote *GameRemote) Pass() {
	gameRemote.M.Lock()
	defer gameRemote.M.Unlock()

	// If both players pass, the game is over
	gameOver := gameRemote.Game.Pass()
	if gameOver {
		gameRemote.State = "GAME_OVER_PASSED"
	}
}

func (gameRemote *GameRemote) JoinGame(userID string, socketClient *SocketClient) bool {
	if len(gameRemote.Players) >= 2 {
		return false
	}

	gameRemote.M.Lock()
	defer gameRemote.M.Unlock()

	player := Player{
		UserID:       userID,
		SocketClient: socketClient,
	}

	gameRemote.Players[userID] = &player
	gameRemote.State = "PLAYING"
	return true
}

func (gameRemote *GameRemote) LeaveGame(userID string) bool {
	// return false if player is not part of game
	if gameRemote.Players[userID] == nil {
		return false
	}
	gameRemote.M.Lock()
	defer gameRemote.M.Unlock()

	if gameRemote.State != "GAME_OVER_PASSED" {
		gameRemote.State = "GAME_OVER_FORFEIT"
	}

	gameRemote.Players[userID].SocketClient = nil
	return true
}

func (gameRemote *GameRemote) RejoinGame(userID string, socketClient *SocketClient) bool {
	// return false if player is not part of game
	if gameRemote.Players[userID] == nil {
		return false
	}

	gameRemote.M.Lock()
	defer gameRemote.M.Unlock()
	gameRemote.Players[userID].SocketClient = socketClient
	return true
}

// Returns all the information that the client needs for the game state
func (gameRemote *GameRemote) GetInfo(userID string) (GameInfoRemote, error) {
	// return error if player is not part of game
	if gameRemote.Players[userID] == nil {
		return GameInfoRemote{}, errors.New("Cannot get game info")
	}

	color := gameRemote.GetPlayerColor(userID)
	Spaces := Spaces{
		BLACK: gameRemote.Game.Board.ListSpacesForColor(gameRemote.Game.Board.Spaces, BLACK),
		WHITE: gameRemote.Game.Board.ListSpacesForColor(gameRemote.Game.Board.Spaces, WHITE),
	}
	opponentId := "NONE"
	for _, player := range gameRemote.Players {
		if player.UserID != userID {
			opponentId = player.UserID
		}
	}

	playerTurn := gameRemote.IsTurn(userID)

	return GameInfoRemote{
		Size:            gameRemote.Game.Board.Size,
		OpponentID:      opponentId,
		PlayerColor:     color,
		PlayerTurn:      playerTurn,
		State:           gameRemote.State,
		ScoreData:       gameRemote.Game.Board.GetScoreData(),
		AvailableSpaces: gameRemote.Game.Board.GetAvailableSpaces(color),
		Spaces:          Spaces,
		Turn:            gameRemote.Game.Turn,
	}, nil
}
