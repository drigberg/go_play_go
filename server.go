package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type GameRoom struct {
	M 		      sync.Mutex
	ID            int
	Players       map[string]*Player
	Turn          int
	Board         Board
	FirstPlayerID string
	IsOver        bool
}

// GameManager handles all requests and game states
type GameManager struct {
	M      		sync.Mutex
	games  		map[int]*GameRoom
	gameIDPointer int
	quit 		chan interface{}
	listener 	net.Listener
	wg 			sync.WaitGroup
}

// NewServer creates a GameManager instance
func NewGameManager() GameManager {
	return GameManager{
		games:  make(map[int]*GameRoom),
		quit: make(chan interface {}),
	}
}

func (gameManager *GameManager) createGame(userID string) int {
	gameManager.M.Lock()
	defer gameManager.M.Unlock()
	defer func() { gameManager.gameIDPointer++ }()

	player := Player{
		UserID:       userID,
	}

	players := make(map[string]*Player)

	players[userID] = &player

	gameManager.games[gameManager.gameIDPointer] = &GameRoom{
		ID:      gameManager.gameIDPointer,
		Players: players,
		Turn:    0,
		Board:   NewBoard(),
	}

	return gameManager.gameIDPointer
}

// PlayMove places a piece
func (game *GameRoom) PlayMove(move Coord, color string) {
	moveStr := move.String()
	game.Board.Spaces[color][moveStr] = true
}


// IsTurn turns if it's a user's turn or not
func IsTurn(game *GameRoom, userID string) bool {
	if userID == game.FirstPlayerID {
		return game.Turn%2 == 1
	}
	return game.Turn%2 == 0
}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")
		c.Next()
	}
}


func RunServer(host string, port string) {
	router := gin.New()
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		s.Emit("data", "This is server, whaddap?")
		return nil
	})

	server.OnEvent("/", "health", func(s socketio.Conn, msg string) {
		log.Println("client-health:", msg)
		s.Emit("data", "Yo yo!")
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("Disconnected:", msg)
	})

	go server.Serve()
	defer server.Close()

	router.Use(GinMiddleware("http://localhost:3000"))

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.PUT("/socket.io/*any", gin.WrapH(server))

	if err := router.Run(); err != nil {
		log.Fatal("failed to run app: ", err)
	}
}