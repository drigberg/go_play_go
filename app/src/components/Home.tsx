import React, { useEffect, useState } from 'react';
import Lobby from './Lobby';
import Game from './Game';
import { incomingMessageGuard } from './types';
import { nanoid } from 'nanoid';
import type {
  GameInfo,
  OutgoingMessage$GetGameInfo,
  OutgoingMessage$JoinGame,
} from './types';

const userIdKey = 'goPlayGo.userId';
const gameIdKey = 'goPlayGo.gameIdKey';

function Home(): JSX.Element {
  const [backoffSeconds, setBackoffSeconds] = useState<number>(0);
  const [connected, setConnected] = useState<boolean>(false);
  const [connectionCounter, setConnectionCounter] = useState<number>(0);
  const [error, setError] = useState<string | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [gameId, setGameId] = useState<string | null>(null);
  const [gameInfo, setGameInfo] = useState<GameInfo | null>(null);
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

  function rejoinGame() {
    if (userId === null || gameId === null || socket === null) {
      return;
    }
    const message: OutgoingMessage$JoinGame = {
      name: 'joinGameRemote',
      data: {
        userID: userId,
        gameID: gameId,
      },
    };
    socket.send(JSON.stringify(message));
  }

  function getGameInfo() {
    if (gameId && userId && socket) {
      const message: OutgoingMessage$GetGameInfo = {
        name: 'getGameInfo',
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
        rejoinGame();
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
        switch (message.name) {
          case 'gameJoined':
            localStorage.setItem(gameIdKey, message.data.GameID.toString());
            setGameId(message.data.GameID);
            setError(null);
            break;
          case 'update':
            getGameInfo();
            break;
          case 'gameInfo':
            setGameInfo(message.data);
            setError(null);
            break;
          case 'gameLeft':
            localStorage.removeItem(gameIdKey);
            setGameId(null);
            setError(null);
            break;
          case 'error':
            switch (message.data.Type) {
              case '400':
                setError(message.data.Message);
                break;
              case 'joinGame':
              case 'getGameInfo':
                localStorage.removeItem(gameIdKey);
                setGameId(null);
                setError(
                  'Game not found! Either you typed in an invalid game ID, or the server restarted.',
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
          {gameId === null ? (
            <Lobby
              userId={userId}
              socket={socket}
              joinGameId={joinGameId}
              setJoinGameId={setJoinGameId}
            />
          ) : (
            <Game
              socket={socket}
              gameId={gameId}
              userId={userId}
              gameInfo={gameInfo}
              getGameInfo={() => getGameInfo()}
            ></Game>
          )}
        </div>
      )}
    </div>
  );
}
export default Home;
