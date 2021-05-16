import {
  array,
  boolean,
  constant,
  either,
  either3,
  either4,
  either5,
  exact,
  guard,
  null_,
  number,
  string,
} from 'decoders';
import type { Guard } from 'decoders';

type IncomingMessage$GameJoined = {
  name: 'gameJoined';
  data: {
    GameID: string;
  };
};

const incomingMessage$GameJoinedDecoder = exact({
  name: constant<'gameJoined'>('gameJoined'),
  data: exact({
    GameID: string,
  }),
});

export type Coord = {
  X: number;
  Y: number;
};

export type GameInfo = {
  Size: number;
  Turn: number;
  PlayerTurn: boolean;
  OpponentID: string;
  PlayerColor: 'BLACK' | 'WHITE';
  State:
    | 'WAITING_FOR_OPPONENT'
    | 'PLAYING'
    | 'GAME_OVER_FORFEIT'
    | 'GAME_OVER_PASSED';
  ScoreData: {
    Winner: 'BLACK' | 'WHITE';
    PointDifference: number;
  };
  AvailableSpaces: Array<Coord>;
  Spaces: {
    BLACK: Array<Coord>;
    WHITE: Array<Coord>;
  };
};

type IncomingMessage$GameInfo = {
  name: 'gameInfo';
  data: GameInfo;
};

const coordDecoder = exact({ X: number, Y: number });
const colorDecoder = either(
  constant<'BLACK'>('BLACK'),
  constant<'WHITE'>('WHITE'),
);

const incomingMessage$GameInfoDecoder = exact({
  name: constant<'gameInfo'>('gameInfo'),
  data: exact({
    Size: number,
    Turn: number,
    PlayerTurn: boolean,
    OpponentID: string,
    PlayerColor: colorDecoder,
    State: either4(
      constant<'WAITING_FOR_OPPONENT'>('WAITING_FOR_OPPONENT'),
      constant<'PLAYING'>('PLAYING'),
      constant<'GAME_OVER_FORFEIT'>('GAME_OVER_FORFEIT'),
      constant<'GAME_OVER_PASSED'>('GAME_OVER_PASSED'),
    ),
    ScoreData: exact({
      Winner: colorDecoder,
      PointDifference: number,
    }),
    AvailableSpaces: array(coordDecoder),
    Spaces: exact({
      BLACK: array(coordDecoder),
      WHITE: array(coordDecoder),
    }),
  }),
});

type IncomingMessage$Update = {
  name: 'update';
  data: null;
};

const incomingMessage$UpdateDecoder = exact({
  name: constant<'update'>('update'),
  data: null_,
});

type IncomingMessage$GameLeft = {
  name: 'gameLeft';
  data: null;
};

const incomingMessage$GameLeftDecoder = exact({
  name: constant<'gameLeft'>('gameLeft'),
  data: null_,
});

type JoinGameError$Data = {
  Type: 'joinGame';
};

type GetGameInfoError$Data = {
  Type: 'getGameInfo';
};

type Error400$Data = {
  Type: '400';
  Message: string;
};

type IncomingMessage$Error = {
  name: 'error';
  data: GetGameInfoError$Data | JoinGameError$Data | Error400$Data;
};

const incomingMessage$ErrorDecoder = exact({
  name: constant<'error'>('error'),
  data: either3(
    exact({
      Type: constant<'400'>('400'),
      Message: string,
    }),
    exact({
      Type: constant<'joinGame'>('joinGame'),
    }),
    exact({
      Type: constant<'getGameInfo'>('getGameInfo'),
    }),
  ),
});

type Message =
  | IncomingMessage$GameJoined
  | IncomingMessage$GameInfo
  | IncomingMessage$GameLeft
  | IncomingMessage$Update
  | IncomingMessage$Error;

const incomingMessageDecoder = either5(
  incomingMessage$GameInfoDecoder,
  incomingMessage$GameJoinedDecoder,
  incomingMessage$GameLeftDecoder,
  incomingMessage$UpdateDecoder,
  incomingMessage$ErrorDecoder,
);

export const incomingMessageGuard: Guard<Message> = guard(
  incomingMessageDecoder,
);

export type OutgoingMessage$GetGameInfo = {
  name: 'getGameInfo';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$JoinGame = {
  name: 'joinGame';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$CreateGame = {
  name: 'createGame';
  data: {
    userID: string;
    size: number;
  };
};

export type OutgoingMessage$LeaveGame = {
  name: 'leaveGame';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$Pass = {
  name: 'pass';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$PlaceStone = {
  name: 'placeStone';
  data: {
    gameID: string;
    userID: string;
    coord: {
      X: number;
      Y: number;
    };
  };
};