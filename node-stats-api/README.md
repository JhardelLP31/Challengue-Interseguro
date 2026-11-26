# `node-stats-api` (Node.js + Express + TypeScript)

Servicio RESTful desarrollado con **Node.js**, **Express.js** y **TypeScript** encargado de validar colecciones de matrices numéricas y calcular estadísticas agregadas junto con la detección de matrices diagonales.

## Estructura del Proyecto

```text
node-stats-api/
├── src/
│   ├── config/env.ts                    # Carga de variables de entorno
│   ├── controllers/
│   │   ├── auth.controller.ts           # Emisión de tokens JWT (HS256)
│   │   └── stats.controller.ts          # Validación de payload con Zod y respuesta HTTP
│   ├── middlewares/
│   │   ├── auth.middleware.ts           # Verificación de Authorization: Bearer <token>
│   │   └── error.middleware.ts          # Manejador centralizado de errores HTTP
│   ├── routes/
│   │   ├── auth.routes.ts               # Rutas de autenticación
│   │   └── stats.routes.ts              # Rutas de cálculo estadístico
│   ├── services/
│   │   └── matrixStats.service.ts       # Lógica estadística O(N) y verificación diagonal
│   ├── types/matrix.types.ts            # Definiciones e interfaces TypeScript
│   ├── app.ts                           # Configuración de Express y middlewares globales
│   └── server.ts                        # Arranque del servidor HTTP
├── tests/
│   ├── unit/matrixStats.service.test.ts # Pruebas unitarias del servicio estadístico
│   └── integration/stats.api.test.ts    # Pruebas de integración HTTP con Supertest
├── Dockerfile
├── tsconfig.json
└── package.json
```

## Variables de Entorno

| Variable | Descripción | Valor por defecto |
| --- | --- | --- |
| `PORT` | Puerto HTTP de escucha | `3000` |
| `NODE_ENV` | Entorno de ejecución (`development` / `production`) | `development` |
| `JWT_SECRET` | Secreto compartido para firma/validación HS256 | `interseguro-coding-challenge-super-secret-key-2024` |
| `JWT_EXPIRES_IN` | Duración del token emitido (`24h`, `1h`, etc.) | `24h` |
| `AUTH_USERNAME` | Usuario habilitado para `POST /api/v1/auth/token` | `interseguro` |
| `AUTH_PASSWORD` | Contraseña habilitada para `POST /api/v1/auth/token` | `challenge2024` |
| `DIAGONAL_EPSILON` | Tolerancia numérica para elementos fuera de la diagonal | `1e-9` |

## Endpoints

### `GET /health` *(público)*
Healthcheck del servicio.

### `POST /api/v1/auth/token` *(público)*
Emite un JWT compartido con `go-matrix-api` para pruebas independientes.

### `POST /api/v1/stats` *(protegido con JWT)*
Acepta un mapa nombrado de matrices (`{ "matrices": { "Q": [...], "R": [...] } }`) o un arreglo de matrices (`{ "matrices": [ [...], [...] ] }`) y calcula:

- `max`: valor máximo encontrado entre todas las matrices.
- `min`: valor mínimo encontrado entre todas las matrices.
- `average`: promedio aritmético de todos los elementos.
- `sum`: suma total de todos los valores.
- `totalElements`: cantidad total de números procesados.
- `isAnyDiagonal`: `true` si al menos una matriz es cuadrada y diagonal dentro de la tolerancia `DIAGONAL_EPSILON`.
- `diagonalByMatrix`: mapa booleano indicando el resultado por cada matriz evaluada.

## Ejecución y Pruebas

```bash
npm install
npm run build
npm test
npm run dev
```
