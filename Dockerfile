# ==========================================
# ETAPA 1: Construir el Frontend (Bun + React)
# ==========================================
FROM oven/bun:1 AS frontend-builder
WORKDIR /app

# Copiamos package.json y bun.lock desde la carpeta frontend local
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

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

# Aquí copiamos desde /app/dist (que es donde bun build deja los archivos si el WORKDIR era /app)
COPY --from=frontend-builder /app/dist ./static

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api

# ==========================================
# ETAPA 3: Imagen Final de Producción (Alpine)
# ==========================================
FROM alpine:3.19
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=backend-builder /app/bin/api /app/api
COPY --from=backend-builder /app/backend/static /app/static

ENV PORT=8080
ENV STATIC_DIR=/app/static

EXPOSE 8080

CMD ["./api"]