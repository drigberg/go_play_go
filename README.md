# go_play_go

Everyone's favorite game (Go) implemented with everyone's favorite language (Go)!

Play it alone or with a friend at: https://go-play-go.herokuapp.com/

## Design details

- Client app runs on React and Typescript
- Board is rendered with a responsive, mobile-friendly svg
- App is configured to run on Heroku

## Gameplay details

- Points are counted using the Ing method (Great explanation at https://senseis.xmp.net/?IngCounting)
- There is no handling for dead groups during counting! We assume that both players pass after capturing any dead stones, or that they already have decided who won between themselves.

## How to run locally

### Development mode

- Shell #1: `go build . && ./go_play_go
- Shell #2: `cd app && npm start`
- Navigate to `http://localhost:3000`

### Production mode

- `npm run build`
- `go build . && ENV=PRODUCTION PORT=3000 ./go_play_go`
- Navigate to `http://localhost:3000`

## Planned features

- Chat
- Tutorial
- CPU mode
- Local multiplayer
- Test coverage for game.go and game_manager.go
