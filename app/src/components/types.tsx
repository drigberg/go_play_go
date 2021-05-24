import {
  array,
  boolean,
  constant,
  either,
  either4,
  either5,
  exact,
  guard,
  null_,
  number,
  string,
} from 'decoders';
import type { Guard } from 'decoders';

// The recommended usage for constants with `decoders` is `constant('someString' as const)`, but the `as` keyword
// is not recognized by this eslint/tsc configuration. The correct configuration would use @typescript-eslint/parser,
// but there is an issue with the latest versions of TypeScript which result in `React` being marked as an unused
// variable.
// - https://github.com/nvie/decoders#constant
// - https://github.com/microsoft/TypeScript/issues/41882
// - https://stackoverflow.com/questions/55807329/why-eslint-throws-no-unused-vars-for-typescript-interface

type IncomingMessage$Remote$GameJoined = {
  name: 'remote/gameJoined';
  data: {
    GameID: string;
  };
};

const incomingMessage$Remote$GameJoinedDecoder = exact({
  name: constant<'remote/gameJoined'>('remote/gameJoined'),
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

type IncomingMessage$Remote$GameInfo = {
  name: 'remote/gameInfo';
  data: GameInfo;
};

const coordDecoder = exact({ X: number, Y: number });
const colorDecoder = either(
  constant<'BLACK'>('BLACK'),
  constant<'WHITE'>('WHITE'),
);

const incomingMessage$GameInfoDecoder = exact({
  name: constant<'remote/gameInfo'>('remote/gameInfo'),
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

type IncomingMessage$Remote$Update = {
  name: 'remote/update';
  data: null;
};

const incomingMessage$UpdateDecoder = exact({
  name: constant<'remote/update'>('remote/update'),
  data: null_,
});

type IncomingMessage$Remote$GameLeft = {
  name: 'remote/gameLeft';
  data: null;
};

const incomingMessage$Remote$GameLeftDecoder = exact({
  name: constant<'remote/gameLeft'>('remote/gameLeft'),
  data: null_,
});

type JoinGameError$Remote$Data = {
  Type: 'remote/joinGame';
};

type RejoinGameError$Remote$Data = {
  Type: 'remote/rejoinGame';
};

type GetGameInfoError$Remote$Data = {
  Type: 'remote/getGameInfo';
};

type Error400$Data = {
  Type: '400';
  Message: string;
};

type IncomingMessage$Error = {
  name: 'error';
  data:
    | GetGameInfoError$Remote$Data
    | RejoinGameError$Remote$Data
    | JoinGameError$Remote$Data
    | Error400$Data;
};

const incomingMessage$ErrorDecoder = exact({
  name: constant<'error'>('error'),
  data: either4(
    exact({
      Type: constant<'400'>('400'),
      Message: string,
    }),
    exact({
      Type: constant<'remote/rejoinGame'>('remote/rejoinGame'),
    }),
    exact({
      Type: constant<'remote/joinGame'>('remote/joinGame'),
    }),
    exact({
      Type: constant<'remote/getGameInfo'>('remote/getGameInfo'),
    }),
  ),
});

type Message =
  | IncomingMessage$Remote$GameJoined
  | IncomingMessage$Remote$GameInfo
  | IncomingMessage$Remote$GameLeft
  | IncomingMessage$Remote$Update
  | IncomingMessage$Error;

const incomingMessageDecoder = either5(
  incomingMessage$GameInfoDecoder,
  incomingMessage$Remote$GameJoinedDecoder,
  incomingMessage$Remote$GameLeftDecoder,
  incomingMessage$UpdateDecoder,
  incomingMessage$ErrorDecoder,
);

export const incomingMessageGuard: Guard<Message> = guard(
  incomingMessageDecoder,
);

export type OutgoingMessage$GetGameInfo$Remote = {
  name: 'remote/getGameInfo';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$RejoinGame$Remote = {
  name: 'remote/rejoinGame';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$JoinGame$Remote = {
  name: 'remote/joinGame';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$CreateGame$Remote = {
  name: 'remote/createGame';
  data: {
    userID: string;
    size: number;
  };
};

export type OutgoingMessage$CreateGame$Local = {
  name: 'local/createGame';
  data: {
    userID: string;
    size: number;
  };
};

export type OutgoingMessage$LeaveGame$Remote = {
  name: 'remote/leaveGame';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$Pass$Remote = {
  name: 'remote/pass';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$PlaceStone$Remote = {
  name: 'remote/placeStone';
  data: {
    gameID: string;
    userID: string;
    coord: {
      X: number;
      Y: number;
    };
  };
};
