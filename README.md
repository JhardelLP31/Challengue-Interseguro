# Interseguro Coding Challenge — Matrix QR & Statistics Microservices

Solución distribuida compuesta por **dos APIs RESTful independientes** comunicadas mediante **HTTP**, protegidas con autenticación **JWT (HS256)**, contenerizadas con **Docker** y respaldadas por pruebas unitarias y de integración.

## Arquitectura de la Solución

```text
Cliente ──(1. POST /api/v1/auth/token)──────────────► go-matrix-api (Go + Fiber :8080)
        ◄─(JWT HS256)────────────────────────────────┘
                                                       │
Cliente ──(2. POST /api/v1/matrices/qr + Bearer JWT)──►│  Valida matriz rectangular,
                                                       │  calcula QR (Householder) o rotación 90°
                                                       │
                                                       ├──(3. HTTP POST /api/v1/stats + Bearer JWT)──► node-stats-api (Express + TS :3000)
                                                       │                                               Calcula max, min, sum, avg, diagonal
                                                       ◄──(4. JSON Estadísticas)───────────────────────┘
        ◄─(5. JSON Resultado + Estadísticas)───────────┘
```

### Estructura del Repositorio

| Componente | Stack | Responsabilidad | Puerto por defecto |
| --- | --- | --- | --- |
| [`go-matrix-api`](./go-matrix-api) | Go 1.22+ · Fiber v2 | Validación de matrices, factorización QR, rotación 90° y orquestación HTTP | `8080` |
| [`node-stats-api`](./node-stats-api) | Node.js 22 · Express · TypeScript | Cálculo estadístico ($O(N)$) y verificación de matriz diagonal | `3000` |
| [`docker-compose.yml`](./docker-compose.yml) | Docker Compose | Orquestación de ambos servicios en red privada con healthcheck | `8080`, `3000` |

---

## Decisiones Técnicas y de Diseño

1. **Resolución de Factorización QR y Rotación de Matriz:**
   - Ante las referencias del enunciado a *factorización QR* (sobre múltiples matrices devueltas y verificación de matriz diagonal) y *rotación de matriz*, se implementó la **Factorización QR como flujo principal** (`POST /api/v1/matrices/qr`) y se expuso adicionalmente el **endpoint de rotación de 90°** (`POST /api/v1/matrices/rotate`), integrando ambos flujos vía HTTP con `node-stats-api`.

2. **Factorización QR completa mediante Reflexiones de Householder:**
   - Se descompone toda matriz rectangular $A \in \mathbb{R}^{m \times n}$ en $A = Q \cdot R$, donde $Q \in \mathbb{R}^{m \times m}$ es ortogonal ($Q^T Q = I_m$) y $R \in \mathbb{R}^{m \times n}$ es triangular/trapezoidal superior.
   - El método de Householder ofrece mayor estabilidad numérica que Gram-Schmidt clásico y soporta matrices cuadradas, altas ($m > n$) y anchas ($m < n$).

3. **Detección de Matriz Diagonal con Tolerancia Numérica ($\varepsilon = 10^{-9}$):**
   - Una matriz se considera diagonal si es cuadrada ($n \times n$) y todos sus elementos fuera de la diagonal principal cumplen $|M_{i,j}| \le \varepsilon$, absorbiendo el ruido infinitesimal de aritmética `float64` propio de las descomposiciones matriciales.

4. **Seguridad con JWT Compartido (HS256):**
   - Ambas APIs validan el header `Authorization: Bearer <token>`. `go-matrix-api` propaga el token del cliente en la petición HTTP interna hacia `node-stats-api`.

---

## Puesta en Marcha

### Opción 1: Docker Compose (Recomendada)

```bash
docker compose up --build
```

- `go-matrix-api`: `http://localhost:8080`
- `node-stats-api`: `http://localhost:3000`

### Opción 2: Ejecución Local Independiente

1. **Iniciar `node-stats-api`:**
   ```bash
   cd node-stats-api
   npm install
   npm run dev
   ```

2. **Iniciar `go-matrix-api` (en otra terminal):**
   ```bash
   cd go-matrix-api
   go mod tidy
   go run ./cmd/server
   ```

---

## Pruebas Automatizadas

### `node-stats-api` (Vitest + Supertest)
```bash
cd node-stats-api
npm test
```

### `go-matrix-api` (Go testing + httptest)
```bash
cd go-matrix-api
go test -v ./...
```

---

## Referencia Rápida de Consumo (HTTP)

### 1. Generar Token JWT
```bash
curl -X POST http://localhost:8080/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{"username":"interseguro","password":"challenge2024"}'
```

### 2. Factorización QR + Estadísticas
```bash
curl -X POST http://localhost:8080/api/v1/matrices/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN_JWT>" \
  -d '{"matrix":[[1,2],[3,4],[5,6]]}'
```

### 3. Rotación de Matriz + Estadísticas
```bash
curl -X POST http://localhost:8080/api/v1/matrices/rotate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN_JWT>" \
  -d '{"matrix":[[1,2,3],[4,5,6]],"direction":"clockwise"}'
```

---

## Despliegue en la Nube

Cada servicio cuenta con un `Dockerfile` multi-stage listo para plataformas como **Render**, **Railway**, **Fly.io** o **Google Cloud Run**:

1. Desplegar primero **`node-stats-api`** (directorio raíz `node-stats-api`, puerto dinámico `PORT`) definiendo `JWT_SECRET`, `AUTH_USERNAME` y `AUTH_PASSWORD`.
2. Desplegar **`go-matrix-api`** (directorio raíz `go-matrix-api`, puerto dinámico `PORT`) configurando el mismo `JWT_SECRET`, credenciales de autenticación, `STATS_API_URL` apuntando a la URL pública de `node-stats-api`, y `HTTP_TIMEOUT_SECONDS=60` si la plataforma suspende instancias inactivas.
