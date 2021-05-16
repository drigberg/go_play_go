# go_play_go

Everyone's favorite game (Go) implemented with everyone's favorite language (Go)!

Play it live with a friend at: https://go-play-go.herokuapp.com/

## Other details

- Client app runs on React and Typescript
- Points are counted using the Ing method (Great explanation at https://senseis.xmp.net/?IngCounting)
- There is no handling for dead groups during counting! We assume that both players pass after capturing any dead stones, or that they already have decided who won between themselves.

## How to run locally

- Don't do it yet! I need to make some edits since completing the struggle to get this to deploy correctly on Heroku.
- If you really want to, run 1. `go build.`, 2. `npm build`, 3. `ENV=PRODUCTION ./go_play_go`

## Planned features

- Chat
- Tutorial
- CPU mode
- Local multiplayer
- Test coverage for game.go and game_manager.go
