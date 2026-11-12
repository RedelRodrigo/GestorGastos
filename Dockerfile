# ==========================================
# ETAPA 1: Construir el Frontend (Bun + React)
# ==========================================
FROM oven/bun:1 AS frontend-builder
WORKDIR /app/frontend

# Copiamos package.json y tu archivo bun.lock actual
COPY frontend/package.json frontend/bun.lock ./
RUN bun install

# Copiamos el resto del código del frontend y compilamos
COPY frontend/ ./
RUN bun run build

# ==========================================
# ETAPA 2: Construir el Backend (Go)
# ==========================================
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app/backend

# Copiamos módulos de Go
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copiamos el código del backend
COPY backend/ ./

# Copiamos los archivos estáticos generados en la Etapa 1 dentro del backend
COPY --from=frontend-builder /app/frontend/dist ./static

# Compilamos el binario de Go de manera estática y optimizada
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api

# ==========================================
# ETAPA 3: Imagen Final de Producción (Alpine)
# ==========================================
FROM alpine:3.19
WORKDIR /app

# Instalamos certificados de seguridad SSL necesarios para HTTPS externo (Turso)
RUN apk --no-cache add ca-certificates

# Copiamos el binario compilado y los estáticos desde las etapas anteriores
COPY --from=backend-builder /app/bin/api /app/api
COPY --from=backend-builder /app/backend/static /app/static

# Definimos variables de entorno por defecto
ENV PORT=8080
ENV STATIC_DIR=/app/static

EXPOSE 8080

# Comando para arrancar el servidor
CMD ["./api"]