import { NextFunction, Request, Response } from 'express';
import { z } from 'zod';
import { MatrixStatsService } from '../services/matrixStats.service';
import { Matrix, NamedMatrices } from '../types/matrix.types';

/** Esquema de validación para una matriz rectangular no vacía de números finitos. */
const matrixSchema = z
  .array(z.array(z.number().finite('Todos los elementos deben ser números finitos.')))
  .min(1, 'La matriz debe tener al menos una fila.')
  .refine((rows) => rows[0].length > 0, {
    message: 'La matriz debe tener al menos una columna.',
  })
  .refine((rows) => rows.every((row) => row.length === rows[0].length), {
    message: 'Todas las filas de la matriz deben tener la misma longitud (matriz rectangular).',
  });

/**
 * Acepta tanto un mapa nombrado `{ matrices: { Q: [...], R: [...] } }` (enviado por la API en Go)
 * como un arreglo de matrices `{ matrices: [ [...], [...] ] }` para facilitar pruebas directas.
 */
const statsPayloadSchema = z.object({
  matrices: z.union([
    z
      .record(z.string().min(1), matrixSchema)
      .refine((obj) => Object.keys(obj).length > 0, {
        message: 'Debe incluir al menos una matriz en el objeto matrices.',
      }),
    z.array(matrixSchema).min(1, 'Debe incluir al menos una matriz en el arreglo matrices.'),
  ]),
});

/**
 * Controlador HTTP para el cálculo de estadísticas matriciales.
 */
export class StatsController {
  constructor(private readonly statsService: MatrixStatsService) {}

  /** Valida el payload de matrices recibido y retorna las estadísticas consolidadas. */
  public computeStats = (req: Request, res: Response, next: NextFunction): void => {
    try {
      const parsed = statsPayloadSchema.parse(req.body);
      const namedMatrices = this.normalizeMatrices(parsed.matrices);
      const statistics = this.statsService.calculateStatistics(namedMatrices);

      res.status(200).json(statistics);
    } catch (error) {
      next(error);
    }
  };

  /** Normaliza arreglos de matrices sin nombre a un mapa con claves deterministas. */
  private normalizeMatrices(input: NamedMatrices | Matrix[]): NamedMatrices {
    if (!Array.isArray(input)) {
      return input;
    }

    const result: NamedMatrices = {};
    input.forEach((matrix, index) => {
      result[`matrix_${index + 1}`] = matrix;
    });
    return result;
  }
}
