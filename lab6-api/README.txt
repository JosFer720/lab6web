# lab6-api

API para gestión de partidos de fútbol.

## Documentación

- **Swagger UI**: [http://localhost:8080/swagger/index.html#/]
- **Colección Postman**: [https://fernando-1614404.postman.co/workspace/Fernando's-Workspace~16733a76-0da6-4135-8247-30881953186e/request/43575536-af6fa617-bfd3-4458-80b5-693ad3d9c8d0?action=share&creator=43575536&ctx=documentation&active-environment=43575536-8e0a080d-ae32-4fa1-aa5e-bb4053b0c033]

## Instalación

1. Clonar el repositorio
2. Ejecutar: `swag init -g main.go`
3. Ejecutar: `go run main.go`

## Endpoints

Ver documentación completa en Swagger UI o en el archivo [llms.txt](llms.txt).

## Base de Datos

SQLite con el archivo `lab6.db`. Estructura:

```sql
CREATE TABLE matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    homeTeam TEXT NOT NULL,
    awayTeam TEXT NOT NULL,
    matchDate TEXT NOT NULL,
    homeGoals INTEGER DEFAULT 0,
    awayGoals INTEGER DEFAULT 0,
    yellowCards INTEGER DEFAULT 0,
    redCards INTEGER DEFAULT 0,
    extraTime INTEGER DEFAULT 0
);