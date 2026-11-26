# Etapa 1: Compilación del binario Go
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .
RUN if [ ! -d vendor ]; then go mod download; fi && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# Etapa 2: Imagen final minimalista y segura
FROM alpine:3.20

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/server ./server

USER appuser

EXPOSE 8080

CMD ["./server"]
