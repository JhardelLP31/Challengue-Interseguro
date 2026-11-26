import dotenv from 'dotenv';

dotenv.config();

const parsePositiveNumber = (value: string | undefined, fallback: number): number => {
  if (!value) return fallback;
  const parsed = Number(value);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
};

/**
 * Configuración centralizada de la API de estadísticas.
 * Se carga desde variables de entorno con valores seguros para desarrollo local.
 */
export const env = {
  port: parsePositiveNumber(process.env.PORT, 3000),
  nodeEnv: process.env.NODE_ENV ?? 'development',
  jwtSecret: process.env.JWT_SECRET ?? 'interseguro-coding-challenge-super-secret-key-2024',
  jwtExpiresIn: process.env.JWT_EXPIRES_IN ?? '24h',
  authUsername: process.env.AUTH_USERNAME ?? 'interseguro',
  authPassword: process.env.AUTH_PASSWORD ?? 'challenge2024',
  diagonalEpsilon: parsePositiveNumber(process.env.DIAGONAL_EPSILON, 1e-9),
};
