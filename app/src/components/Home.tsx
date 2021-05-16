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
  const [connected, setConnected] = useState<boolean>(false);
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
      name: 'joinGame',
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
    // TODO: reconnect with backoff if disconnected
    const HOST =
      process.env.NODE_ENV === 'production'
        ? location.origin.replace(/^http/, 'ws')
        : 'ws://localhost:3001';
    const s = new WebSocket(`${HOST}/socket`);
    setSocket(s);
  }, []);

  useEffect(() => {
    // Reset callbacks whenever gameId or userId are updated, so that we never reference
    // stale values
    if (socket) {
      socket.onopen = () => {
        console.log('Connected!');
        setConnected(true);
        setError(null);
        // Ensure that the server is aware of the new socket connection: we might have refreshed
        rejoinGame();
      };

      socket.onclose = () => {
        console.log('Disconnected!');
        setConnected(false);
      };

      socket.onerror = (event) => {
        console.error('Error:', event);
        setError('Something went wrong!');
      };

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
                  "Current game not found on server! The developer didn't want to integrate any sort of database or file store, so this is what you get when the server restarts. Hope you weren't in the middle of a game!!!!",
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
  }, [socket, gameId, userId]);

  return (
    <div style={{ textAlign: 'center' }}>
      <h1>GoPlayGo!</h1>
      {error && <h3 style={{ color: 'red' }}>Error: {error}</h3>}
      {!connected && (
        <div>
          <p>Not connected! Try refreshing.</p>
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
