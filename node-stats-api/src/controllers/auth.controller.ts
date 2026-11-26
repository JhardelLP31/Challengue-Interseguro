import { NextFunction, Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import { z } from 'zod';
import { env } from '../config/env';

const loginSchema = z.object({
  username: z.string().min(1, 'El campo username es requerido.'),
  password: z.string().min(1, 'El campo password es requerido.'),
});

/**
 * Controlador de autenticación JWT.
 * Emite tokens HS256 firmados con el mismo secreto compartido con la API en Go.
 */
export class AuthController {
  /** Valida credenciales de cliente y genera un token JWT Bearer firmado. */
  public issueToken = (req: Request, res: Response, next: NextFunction): void => {
    try {
      const { username, password } = loginSchema.parse(req.body);

      if (username !== env.authUsername || password !== env.authPassword) {
        res.status(401).json({
          error: 'Unauthorized',
          message: 'Credenciales inválidas.',
        });
        return;
      }

      const token = jwt.sign(
        { sub: username, role: 'api-client' },
        env.jwtSecret,
        { expiresIn: env.jwtExpiresIn as jwt.SignOptions['expiresIn'] },
      );

      res.status(200).json({
        token,
        tokenType: 'Bearer',
        expiresIn: env.jwtExpiresIn,
      });
    } catch (error) {
      next(error);
    }
  };
}
