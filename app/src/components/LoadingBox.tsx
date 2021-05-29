import React, { useEffect, useState } from 'react';
import type { Coord } from './types';

type Placement = {
  color: 'black' | 'white';
  coord: Coord;
};

type Step = {
  type: 'place' | 'remove';
  placement: Placement;
};

function Board(): JSX.Element {
  const [windowWidth, setWindowWidth] = useState<number>(window.innerWidth);
  const [stepNumber, setStepNumber] = useState<number>(0);
  const [placements, setPlacements] = useState<Array<Placement>>([]);

  const steps: Array<Step> = [
    { type: 'place', placement: { coord: { X: 0, Y: 0 }, color: 'black' } },
    { type: 'place', placement: { coord: { X: 1, Y: 0 }, color: 'white' } },
    { type: 'place', placement: { coord: { X: 2, Y: 0 }, color: 'black' } },
    { type: 'place', placement: { coord: { X: 2, Y: 1 }, color: 'white' } },
    { type: 'place', placement: { coord: { X: 1, Y: 1 }, color: 'black' } },
    { type: 'place', placement: { coord: { X: 0, Y: 1 }, color: 'white' } },
    { type: 'place', placement: { coord: { X: 0, Y: 2 }, color: 'black' } },
    { type: 'place', placement: { coord: { X: 1, Y: 2 }, color: 'white' } },
    { type: 'place', placement: { coord: { X: 2, Y: 2 }, color: 'black' } },
    { type: 'remove', placement: { coord: { X: 0, Y: 0 }, color: 'black' } },
    { type: 'remove', placement: { coord: { X: 1, Y: 0 }, color: 'white' } },
    { type: 'remove', placement: { coord: { X: 2, Y: 0 }, color: 'black' } },
    { type: 'remove', placement: { coord: { X: 2, Y: 1 }, color: 'white' } },
    { type: 'remove', placement: { coord: { X: 1, Y: 1 }, color: 'black' } },
    { type: 'remove', placement: { coord: { X: 0, Y: 1 }, color: 'white' } },
    { type: 'remove', placement: { coord: { X: 0, Y: 2 }, color: 'black' } },
    { type: 'remove', placement: { coord: { X: 1, Y: 2 }, color: 'white' } },
    { type: 'remove', placement: { coord: { X: 2, Y: 2 }, color: 'black' } },
  ];

  useEffect(() => {
    window.addEventListener('resize', () => {
      setWindowWidth(window.innerWidth);
    });
  }, []);

  useEffect(() => {
    const intervalId = setInterval(() => {
      console.log('Step!');
      const step = steps[stepNumber];
      if (step.type === 'place') {
        setPlacements([...placements, step.placement]);
      } else {
        let newPlacements = [...placements];
        for (let i = 0; i < placements.length; i++) {
          const placement = placements[i];
          if (
            placement.color === step.placement.color &&
            placement.coord.X === step.placement.coord.X &&
            placement.coord.Y === step.placement.coord.Y
          ) {
            newPlacements = [
              ...newPlacements.slice(0, i),
              ...newPlacements.slice(i + 1, newPlacements.length),
            ];
          }
        }
        setPlacements(newPlacements);
      }
      if (stepNumber < steps.length - 1) {
        setStepNumber(stepNumber + 1);
      } else {
        setStepNumber(0);
      }
    }, 200);
    return () => {
      clearInterval(intervalId);
    };
  }, [stepNumber]);

  const width = windowWidth > 200 ? 200 : windowWidth - 60;
  const size = 3;
  const rowWidth = width / (size + 1);
  const strokeWidth = width / 200 / (size / 9);
  const stoneRadius = width / 32 / (size / 9);

  return (
    <div style={{ margin: '20px' }}>
      <svg
        width={width}
        height={width}
        style={{ backgroundColor: '#ffc4fb', borderRadius: '10px' }}
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
          {placements
            .filter((p) => p.color === 'black')
            .map((p) => (
              <circle
                key={`black-${p.coord.X}-${p.coord.Y}`}
                id={`black-${p.coord.X}-${p.coord.Y}`}
                cx={rowWidth * (p.coord.X + 1)}
                cy={rowWidth * (p.coord.Y + 1)}
                r={stoneRadius}
                fill="black"
                strokeWidth={strokeWidth}
                stroke="black"
              />
            ))}
          {placements
            .filter((p) => p.color === 'white')
            .map((p) => (
              <circle
                key={`white-${p.coord.X}-${p.coord.Y}`}
                id={`white-${p.coord.X}-${p.coord.Y}`}
                cx={rowWidth * (p.coord.X + 1)}
                cy={rowWidth * (p.coord.Y + 1)}
                r={stoneRadius}
                fill="white"
                strokeWidth={strokeWidth}
                stroke="black"
              />
            ))}
        </g>
      </svg>
    </div>
  );
}
export default Board;
