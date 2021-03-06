import React, { useEffect, useLayoutEffect, useState } from 'react';
import type { Coord, Spaces } from './types';

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
  size: number;
  playerColor: 'BLACK' | 'WHITE';
  canPlaceStone: boolean;
  spaces: Spaces;
  availableSpaces: Array<Coord>;
  lastCoord: Coord;
};

function Board(props: Props): JSX.Element {
  const [stoneToPlace, setStoneToPlace] = useState<Coord | null>(null);
  const [windowWidth, setWindowWidth] = useState<number>(window.innerWidth);
  const [strokeDashoffset, setStrokeDashOffset] = useState<number>(0);
  const [readyToPlace, setReadyToPlace] = useState<boolean>(false);

  useEffect(() => {
    const intervalId = setInterval(() => {
      const step = readyToPlace ? -0.6 : 0.1;
      setStrokeDashOffset(strokeDashoffset + step);
    }, 10);
    return () => clearInterval(intervalId);
  }, [strokeDashoffset]);

  useLayoutEffect(() => {
    window.addEventListener('resize', () => {
      setWindowWidth(window.innerWidth);
    });
  }, []);

  const width = windowWidth > 800 ? 800 : windowWidth - 60;
  const size = props.size;
  const rowWidth = width / (size + 1);
  const strokeWidth = width / 200 / (size / 9);
  const stoneRadius = width / 32 / (size / 9);
  const hoshiRadius = width / 80 / (size / 9);
  const hoshiPositions = getHoshiPositions(size);
  const dashArraySize = width / 80 / (size / 9);

  function isStoneToPlace(coord: Coord): boolean {
    return Boolean(
      stoneToPlace && stoneToPlace.X === coord.X && stoneToPlace.Y === coord.Y,
    );
  }

  function isLastCoord(coord: Coord): boolean {
    return props.lastCoord.X === coord.X && props.lastCoord.Y === coord.Y;
  }

  const stoneColor = props.playerColor === 'BLACK' ? 'black' : 'white';

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
          {props.spaces.BLACK.map((coord) => (
            <circle
              key={`black-${coord.X}-${coord.Y}`}
              id={`black-${coord.X}-${coord.Y}`}
              cx={rowWidth * (coord.X + 1)}
              cy={rowWidth * (coord.Y + 1)}
              r={stoneRadius}
              fill="black"
              strokeWidth={strokeWidth}
              stroke={isLastCoord(coord) ? '#00d619' : 'black'}
            />
          ))}
          {props.spaces.WHITE.map((coord) => (
            <circle
              key={`white-${coord.X}-${coord.Y}`}
              id={`white-${coord.X}-${coord.Y}`}
              cx={rowWidth * (coord.X + 1)}
              cy={rowWidth * (coord.Y + 1)}
              r={stoneRadius}
              fill="white"
              strokeWidth={strokeWidth}
              stroke={isLastCoord(coord) ? '#00d619' : 'black'}
            />
          ))}
          {props.availableSpaces.map((coord) => (
            <circle
              key={`available-${coord.X}-${coord.Y}`}
              id={`available-${coord.X}-${coord.Y}`}
              cx={rowWidth * (coord.X + 1)}
              cy={rowWidth * (coord.Y + 1)}
              r={isStoneToPlace(coord) ? stoneRadius : stoneRadius * 1.5}
              strokeWidth={strokeWidth}
              strokeDasharray={dashArraySize}
              strokeDashoffset={strokeDashoffset}
              stroke={isStoneToPlace(coord) ? '#00d619' : 'rgba(0, 0, 0, 0)'}
              fill={isStoneToPlace(coord) ? stoneColor : 'rgba(0, 0, 0, 0)'}
              onMouseEnter={() => {
                if (props.canPlaceStone) {
                  setStoneToPlace(coord);
                  setReadyToPlace(false);
                }
              }}
              onClick={() => {
                if (props.canPlaceStone && props.placeStone) {
                  if (readyToPlace) {
                    props.placeStone(coord);
                  } else {
                    setReadyToPlace(true);
                  }
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
