import React, { useState } from 'react';
import type {
  OutgoingMessage$JoinGame$Remote,
  OutgoingMessage$CreateGame$Local,
  OutgoingMessage$CreateGame$Remote,
} from './types';

type Props = {
  userId: string;
  socket: WebSocket;
  joinGameId: string | null;
  setJoinGameId: (joinGameId: string) => void;
};

function Lobby(props: Props): JSX.Element {
  const [createGameSize, setCreateGameSize] = useState<number | null>(null);

  function createRemoteGame() {
    if (createGameSize === null) {
      return;
    }
    const message: OutgoingMessage$CreateGame$Remote = {
      name: 'remote/createGame',
      data: {
        userID: props.userId,
        size: createGameSize,
      },
    };
    props.socket.send(JSON.stringify(message));
  }

  function createLocalGame() {
    if (createGameSize === null) {
      return;
    }
    const message: OutgoingMessage$CreateGame$Local = {
      name: 'local/createGame',
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
    const message: OutgoingMessage$JoinGame$Remote = {
      name: 'remote/joinGame',
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
      <h2>Create a game...</h2>
      <div style={{ margin: '5px' }}>
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
      </div>
      <div style={{ marginTop: '15px' }}>
        <button disabled={createGameSize === null} onClick={createLocalGame}>
          Create Local Game
        </button>
        <button disabled={createGameSize === null} onClick={createRemoteGame}>
          Create Online Game
        </button>
      </div>
      <h2>...or join using a game id sent by a friend</h2>
      <div style={{ margin: '5px' }}>
        <input
          type="text"
          placeholder="game id"
          onChange={(e) => props.setJoinGameId(e.target.value)}
        />
        <button onClick={joinGame} disabled={props.joinGameId === null}>
          Join Online Game
        </button>
      </div>
    </div>
  );
}
export default Lobby;
