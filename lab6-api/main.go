package main

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Match struct {
	ID        int    `json:"id"`
	HomeTeam  string `json:"homeTeam" binding:"required"`
	AwayTeam  string `json:"awayTeam" binding:"required"`
	MatchDate string `json:"matchDate" binding:"required"`
}

var db *sql.DB

func main() {
	// Configuración inicial
	setupDatabase()
	router := setupRouter()
	router.Run(":8080")
}

func setupDatabase() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "lab6.db"
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic("Error al conectar a SQLite: " + err.Error())
	}

	// Crear tabla con datos iniciales
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS matches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		homeTeam TEXT NOT NULL,
		awayTeam TEXT NOT NULL,
		matchDate TEXT NOT NULL
	);

	INSERT INTO matches (homeTeam, awayTeam, matchDate) VALUES
	('Real Madrid', 'Barcelona', '2024-06-01'),
	('Atlético de Madrid', 'Sevilla', '2024-06-02')`)

	if err != nil {
		panic("Error al crear tabla: " + err.Error())
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Rutas
	router.GET("/api/matches", listMatches)
	router.GET("/api/matches/:id", getMatchByID)
	router.POST("/api/matches", createMatch)
	router.PUT("/api/matches/:id", updateMatch)
	router.DELETE("/api/matches/:id", deleteMatch)

	return router
}

// Handlers
func listMatches(c *gin.Context) {
	rows, err := db.Query("SELECT id, homeTeam, awayTeam, matchDate FROM matches")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar partidos"})
		return
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer datos"})
			return
		}
		matches = append(matches, m)
	}
	c.JSON(http.StatusOK, matches)
}

func getMatchByID(c *gin.Context) {
	id := c.Param("id")
	var match Match

	err := db.QueryRow("SELECT id, homeTeam, awayTeam, matchDate FROM matches WHERE id = ?", id).Scan(
		&match.ID, &match.HomeTeam, &match.AwayTeam, &match.MatchDate)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la base de datos"})
		}
		return
	}
	c.JSON(http.StatusOK, match)
}

func createMatch(c *gin.Context) {
	var newMatch Match
	if err := c.ShouldBindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	result, err := db.Exec(
		"INSERT INTO matches (homeTeam, awayTeam, matchDate) VALUES (?, ?, ?)",
		newMatch.HomeTeam, newMatch.AwayTeam, newMatch.MatchDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear partido"})
		return
	}

	id, _ := result.LastInsertId()
	newMatch.ID = int(id)
	c.JSON(http.StatusCreated, newMatch)
}

func updateMatch(c *gin.Context) {
	id := c.Param("id")
	var updated Match

	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Convertir ID a número
	matchID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	_, err = db.Exec(
		"UPDATE matches SET homeTeam = ?, awayTeam = ?, matchDate = ? WHERE id = ?",
		updated.HomeTeam, updated.AwayTeam, updated.MatchDate, matchID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
		return
	}

	updated.ID = matchID
	c.JSON(http.StatusOK, updated)
}

func deleteMatch(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM matches WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partido eliminado correctamente"})
}
