package main

import (
	"strconv"
)


type Coord struct {
	X int
	Y int
}

func (coord Coord) String() string {
	return strconv.Itoa(coord.X) + " " + strconv.Itoa(coord.Y)
}

type Message struct {
	Content string
	Sender  string
}

type OpenRoom struct {
	ID     int
	UserID string
}

type Player struct {
	UserID       string
	Color        string
}
