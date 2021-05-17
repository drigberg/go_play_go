import React, { useEffect, useState } from 'react';
import type { Coord, GameInfo } from './types';

function getHoshiPositions(size: number): Array<number> {
  switch (size) {
    case 9:
      return [2, 4, 6];
    case 13:
      return [3, 6, 9];
    case 19:
      return [3, 9, 15];
    default:
      throw new Error('Unrecognized board size: this should not be possible');
  }
}

type Props = {
  placeStone?: (coord: Coord) => void;
  gameInfo: GameInfo;
};

function Board(props: Props): JSX.Element {
  const [stoneToPlace, setStoneToPlace] = useState<Coord | null>(null);
  const [windowWidth, setWindowWidth] = useState<number>(window.innerWidth);

  useEffect(() => {
    window.addEventListener('resize', () => {
      setWindowWidth(window.innerWidth);
    });
  }, []);

  const width = windowWidth > 800 ? 800 : windowWidth - 60;
  const size = props.gameInfo.Size;
  const rowWidth = width / (size + 1);
  const strokeWidth = width / 200 / (size / 9);
  const stoneRadius = width / 32 / (size / 9);
  const hoshiRadius = width / 80 / (size / 9);
  const hoshiPositions = getHoshiPositions(size);

  function isStoneToPlace(coord: Coord): boolean {
    return Boolean(
      stoneToPlace && stoneToPlace.X === coord.X && stoneToPlace.Y === coord.Y,
    );
  }

  const stoneColor = props.gameInfo.PlayerColor === 'BLACK' ? 'black' : 'white';
  const canPlaceStone =
    props.gameInfo.PlayerTurn && props.gameInfo.State === 'PLAYING';

  return (
    <div style={{ margin: '20px' }}>
      <svg
        width={width}
        height={width}
        onMouseLeave={() => {
          setStoneToPlace(null);
        }}
        style={{ backgroundColor: '#ffc4fb' }}
        xmlns="http://www.w3.org/2000/svg"
        version="1.1"
      >
        <g stroke="black">
          {new Array(size).fill(null).map((_, x) => (
            <line
              id={`row-${x}`}
              key={`row-${x}`}
              x1={rowWidth}
              y1={rowWidth * (x + 1)}
              x2={rowWidth * size}
              y2={rowWidth * (x + 1)}
              strokeWidth={strokeWidth}
            />
          ))}
          {new Array(size).fill(null).map((_, x) => (
            <line
              id={`column-${x}`}
              key={`column-${x}`}
              x1={rowWidth * (x + 1)}
              y1={rowWidth}
              x2={rowWidth * (x + 1)}
              y2={rowWidth * size}
              strokeWidth={strokeWidth}
            />
          ))}
          {hoshiPositions.map((x) =>
            hoshiPositions.map((y) => (
              <circle
                key={`circle-${x}-${y}`}
                cx={rowWidth * (x + 1)}
                cy={rowWidth * (y + 1)}
                r={hoshiRadius}
                fill="black"
              />
            )),
          )}
          {props.gameInfo.Spaces.BLACK.map((coord) => (
            <circle
              key={`black-${coord.X}-${coord.Y}`}
              id={`black-${coord.X}-${coord.Y}`}
              cx={rowWidth * (coord.X + 1)}
              cy={rowWidth * (coord.Y + 1)}
              r={stoneRadius}
              fill="black"
              strokeWidth={strokeWidth}
              stroke="black"
            />
          ))}
          {props.gameInfo.Spaces.WHITE.map((coord) => (
            <circle
              key={`white-${coord.X}-${coord.Y}`}
              id={`white-${coord.X}-${coord.Y}`}
              cx={rowWidth * (coord.X + 1)}
              cy={rowWidth * (coord.Y + 1)}
              r={stoneRadius}
              fill="white"
              strokeWidth={strokeWidth}
              stroke="black"
            />
          ))}
          {props.gameInfo.AvailableSpaces.map((coord) => (
            <circle
              key={`available-${coord.X}-${coord.Y}`}
              id={`available-${coord.X}-${coord.Y}`}
              cx={rowWidth * (coord.X + 1)}
              cy={rowWidth * (coord.Y + 1)}
              r={isStoneToPlace(coord) ? stoneRadius : stoneRadius * 1.5}
              strokeWidth={strokeWidth}
              stroke={isStoneToPlace(coord) ? '#00d619' : 'rgba(0, 0, 0, 0)'}
              fill={isStoneToPlace(coord) ? stoneColor : 'rgba(0, 0, 0, 0)'}
              onMouseEnter={() => {
                if (canPlaceStone) {
                  setStoneToPlace(coord);
                }
              }}
              onClick={() => {
                if (canPlaceStone && props.placeStone) {
                  props.placeStone(coord);
                }
              }}
            />
          ))}
        </g>
      </svg>
    </div>
  );
}
export default Board;
