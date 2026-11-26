import { env } from '../config/env';
import { Matrix, MatrixStatisticsResult, NamedMatrices } from '../types/matrix.types';

/**
 * Servicio de dominio encargado de evaluar matrices y calcular sus métricas estadísticas.
 */
export class MatrixStatsService {
  private readonly epsilon: number;

  constructor(epsilon: number = env.diagonalEpsilon) {
    this.epsilon = epsilon;
  }

  /**
   * Verifica si una matriz es diagonal.
   * Por definición de álgebra lineal, una matriz es diagonal si es cuadrada (n x n)
   * y todos sus elementos fuera de la diagonal principal son cero (dentro de una tolerancia epsilon).
   */
  public isDiagonalMatrix(matrix: Matrix): boolean {
    const rows = matrix.length;
    if (rows === 0) return false;

    const cols = matrix[0].length;
    if (rows !== cols) {
      return false;
    }

    for (let i = 0; i < rows; i++) {
      for (let j = 0; j < cols; j++) {
        if (i !== j && Math.abs(matrix[i][j]) > this.epsilon) {
          return false;
        }
      }
    }

    return true;
  }

  /**
   * Calcula el valor máximo, mínimo, suma total, promedio y si alguna matriz es diagonal
   * recorriendo todas las matrices recibidas en una sola pasada O(N) para eficiencia de memoria.
   */
  public calculateStatistics(matrices: NamedMatrices): MatrixStatisticsResult {
    const entries = Object.entries(matrices);
    if (entries.length === 0) {
      throw new Error('Debe proporcionarse al menos una matriz para calcular estadísticas.');
    }

    let max = -Infinity;
    let min = Infinity;
    let sum = 0;
    let totalElements = 0;
    let isAnyDiagonal = false;
    const diagonalByMatrix: Record<string, boolean> = {};

    for (const [name, matrix] of entries) {
      const isDiagonal = this.isDiagonalMatrix(matrix);
      diagonalByMatrix[name] = isDiagonal;
      if (isDiagonal) {
        isAnyDiagonal = true;
      }

      for (const row of matrix) {
        for (const value of row) {
          if (value > max) max = value;
          if (value < min) min = value;
          sum += value;
          totalElements++;
        }
      }
    }

    const average = totalElements > 0 ? sum / totalElements : 0;

    return {
      max: this.roundPrecision(max),
      min: this.roundPrecision(min),
      average: this.roundPrecision(average),
      sum: this.roundPrecision(sum),
      totalElements,
      isAnyDiagonal,
      diagonalByMatrix,
    };
  }

  /** Normaliza ceros negativos (-0) y mantiene alta precisión numérica en JSON. */
  private roundPrecision(value: number): number {
    if (Object.is(value, -0) || Math.abs(value) < 1e-12) {
      return 0;
    }
    return Number(value.toFixed(10));
  }
}
