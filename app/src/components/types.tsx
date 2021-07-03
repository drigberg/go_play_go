import {
  array,
  boolean,
  constant,
  either,
  either4,
  either6,
  either9,
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

type IncomingMessage$Local$GameJoined = {
  name: 'local/gameJoined';
  data: {
    GameID: string;
  };
};

const incomingMessage$Local$GameJoinedDecoder = exact({
  name: constant<'local/gameJoined'>('local/gameJoined'),
  data: exact({
    GameID: string,
  }),
});

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

type Color = 'BLACK' | 'WHITE';

type ScoreData = {
  Winner: Color;
  PointDifference: number;
};

export type Spaces = {
  BLACK: Array<Coord>;
  WHITE: Array<Coord>;
};

const coordDecoder = exact({ X: number, Y: number });
const colorDecoder = either(
  constant<'BLACK'>('BLACK'),
  constant<'WHITE'>('WHITE'),
);

const scoreDataDecoder = exact({
  Winner: colorDecoder,
  PointDifference: number,
});

const spacesDecoder = exact({
  BLACK: array(coordDecoder),
  WHITE: array(coordDecoder),
});

export type GameInfo$Local = {
  Size: number;
  Turn: number;
  ScoreData: ScoreData;
  State: 'PLAYING' | 'GAME_OVER';
  CurrentTurnColor: Color;
  AvailableSpaces: Array<Coord>;
  Spaces: Spaces;
  LastCoord: Coord;
};

type IncomingMessage$Local$GameInfo = {
  name: 'local/gameInfo';
  data: GameInfo$Local;
};

const incomingMessage$GameInfo$LocalDecoder = exact({
  name: constant<'local/gameInfo'>('local/gameInfo'),
  data: exact({
    Size: number,
    Turn: number,
    ScoreData: scoreDataDecoder,
    State: either(
      constant<'PLAYING'>('PLAYING'),
      constant<'GAME_OVER'>('GAME_OVER'),
    ),
    CurrentTurnColor: colorDecoder,
    AvailableSpaces: array(coordDecoder),
    Spaces: spacesDecoder,
    LastCoord: coordDecoder,
  }),
});

export type GameInfo$Remote = {
  Size: number;
  Turn: number;
  PlayerTurn: boolean;
  OpponentID: string;
  PlayerColor: Color;
  State:
    | 'WAITING_FOR_OPPONENT'
    | 'PLAYING'
    | 'GAME_OVER_PASSED'
    | 'GAME_OVER_FORFEIT';
  ScoreData: ScoreData;
  AvailableSpaces: Array<Coord>;
  Spaces: Spaces;
  LastCoord: Coord;
};

type IncomingMessage$Remote$GameInfo = {
  name: 'remote/gameInfo';
  data: GameInfo$Remote;
};

const incomingMessage$GameInfo$RemoteDecoder = exact({
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
    ScoreData: scoreDataDecoder,
    AvailableSpaces: array(coordDecoder),
    Spaces: spacesDecoder,
    LastCoord: coordDecoder,
  }),
});

type IncomingMessage$Local$Update = {
  name: 'local/update';
  data: null;
};

const incomingMessage$Update$LocalDecoder = exact({
  name: constant<'local/update'>('local/update'),
  data: null_,
});

type IncomingMessage$Local$GameLeft = {
  name: 'local/gameLeft';
  data: null;
};

const incomingMessage$Local$GameLeftDecoder = exact({
  name: constant<'local/gameLeft'>('local/gameLeft'),
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

type IncomingMessage$Remote$Update = {
  name: 'remote/update';
  data: null;
};

const incomingMessage$Update$RemoteDecoder = exact({
  name: constant<'remote/update'>('remote/update'),
  data: null_,
});

type RejoinGameError$Local$Data = {
  Type: 'local/rejoinGame';
};

type GetGameInfoError$Local$Data = {
  Type: 'local/getGameInfo';
};

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
    | RejoinGameError$Local$Data
    | GetGameInfoError$Local$Data
    | GetGameInfoError$Remote$Data
    | RejoinGameError$Remote$Data
    | JoinGameError$Remote$Data
    | Error400$Data;
};

const incomingMessage$ErrorDecoder = exact({
  name: constant<'error'>('error'),
  data: either6(
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
    exact({
      Type: constant<'local/rejoinGame'>('local/rejoinGame'),
    }),
    exact({
      Type: constant<'local/getGameInfo'>('local/getGameInfo'),
    }),
  ),
});

type Message =
  | IncomingMessage$Local$GameInfo
  | IncomingMessage$Local$GameLeft
  | IncomingMessage$Local$GameJoined
  | IncomingMessage$Local$Update
  | IncomingMessage$Remote$GameInfo
  | IncomingMessage$Remote$GameJoined
  | IncomingMessage$Remote$GameLeft
  | IncomingMessage$Remote$Update
  | IncomingMessage$Error;

const incomingMessageDecoder = either9(
  incomingMessage$GameInfo$LocalDecoder,
  incomingMessage$Local$GameJoinedDecoder,
  incomingMessage$Local$GameLeftDecoder,
  incomingMessage$Update$LocalDecoder,
  incomingMessage$GameInfo$RemoteDecoder,
  incomingMessage$Remote$GameJoinedDecoder,
  incomingMessage$Remote$GameLeftDecoder,
  incomingMessage$Update$RemoteDecoder,
  incomingMessage$ErrorDecoder,
);

export const incomingMessageGuard: Guard<Message> = guard(
  incomingMessageDecoder,
);

export type OutgoingMessage$GetGameInfo$Local = {
  name: 'local/getGameInfo';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$RejoinGame$Local = {
  name: 'local/rejoinGame';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$CreateGame$Local = {
  name: 'local/createGame';
  data: {
    userID: string;
    size: number;
  };
};

export type OutgoingMessage$LeaveGame$Local = {
  name: 'local/leaveGame';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$Pass$Local = {
  name: 'local/pass';
  data: {
    userID: string;
    gameID: string;
  };
};

export type OutgoingMessage$PlaceStone$Local = {
  name: 'local/placeStone';
  data: {
    gameID: string;
    userID: string;
    coord: {
      X: number;
      Y: number;
    };
  };
};

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
