import React, { useEffect, useState } from 'react';
import Board from './Board';
import type {
  Coord,
  GameInfo$Local,
  OutgoingMessage$Pass$Local,
  OutgoingMessage$PlaceStone$Local,
} from './types';

type Props = {
  socket: WebSocket;
  userId: string;
  gameId: string;
  gameInfo: GameInfo$Local | null;
  getGameInfo: () => void;
  leaveGame: () => void;
};

const WAIT_SECONDS = 3;

function GameLocal(props: Props): JSX.Element {
  // If waiting for a server response, don't let player click
  const [waiting, setWaiting] = useState<boolean>(false);

  useEffect(() => {
    props.getGameInfo();
  }, []);

  // If gameInfo is updated, enable player to click again
  useEffect(() => {
    setWaiting(false);
  }, [props.gameInfo]);

  // Wait for up to 3 seconds after placing stones
  useEffect(() => {
    if (waiting) {
      const timeoutId = setTimeout(() => {
        setWaiting(false);
      }, WAIT_SECONDS * 1000);
      return () => {
        // clear timeout if `waiting` is updated somewhere else
        clearTimeout(timeoutId);
      };
    }
  }, [waiting]);

  if (props.gameInfo === null) {
    return <div>Loading...</div>;
  }

  function placeStone(coord: Coord) {
    const message: OutgoingMessage$PlaceStone$Local = {
      name: 'local/placeStone',
      data: {
        userID: props.userId,
        gameID: props.gameId,
        coord,
      },
    };
    props.socket.send(JSON.stringify(message));
    setWaiting(true);
  }

  function pass() {
    const message: OutgoingMessage$Pass$Local = {
      name: 'local/pass',
      data: {
        userID: props.userId,
        gameID: props.gameId,
      },
    };
    setWaiting(true);
    props.socket.send(JSON.stringify(message));
  }

  const gameOver = props.gameInfo.State.startsWith('GAME_OVER');

  return (
    <div>
      <div>
        {gameOver ? (
          <div>
            <h2>Game over!</h2>
            <h3>
              {props.gameInfo.ScoreData.Winner === 'BLACK'
                ? 'Black won'
                : 'White won'}{' '}
              by {props.gameInfo.ScoreData.PointDifference} points!
            </h3>
          </div>
        ) : (
          <p>
            {props.gameInfo.CurrentTurnColor === 'BLACK'
              ? "Black's turn to play"
              : "White's turn to play"}
          </p>
        )}
        {!gameOver && <button onClick={() => pass()}>Pass</button>}
        <Board
          size={props.gameInfo.Size}
          canPlaceStone={props.gameInfo.State === 'PLAYING' && !waiting}
          placeStone={placeStone}
          spaces={props.gameInfo.Spaces}
          availableSpaces={props.gameInfo.AvailableSpaces}
          playerColor={props.gameInfo.CurrentTurnColor}
          lastCoord={props.gameInfo.LastCoord}
        />
        <button onClick={() => props.leaveGame()}>Quit Game</button>
      </div>
    </div>
  );
}
export default GameLocal;
