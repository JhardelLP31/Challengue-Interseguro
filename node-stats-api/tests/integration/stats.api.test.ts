import request from 'supertest';
import { describe, expect, it } from 'vitest';
import { createApp } from '../../src/app';

describe('node-stats-api Integration Tests', () => {
  const app = createApp();

  const getJwtToken = async (): Promise<string> => {
    const response = await request(app).post('/api/v1/auth/token').send({
      username: 'interseguro',
      password: 'challenge2024',
    });
    expect(response.status).toBe(200);
    return response.body.token as string;
  };

  it('GET /health responde 200 OK sin requerir autenticación', async () => {
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
    expect(res.body.service).toBe('node-stats-api');
  });

  it('POST /api/v1/stats rechaza peticiones sin token JWT (401)', async () => {
    const res = await request(app)
      .post('/api/v1/stats')
      .send({
        matrices: {
          Q: [[1, 0], [0, 1]],
        },
      });

    expect(res.status).toBe(401);
  });

  it('POST /api/v1/stats rechaza tokens JWT inválidos (401)', async () => {
    const res = await request(app)
      .post('/api/v1/stats')
      .set('Authorization', 'Bearer invalid.jwt.token')
      .send({
        matrices: {
          Q: [[1, 0], [0, 1]],
        },
      });

    expect(res.status).toBe(401);
  });

  it('POST /api/v1/stats calcula estadísticas con formato objeto nombrado { Q, R }', async () => {
    const token = await getJwtToken();

    const res = await request(app)
      .post('/api/v1/stats')
      .set('Authorization', `Bearer ${token}`)
      .send({
        matrices: {
          Q: [
            [1, 0],
            [0, 1],
          ],
          R: [
            [4, 2],
            [0, 3],
          ],
        },
      });

    expect(res.status).toBe(200);
    expect(res.body).toEqual({
      max: 4,
      min: 0,
      average: 1.375,
      sum: 11,
      totalElements: 8,
      isAnyDiagonal: true,
      diagonalByMatrix: {
        Q: true,
        R: false,
      },
    });
  });

  it('POST /api/v1/stats también soporta matrices en arreglo y valida matrices irregulares (400)', async () => {
    const token = await getJwtToken();

    const okRes = await request(app)
      .post('/api/v1/stats')
      .set('Authorization', `Bearer ${token}`)
      .send({
        matrices: [
          [
            [2, 0],
            [0, 5],
          ],
        ],
      });

    expect(okRes.status).toBe(200);
    expect(okRes.body.isAnyDiagonal).toBe(true);
    expect(okRes.body.max).toBe(5);

    const badRes = await request(app)
      .post('/api/v1/stats')
      .set('Authorization', `Bearer ${token}`)
      .send({
        matrices: {
          Q: [
            [1, 2],
            [3],
          ],
        },
      });

    expect(badRes.status).toBe(400);
    expect(badRes.body.error).toBe('ValidationError');
  });
});
