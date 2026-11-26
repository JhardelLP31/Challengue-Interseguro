import cors from 'cors';
import express, { Application } from 'express';
import helmet from 'helmet';
import { errorHandler } from './middlewares/error.middleware';
import { createAuthRouter } from './routes/auth.routes';
import { createStatsRouter } from './routes/stats.routes';

/**
 * Construye y configura la instancia de Express.
 * Separada de `server.ts` para facilitar las pruebas de integración con Supertest.
 */
export const createApp = (): Application => {
  const app = express();

  app.use(helmet());
  app.use(cors());
  app.use(express.json({ limit: '5mb' }));

  app.get('/health', (_req, res) => {
    res.status(200).json({
      status: 'ok',
      service: 'node-stats-api',
      timestamp: new Date().toISOString(),
    });
  });

  app.use('/api/v1/auth', createAuthRouter());
  app.use('/api/v1/stats', createStatsRouter());

  app.use((_req, res) => {
    res.status(404).json({
      error: 'NotFound',
      message: 'Ruta no encontrada.',
    });
  });

  app.use(errorHandler);

  return app;
};
