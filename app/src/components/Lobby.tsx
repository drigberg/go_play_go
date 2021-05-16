import React, { useState } from 'react';
import type {
  OutgoingMessage$JoinGame,
  OutgoingMessage$CreateGame,
} from './types';

type Props = {
  userId: string;
  socket: WebSocket;
  joinGameId: string | null;
  setJoinGameId: (joinGameId: string) => void;
};

function Lobby(props: Props): JSX.Element {
  const [createGameSize, setCreateGameSize] = useState<number | null>(null);

  function createGame() {
    if (createGameSize === null) {
      return;
    }
    const message: OutgoingMessage$CreateGame = {
      name: 'createGame',
      data: {
        userID: props.userId,
        size: createGameSize,
      },
    };
    props.socket.send(JSON.stringify(message));
  }

  function joinGame() {
    if (props.joinGameId === null) {
      return;
    }
    const message: OutgoingMessage$JoinGame = {
      name: 'joinGame',
      data: {
        userID: props.userId,
        gameID: props.joinGameId,
      },
    };
    props.socket.send(JSON.stringify(message));
  }

  const selectedStyle = {
    backgroundColor: '#ffc4fb',
  };

  function getSizeButtonStyle(size: number): { backgroundColor?: string } {
    return createGameSize === size ? selectedStyle : {};
  }

  return (
    <div>
      <h2>Create or join a game!</h2>
      <div>
        <button
          style={getSizeButtonStyle(9)}
          onClick={() => setCreateGameSize(9)}
        >
          9x9
        </button>
        <button
          style={getSizeButtonStyle(13)}
          onClick={() => setCreateGameSize(13)}
        >
          13x13
        </button>
        <button
          style={getSizeButtonStyle(19)}
          onClick={() => setCreateGameSize(19)}
        >
          19x19
        </button>
        <button disabled={createGameSize === null} onClick={createGame}>
          Create Game
        </button>
      </div>
      <div>
        <input
          type="text"
          placeholder="gameId"
          onChange={(e) => props.setJoinGameId(e.target.value)}
        />
        <button onClick={joinGame} disabled={props.joinGameId === null}>
          Join Game
        </button>
      </div>
    </div>
  );
}
export default Lobby;
