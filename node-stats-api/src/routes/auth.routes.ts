import { Router } from 'express';
import { AuthController } from '../controllers/auth.controller';

/**
 * Define las rutas de autenticación para emisión de tokens JWT.
 */
export const createAuthRouter = (): Router => {
  const router = Router();
  const controller = new AuthController();

  router.post('/token', controller.issueToken);

  return router;
};
