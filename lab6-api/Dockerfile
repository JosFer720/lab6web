# ========= ETAPA DE CONSTRUCCIÓN =========
FROM golang:1.24.1-alpine AS builder

# Instalar dependencias de compilación y herramientas para Swagger
RUN apk add --no-cache \
    gcc \
    musl-dev \
    git

# Instalar swag CLI para documentación
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Copiar módulos primero para mejor caché
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Generar documentación Swagger
RUN /go/bin/swag init -g main.go --output docs --parseDependency --parseInternal

# Construir la aplicación con soporte CGO
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /lab6-api

# ========= ETAPA DE EJECUCIÓN =========
FROM alpine:3.18

WORKDIR /app

# Instalar dependencias de runtime
RUN apk add --no-cache \
    sqlite \
    ca-certificates \
    tzdata && \
    mkdir -p /data && \
    chown nobody:nobody /data

# Copiar artefactos desde builder
COPY --from=builder /lab6-api /app/
COPY --from=builder /app/docs /app/docs/
COPY --from=builder /app/LaLigaTracker.html /app/

# Configurar variables de entorno
ENV DB_PATH=/data/lab6.db \
    GIN_MODE=release \
    TZ=UTC

# Configurar permisos
RUN chown -R nobody:nobody /app

# Volumen para persistencia
VOLUME /data

# Puerto expuesto
EXPOSE 8080

# Usuario no privilegiado
USER nobody

# Comando de inicio
CMD ["/app/lab6-api"]