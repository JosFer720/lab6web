# Dockerfile para lab6-api con SQLite persistente (Alpine + CGO)
# ========= ETAPA DE CONSTRUCCIÓN =========
FROM golang:1.24-alpine AS builder

# Instalar dependencias de compilación C para SQLite
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# 1. Copiar dependencias primero para mejor caché
COPY go.mod go.sum ./
RUN go mod download

# 2. Copiar todo el código fuente
COPY . .

# 3. Construir la aplicación con soporte CGO
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /lab6-api

# ========= ETAPA DE EJECUCIÓN =========
FROM alpine:3.18

WORKDIR /app

# 1. Instalar SQLite y crear directorio para datos
RUN apk add --no-cache \
    sqlite \
    ca-certificates && \
    mkdir -p /data && \
    chmod 777 /data

# 2. Copiar artefactos desde builder
COPY --from=builder /lab6-api /app/
COPY --from=builder /app/LaLigaTracker.html /app/

# 3. Configurar variables de entorno
ENV DB_PATH=/data/lab6.db \
    GIN_MODE=release

# 4. Volumen para persistencia
VOLUME /data

# 5. Puerto expuesto
EXPOSE 8080

# 6. Comando de inicio
CMD ["/app/lab6-api"]