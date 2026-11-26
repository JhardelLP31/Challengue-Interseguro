import { Router } from 'express';
import { StatsController } from '../controllers/stats.controller';
import { authenticateJwt } from '../middlewares/auth.middleware';
import { MatrixStatsService } from '../services/matrixStats.service';

/**
 * Define las rutas protegidas para el cálculo estadístico sobre matrices.
 */
export const createStatsRouter = (): Router => {
  const router = Router();
  const statsService = new MatrixStatsService();
  const controller = new StatsController(statsService);

  router.post('/', authenticateJwt, controller.computeStats);

  return router;
};
