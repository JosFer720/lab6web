## Endpoints
- **GET** `/api/matches` Recupera la lista de todos los partidos almacenados en la base de datos.
- **GET** `/api/matches/:id` Obtiene los detalles de un partido específico a través de su ID.
- **POST** `/api/matches` Crea un nuevo partido proporcionando los equipos y la fecha del partido.
- **PUT** `/api/matches/:id` Actualiza los datos de un partido existente identificado por su ID.
- **DELETE** `/api/matches/:id` Elimina un partido específico de la base de datos según su ID.

## Base de Datos
La API utiliza SQLite para almacenar los datos de los partidos en un archivo llamado `lab6.db`. La estructura de la tabla `matches`:

```sql
CREATE TABLE matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    homeTeam TEXT NOT NULL,
    awayTeam TEXT NOT NULL,
    matchDate TEXT NOT NULL
);
```
Cada partido tiene un ID único, los nombres de los equipos local y visitante, y la fecha del partido.