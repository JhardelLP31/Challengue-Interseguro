import { NextFunction, Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import { env } from '../config/env';
import { JwtPayload } from '../types/matrix.types';

/** Extiende Request de Express para adjuntar el usuario autenticado por JWT. */
export interface AuthenticatedRequest extends Request {
  user?: JwtPayload;
}

/**
 * Middleware que valida el token JWT enviado en el header `Authorization: Bearer <token>`.
 */
export const authenticateJwt = (
  req: AuthenticatedRequest,
  res: Response,
  next: NextFunction,
): void => {
  const authHeader = req.headers.authorization;
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    res.status(401).json({
      error: 'Unauthorized',
      message: 'Falta el header Authorization con formato Bearer <token>.',
    });
    return;
  }

  const token = authHeader.slice('Bearer '.length).trim();
  if (!token) {
    res.status(401).json({
      error: 'Unauthorized',
      message: 'El token JWT no puede estar vacío.',
    });
    return;
  }

  try {
    const decoded = jwt.verify(token, env.jwtSecret) as JwtPayload;
    req.user = decoded;
    next();
  } catch {
    res.status(401).json({
      error: 'Unauthorized',
      message: 'Token JWT inválido o expirado.',
    });
  }
};
