package main

type Message struct {
	Content string
	Sender  string
}

type OpenRoom struct {
	ID     int
	UserID string
}

type Player struct {
	UserID string
	Color  string
}
