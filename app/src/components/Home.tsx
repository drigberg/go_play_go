import React, { useEffect, useState } from 'react';
import Lobby from './Lobby';
import GameLocal from './GameLocal';
import GameRemote from './GameRemote';
import { incomingMessageGuard } from './types';
import { nanoid } from 'nanoid';
import type {
  GameInfo$Local,
  GameInfo$Remote,
  OutgoingMessage$GetGameInfo$Local,
  OutgoingMessage$RejoinGame$Local,
  OutgoingMessage$GetGameInfo$Remote,
  OutgoingMessage$RejoinGame$Remote,
} from './types';

const userIdKey = 'goPlayGo.userId';
const gameIdKey = 'goPlayGo.gameId';
const gameTypeKey = 'goPlayGo.gameType';

function Home(): JSX.Element {
  const [backoffSeconds, setBackoffSeconds] = useState<number>(0);
  const [connected, setConnected] = useState<boolean>(false);
  const [connectionCounter, setConnectionCounter] = useState<number>(0);
  const [error, setError] = useState<string | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [gameId, setGameId] = useState<string | null>(null);
  const [gameInfoRemote, setGameInfoRemote] =
    useState<GameInfo$Remote | null>(null);
  const [gameInfoLocal, setGameInfoLocal] =
    useState<GameInfo$Local | null>(null);
  const [joinGameId, setJoinGameId] = useState<string | null>(null);
  const [socket, setSocket] = useState<WebSocket | null>(null);

  if (userId === null) {
    const userIdStored = localStorage.getItem(userIdKey);
    if (userIdStored) {
      setUserId(userIdStored);
    } else {
      const newUserId = nanoid(5);
      localStorage.setItem(userIdKey, newUserId);
      setUserId(newUserId);
    }
  }

  if (gameId === null) {
    const gameIdStored = localStorage.getItem(gameIdKey);
    if (gameIdStored) {
      setGameId(gameIdStored);
    }
  }

  function rejoinGameLocal() {
    if (userId === null || gameId === null || socket === null) {
      return;
    }
    const message: OutgoingMessage$RejoinGame$Local = {
      name: 'local/rejoinGame',
      data: {
        userID: userId,
        gameID: gameId,
      },
    };
    socket.send(JSON.stringify(message));
  }

  function rejoinGameRemote() {
    if (userId === null || gameId === null || socket === null) {
      return;
    }
    const message: OutgoingMessage$RejoinGame$Remote = {
      name: 'remote/rejoinGame',
      data: {
        userID: userId,
        gameID: gameId,
      },
    };
    socket.send(JSON.stringify(message));
  }

  function getGameInfoRemote() {
    if (gameId && userId && socket) {
      const message: OutgoingMessage$GetGameInfo$Remote = {
        name: 'remote/getGameInfo',
        data: {
          userID: userId,
          gameID: gameId,
        },
      };
      socket.send(JSON.stringify(message));
    } else {
      setError("Can't get game info: try refreshing the page");
    }
  }

  function getGameInfoLocal() {
    if (gameId && userId && socket) {
      const message: OutgoingMessage$GetGameInfo$Local = {
        name: 'local/getGameInfo',
        data: {
          userID: userId,
          gameID: gameId,
        },
      };
      socket.send(JSON.stringify(message));
    } else {
      setError("Can't get game info: try refreshing the page");
    }
  }

  useEffect(() => {
    if (!socket) {
      console.log(`Waiting ${backoffSeconds} seconds before connecting...`);
      setTimeout(() => {
        console.log('Connecting...');
        // TODO: reconnect with backoff if disconnected
        const HOST =
          process.env.NODE_ENV === 'production'
            ? location.origin.replace(/^http/, 'ws')
            : 'ws://localhost:3001';
        const s = new WebSocket(`${HOST}/socket`);
        setSocket(s);
      }, backoffSeconds * 1000);
    }
  }, [socket]);

  useEffect(() => {
    // Reset callbacks whenever gameId or userId are updated, so that we never reference
    // stale values
    if (socket) {
      socket.onopen = () => {
        console.log('Connected!');
        setBackoffSeconds(0);
        setConnected(true);
        setError(null);
        // Ensure that the server is aware of the new socket connection: we might have refreshed
        const gameType = localStorage.getItem(gameTypeKey);
        if (gameType === 'LOCAL') {
          rejoinGameLocal();
        } else if (gameType === 'REMOTE') {
          rejoinGameRemote();
        }
      };

      socket.onclose = () => {
        let newBackoffSeconds = backoffSeconds + 1;
        if (newBackoffSeconds > 5) {
          newBackoffSeconds = 5;
        }
        setBackoffSeconds(newBackoffSeconds);
        console.log('Disconnected!');
        setSocket(null);
        setConnected(false);
      };

      // we don't have a `socket.onerror` callback because the `event` information is not helpful

      socket.onmessage = (event) => {
        const message = incomingMessageGuard(JSON.parse(event.data));
        console.log(message);
        switch (message.name) {
          case 'local/gameInfo':
            setGameInfoLocal(message.data);
            setError(null);
            break;
          case 'local/gameJoined':
            localStorage.setItem(gameIdKey, message.data.GameID.toString());
            localStorage.setItem(gameTypeKey, 'LOCAL');
            setGameId(message.data.GameID);
            setGameInfoRemote(null);
            setError(null);
            break;
          case 'local/update':
            getGameInfoLocal();
            break;
          case 'remote/gameJoined':
            localStorage.setItem(gameIdKey, message.data.GameID.toString());
            localStorage.setItem(gameTypeKey, 'REMOTE');
            setGameId(message.data.GameID);
            setGameInfoLocal(null);
            setError(null);
            break;
          case 'remote/update':
            getGameInfoRemote();
            break;
          case 'remote/gameInfo':
            setGameInfoRemote(message.data);
            setError(null);
            break;
          case 'remote/gameLeft':
            // For remote games, we don't leave until we're sure that the server is aware that we've left, so that
            // the other player is made aware that the game is over.
            localStorage.removeItem(gameIdKey);
            localStorage.removeItem(gameTypeKey);
            setGameId(null);
            setError(null);
            setGameInfoRemote(null);
            break;
          case 'error':
            switch (message.data.Type) {
              case '400':
                setError(message.data.Message);
                break;
              case 'local/rejoinGame':
              case 'local/getGameInfo':
                localStorage.removeItem(gameIdKey);
                setGameId(null);
                setGameInfoLocal(null);
                setError(
                  'Game not found! Either you submitted an invalid game ID, or the server restarted.',
                );
                break;
              case 'remote/rejoinGame':
              case 'remote/joinGame':
              case 'remote/getGameInfo':
                localStorage.removeItem(gameIdKey);
                setGameId(null);
                setGameInfoRemote(null);
                setError(
                  'Game not found! Either you submitted an invalid game ID, or the server restarted.',
                );
                break;
              default:
                throw new Error('Error type not recognized');
            }
            break;
          default:
            throw new Error('Event not recognized!');
        }
      };
    }
  }, [socket, gameId, userId, backoffSeconds]);

  useEffect(() => {
    if (!connected) {
      const intervalId = setInterval(() => {
        setConnectionCounter(connectionCounter + 1);
      }, 500);

      return () => {
        clearInterval(intervalId);
      };
    }
  }, [connected, connectionCounter]);

  function getConnectingDots() {
    let dots = '';
    for (let i = 0; i < (connectionCounter % 3) + 1; i++) {
      dots += '.';
    }
    return dots;
  }

  return (
    <div style={{ textAlign: 'center' }}>
      <h1>GoPlayGo!</h1>
      {error && <h3 style={{ color: 'red' }}>Error: {error}</h3>}
      {!connected && (
        <div>
          <p>Connecting{getConnectingDots()}</p>
        </div>
      )}
      {socket && connected && userId && (
        <div>
          {gameId === null && (
            <Lobby
              userId={userId}
              socket={socket}
              joinGameId={joinGameId}
              setJoinGameId={setJoinGameId}
            />
          )}
          {gameId !== null && gameInfoRemote && (
            <GameRemote
              socket={socket}
              gameId={gameId}
              userId={userId}
              gameInfo={gameInfoRemote}
              getGameInfo={() => getGameInfoRemote()}
            />
          )}
          {gameId !== null && gameInfoLocal && (
            <GameLocal
              socket={socket}
              gameId={gameId}
              userId={userId}
              gameInfo={gameInfoLocal}
              getGameInfo={() => getGameInfoLocal()}
              leaveGame={() => {
                localStorage.removeItem(gameIdKey);
                localStorage.removeItem(gameTypeKey);
                setGameId(null);
                setError(null);
                setGameInfoLocal(null);
              }}
            />
          )}
        </div>
      )}
    </div>
  );
}
export default Home;
