import React, { useEffect } from 'react';
import Board from './Board';
import type {
  Coord,
  GameInfo,
  OutgoingMessage$LeaveGame,
  OutgoingMessage$Pass,
  OutgoingMessage$PlaceStone,
} from './types';

type Props = {
  socket: WebSocket;
  userId: string;
  gameId: string;
  gameInfo: GameInfo | null;
  getGameInfo: () => void;
};

function Game(props: Props): JSX.Element {
  useEffect(() => {
    props.getGameInfo();
  }, []);

  function leaveGame() {
    const message: OutgoingMessage$LeaveGame = {
      name: 'leaveGame',
      data: {
        userID: props.userId,
        gameID: props.gameId,
      },
    };
    props.socket.send(JSON.stringify(message));
  }

  if (props.gameInfo === null) {
    return <div>Loading...</div>;
  }

  if (props.gameInfo.OpponentID === 'NONE') {
    return (
      <div>
        <h2>
          You are in game {props.gameId}! Tell a friend so that they can join!
        </h2>
        <Board gameInfo={props.gameInfo} />
        <button onClick={() => leaveGame()}>Leave Game</button>
      </div>
    );
  }

  function placeStone(coord: Coord) {
    const message: OutgoingMessage$PlaceStone = {
      name: 'placeStone',
      data: {
        userID: props.userId,
        gameID: props.gameId,
        coord,
      },
    };
    props.socket.send(JSON.stringify(message));
  }

  function pass() {
    const message: OutgoingMessage$Pass = {
      name: 'pass',
      data: {
        userID: props.userId,
        gameID: props.gameId,
      },
    };
    props.socket.send(JSON.stringify(message));
  }

  const gameOver = props.gameInfo.State.startsWith('GAME_OVER');

  return (
    <div>
      <div>
        {gameOver ? (
          <div>
            <h2>Game over!</h2>
            {props.gameInfo.State === 'GAME_OVER_FORFEIT' && (
              <h3>Opponent left the game.</h3>
            )}
            <h3>
              {props.gameInfo.ScoreData.Winner === props.gameInfo.PlayerColor
                ? 'You won'
                : 'Opponent won'}{' '}
              by {props.gameInfo.ScoreData.PointDifference} points
            </h3>
          </div>
        ) : (
          <p>
            {props.gameInfo.PlayerTurn
              ? 'Your turn!'
              : 'Waiting for opponent to play...'}
          </p>
        )}
        {!gameOver && props.gameInfo.PlayerTurn && (
          <button onClick={() => pass()}>Pass</button>
        )}
        <Board gameInfo={props.gameInfo} placeStone={placeStone} />
        <div>{`Game ID: ${props.gameId}`}</div>
        <button onClick={() => leaveGame()}>
          {gameOver ? 'Leave Game' : 'Forfeit Game'}
        </button>
      </div>
    </div>
  );
}
export default Game;
