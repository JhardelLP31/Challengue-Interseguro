# `go-matrix-api` (Go + Fiber)

Servicio RESTful implementado en **Go** con **Fiber v2** responsable de validar matrices rectangulares, calcular su **Factorización QR** (y rotación de 90°), gestionar autenticación JWT y comunicarse por **HTTP** con `node-stats-api` para obtener métricas consolidadas.

## Estructura del Proyecto

```text
go-matrix-api/
├── cmd/server/main.go               # Punto de entrada del servidor HTTP
├── internal/
│   ├── app/
│   │   ├── app.go                   # Inicialización de Fiber, middlewares y enrutamiento
│   │   └── app_test.go              # Pruebas de integración HTTP (con mock de stats API)
│   ├── client/
│   │   └── stats_client.go          # Cliente HTTP hacia node-stats-api
│   ├── config/
│   │   └── config.go                # Carga tipada de variables de entorno
│   ├── handler/
│   │   ├── auth_handler.go          # Generación de tokens JWT (HS256)
│   │   └── matrix_handler.go        # Controladores de factorización QR y rotación
│   ├── middleware/
│   │   └── auth.go                  # Middleware de verificación Bearer JWT
│   ├── model/
│   │   └── matrix.go                # Estructuras de petición, respuesta y errores
│   └── service/
│       ├── matrix_service.go        # Lógica matemática (QR Householder y rotación)
│       └── matrix_service_test.go   # Pruebas unitarias de álgebra matricial
├── Dockerfile
├── .env.example
└── go.mod
```

## Variables de Entorno

| Variable | Descripción | Valor por defecto |
| --- | --- | --- |
| `PORT` | Puerto HTTP de escucha | `8080` |
| `STATS_API_URL` | URL base del servicio `node-stats-api` | `http://localhost:3000` |
| `JWT_SECRET` | Secreto compartido para firma/validación HS256 | `interseguro-coding-challenge-super-secret-key-2024` |
| `JWT_EXPIRES_IN_HOURS` | Vigencia del token emitido en horas | `24` |
| `AUTH_USERNAME` | Usuario habilitado para `POST /api/v1/auth/token` | `interseguro` |
| `AUTH_PASSWORD` | Contraseña habilitada para `POST /api/v1/auth/token` | `challenge2024` |
| `HTTP_TIMEOUT_SECONDS` | Timeout del cliente HTTP hacia `node-stats-api` | `10` |

## Endpoints

### `GET /health` *(público)*
Verifica la disponibilidad del servicio.

### `POST /api/v1/auth/token` *(público)*
Emite un token JWT `Bearer` válido para consumir ambas APIs.
```json
{
  "username": "interseguro",
  "password": "challenge2024"
}
```

### `POST /api/v1/matrices/qr` *(protegido con JWT)*
Recibe una matriz rectangular (`rows >= 1`, `cols >= 1`), calcula $Q$ y $R$ mediante reflexiones de Householder y delega el cálculo estadístico a `node-stats-api`.

**Request:**
```json
{
  "matrix": [
    [1, 2],
    [3, 4],
    [5, 6]
  ]
}
```

**Response (`200 OK`):**
```json
{
  "operation": "qr_factorization",
  "input": { "rows": 3, "cols": 2, "matrix": [[1, 2], [3, 4], [5, 6]] },
  "qr": { "q": [[...]], "r": [[...]] },
  "statistics": {
    "max": 0.8970852271,
    "min": -7.4373574416,
    "average": -0.8812371694,
    "sum": -13.2185575408,
    "totalElements": 15,
    "isAnyDiagonal": false,
    "diagonalByMatrix": { "Q": false, "R": false }
  }
}
```

### `POST /api/v1/matrices/rotate` *(protegido con JWT)*
Rota la matriz 90° (`"clockwise"` por defecto o `"counterclockwise"`) y consulta las estadísticas a `node-stats-api`.

## Ejecución y Pruebas

```bash
go mod tidy
go test -v ./...
go run ./cmd/server
```
