import { describe, expect, it } from 'vitest';
import { MatrixStatsService } from '../../src/services/matrixStats.service';

describe('MatrixStatsService', () => {
  const service = new MatrixStatsService(1e-9);

  it('calcula max, min, suma, promedio y diagonal correctamente para múltiples matrices', () => {
    const result = service.calculateStatistics({
      Q: [
        [1, 0],
        [0, 1],
      ],
      R: [
        [2, 3],
        [0, 4],
      ],
    });

    expect(result.max).toBe(4);
    expect(result.min).toBe(0);
    expect(result.sum).toBe(11);
    expect(result.totalElements).toBe(8);
    expect(result.average).toBeCloseTo(1.375, 10);
    expect(result.isAnyDiagonal).toBe(true);
    expect(result.diagonalByMatrix).toEqual({
      Q: true,
      R: false,
    });
  });

  it('maneja valores negativos y decimales correctamente', () => {
    const result = service.calculateStatistics({
      Q: [
        [-0.8, 0.6],
        [-0.6, -0.8],
      ],
      R: [
        [-5, 1],
        [0, -2],
      ],
    });

    expect(result.max).toBe(1);
    expect(result.min).toBe(-5);
    expect(result.sum).toBeCloseTo(-7.6, 10);
    expect(result.totalElements).toBe(8);
    expect(result.average).toBeCloseTo(-0.95, 10);
    expect(result.isAnyDiagonal).toBe(false);
  });

  it('no marca como diagonal una matriz rectangular aunque tenga ceros fuera de la diagonal', () => {
    const rectangular = [
      [3, 0, 0],
      [0, 5, 0],
    ];
    expect(service.isDiagonalMatrix(rectangular)).toBe(false);
  });

  it('considera tolerancia epsilon para errores numéricos de punto flotante en matrices diagonales', () => {
    const almostDiagonal = [
      [5, 1e-12],
      [-2e-12, -3],
    ];
    expect(service.isDiagonalMatrix(almostDiagonal)).toBe(true);
  });
});
