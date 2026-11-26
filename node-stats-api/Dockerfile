# Etapa 1: Compilación de TypeScript
FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY tsconfig.json ./
COPY src ./src
RUN npm run build

# Etapa 2: Imagen final liviana para producción
FROM node:22-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production

COPY package*.json ./
RUN npm ci --omit=dev && npm cache clean --force

COPY --from=builder /app/dist ./dist

# Ejecutar como usuario no-root por seguridad
USER node

EXPOSE 3000

CMD ["node", "dist/server.js"]
