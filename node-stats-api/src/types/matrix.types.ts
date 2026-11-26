/** Representa una matriz numérica bidimensional (array de filas). */
export type Matrix = number[][];

/** Mapa nombrado de matrices (ejemplo: { Q: [...], R: [...] }). */
export type NamedMatrices = Record<string, Matrix>;

/** Resultado consolidado de las estadísticas calculadas sobre las matrices. */
export interface MatrixStatisticsResult {
  max: number;
  min: number;
  average: number;
  sum: number;
  totalElements: number;
  isAnyDiagonal: boolean;
  diagonalByMatrix: Record<string, boolean>;
}

/** Payload del token JWT decodificado. */
export interface JwtPayload {
  sub: string;
  role: string;
  iat?: number;
  exp?: number;
}
