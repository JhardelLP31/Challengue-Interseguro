import { NextFunction, Request, Response } from 'express';
import { ZodError } from 'zod';

/** Error HTTP personalizado para respuestas controladas de negocio/validación. */
export class HttpError extends Error {
  public readonly statusCode: number;

  constructor(statusCode: number, message: string) {
    super(message);
    this.statusCode = statusCode;
    this.name = 'HttpError';
  }
}

/**
 * Middleware global de manejo de errores de Express.
 * Convierte errores de Zod, JSON malformado y HttpError en respuestas HTTP uniformes.
 */
export const errorHandler = (
  err: unknown,
  _req: Request,
  res: Response,
  _next: NextFunction,
): void => {
  if (err instanceof HttpError) {
    res.status(err.statusCode).json({
      error: err.statusCode >= 500 ? 'InternalServerError' : 'BadRequest',
      message: err.message,
    });
    return;
  }

  if (err instanceof ZodError) {
    res.status(400).json({
      error: 'ValidationError',
      message: 'El payload enviado no cumple con el esquema esperado.',
      details: err.errors.map((issue) => ({
        path: issue.path.join('.'),
        message: issue.message,
      })),
    });
    return;
  }

  if (err instanceof SyntaxError && 'body' in err) {
    res.status(400).json({
      error: 'InvalidJson',
      message: 'El cuerpo de la petición no es un JSON válido.',
    });
    return;
  }

  res.status(500).json({
    error: 'InternalServerError',
    message: 'Ocurrió un error interno inesperado.',
  });
};
