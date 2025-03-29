package main

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           API de Partidos de Fútbol
// @version         1.0
// @description     Sistema de gestión de partidos de fútbol con estadísticas
// @termsOfService  http://swagger.io/terms/

// @contact.name   Soporte API
// @contact.url    http://www.swagger.io/support
// @contact.email  soporte@futbolapi.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth
type Match struct {
	ID          int    `json:"id"`
	HomeTeam    string `json:"homeTeam" binding:"required"`
	AwayTeam    string `json:"awayTeam" binding:"required"`
	MatchDate   string `json:"matchDate" binding:"required"`
	HomeGoals   int    `json:"homeGoals"`
	AwayGoals   int    `json:"awayGoals"`
	YellowCards int    `json:"yellowCards"`
	RedCards    int    `json:"redCards"`
	ExtraTime   int    `json:"extraTime"`
}

var db *sql.DB
var dbPath string

func main() {
	setupDatabase()
	router := setupRouter()
	router.Run(":8080")
}

func setupDatabase() {
	var err error
	dbPath = getDBPath()

	db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		panic("Error al conectar a SQLite: " + err.Error())
	}

	var tableExists bool
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='matches'").Scan(&tableExists)
	if err != nil {
		panic("Error al verificar tablas existentes: " + err.Error())
	}

	if !tableExists {
		_, err = db.Exec(`
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
		)`)
		if err != nil {
			panic("Error al crear tabla: " + err.Error())
		}

		_, err = db.Exec(`
		INSERT INTO matches (homeTeam, awayTeam, matchDate) VALUES
		('Real Madrid', 'Barcelona', '2024-06-01'),
		('Atlético de Madrid', 'Sevilla', '2024-06-02')`)
		if err != nil {
			panic("Error al insertar datos iniciales: " + err.Error())
		}
	}
}

func getDBPath() string {
	if path := os.Getenv("DB_PATH"); path != "" {
		return path
	}

	exePath, err := os.Executable()
	if err != nil {
		return "lab6.db"
	}
	exeDir := filepath.Dir(exePath)

	if _, err := os.Stat(exeDir); os.IsPermission(err) {
		return "lab6.db"
	}

	return filepath.Join(exeDir, "lab6.db")
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	// Configuración CORS (mantén tu configuración actual)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Total-Count")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.Static("/docs", "./docs")

	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/docs/swagger.json"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	api := router.Group("/api")
	{
		api.GET("/matches", listMatches)
		api.GET("/matches/:id", getMatchByID)
		api.POST("/matches", createMatch)
		api.PUT("/matches/:id", updateMatch)
		api.DELETE("/matches/:id", deleteMatch)
		api.PATCH("/matches/:id/goals", updateGoals)
		api.PATCH("/matches/:id/yellowcards", updateYellowCards)
		api.PATCH("/matches/:id/redcards", updateRedCards)
		api.PATCH("/matches/:id/extratime", updateExtraTime)
	}

	return router
}

// @Summary Listar todos los partidos
// @Tags Partidos
// @Produce json
// @Success 200 {array} Match
// @Router /api/matches [get]
func listMatches(c *gin.Context) {
	rows, err := db.Query("SELECT id, homeTeam, awayTeam, matchDate, homeGoals, awayGoals, yellowCards, redCards, extraTime FROM matches")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar partidos: " + err.Error()})
		return
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		err := rows.Scan(
			&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate,
			&m.HomeGoals, &m.AwayGoals, &m.YellowCards, &m.RedCards, &m.ExtraTime,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer datos: " + err.Error()})
			return
		}
		matches = append(matches, m)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error después de leer filas: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, matches)
}

// @Summary Obtener un partido por ID
// @Tags Partidos
// @Produce json
// @Param id path int true "ID del partido"
// @Success 200 {object} Match
// @Failure 404 {object} object
// @Router /api/matches/{id} [get]
func getMatchByID(c *gin.Context) {
	id := c.Param("id")
	var match Match

	err := db.QueryRow(`
		SELECT id, homeTeam, awayTeam, matchDate, 
		homeGoals, awayGoals, yellowCards, redCards, extraTime 
		FROM matches WHERE id = ?`, id).Scan(
		&match.ID, &match.HomeTeam, &match.AwayTeam, &match.MatchDate,
		&match.HomeGoals, &match.AwayGoals, &match.YellowCards, &match.RedCards, &match.ExtraTime,
	)

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

// @Summary Actualizar goles del partido
// @Tags Estadísticas
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param homeGoals query int false "Goles del equipo local"
// @Param awayGoals query int false "Goles del equipo visitante"
// @Success 200 {object} Match
// @Router /api/matches/{id}/goals [patch]
func updateGoals(c *gin.Context) {
	id := c.Param("id")
	homeGoals, _ := strconv.Atoi(c.DefaultQuery("homeGoals", "0"))
	awayGoals, _ := strconv.Atoi(c.DefaultQuery("awayGoals", "0"))

	_, err := db.Exec(
		"UPDATE matches SET homeGoals = ?, awayGoals = ? WHERE id = ?",
		homeGoals, awayGoals, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var match Match
	db.QueryRow("SELECT * FROM matches WHERE id = ?", id).Scan(
		&match.ID, &match.HomeTeam, &match.AwayTeam, &match.MatchDate,
		&match.HomeGoals, &match.AwayGoals, &match.YellowCards, &match.RedCards, &match.ExtraTime)

	c.JSON(http.StatusOK, match)
}

// Implementaciones similares para:
// updateYellowCards, updateRedCards, updateExtraTime
// createMatch, updateMatch, deleteMatch

// @Summary Registrar tarjeta amarilla
// @Tags Estadísticas
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param count query int false "Cantidad de tarjetas a añadir" default(1)
// @Success 200 {object} Match
// @Router /api/matches/{id}/yellowcards [patch]
func updateYellowCards(c *gin.Context) {
	id := c.Param("id")
	count, _ := strconv.Atoi(c.DefaultQuery("count", "1"))

	_, err := db.Exec(
		"UPDATE matches SET yellowCards = yellowCards + ? WHERE id = ?",
		count, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var match Match
	db.QueryRow("SELECT * FROM matches WHERE id = ?", id).Scan(
		&match.ID, &match.HomeTeam, &match.AwayTeam, &match.MatchDate,
		&match.HomeGoals, &match.AwayGoals, &match.YellowCards, &match.RedCards, &match.ExtraTime)

	c.JSON(http.StatusOK, match)
}

// @Summary Registrar tarjeta roja
// @Tags Estadísticas
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param count query int false "Cantidad de tarjetas a añadir" default(1)
// @Success 200 {object} Match
// @Router /api/matches/{id}/redcards [patch]
func updateRedCards(c *gin.Context) {
	id := c.Param("id")
	count, _ := strconv.Atoi(c.DefaultQuery("count", "1"))

	_, err := db.Exec(
		"UPDATE matches SET redCards = redCards + ? WHERE id = ?",
		count, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var match Match
	db.QueryRow("SELECT * FROM matches WHERE id = ?", id).Scan(
		&match.ID, &match.HomeTeam, &match.AwayTeam, &match.MatchDate,
		&match.HomeGoals, &match.AwayGoals, &match.YellowCards, &match.RedCards, &match.ExtraTime)

	c.JSON(http.StatusOK, match)
}

// @Summary Registrar tiempo extra
// @Tags Estadísticas
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param minutes query int false "Minutos de tiempo extra" default(2)
// @Success 200 {object} Match
// @Router /api/matches/{id}/extratime [patch]
func updateExtraTime(c *gin.Context) {
	id := c.Param("id")
	minutes, _ := strconv.Atoi(c.DefaultQuery("minutes", "2"))

	_, err := db.Exec(
		"UPDATE matches SET extraTime = ? WHERE id = ?",
		minutes, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var match Match
	db.QueryRow("SELECT * FROM matches WHERE id = ?", id).Scan(
		&match.ID, &match.HomeTeam, &match.AwayTeam, &match.MatchDate,
		&match.HomeGoals, &match.AwayGoals, &match.YellowCards, &match.RedCards, &match.ExtraTime)

	c.JSON(http.StatusOK, match)
}

// @Summary Crear nuevo partido
// @Tags Partidos
// @Accept json
// @Produce json
// @Param match body Match true "Datos del partido"
// @Success 201 {object} Match
// @Router /api/matches [post]
func createMatch(c *gin.Context) {
	var newMatch Match
	if err := c.ShouldBindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	result, err := db.Exec(
		"INSERT INTO matches (homeTeam, awayTeam, matchDate) VALUES (?, ?, ?)",
		newMatch.HomeTeam, newMatch.AwayTeam, newMatch.MatchDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear partido: " + err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	newMatch.ID = int(id)
	c.JSON(http.StatusCreated, newMatch)
}

// @Summary Actualizar partido
// @Tags Partidos
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param match body Match true "Datos actualizados"
// @Success 200 {object} Match
// @Router /api/matches/{id} [put]
func updateMatch(c *gin.Context) {
	id := c.Param("id")
	var updated Match

	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	_, err := db.Exec(
		"UPDATE matches SET homeTeam = ?, awayTeam = ?, matchDate = ? WHERE id = ?",
		updated.HomeTeam, updated.AwayTeam, updated.MatchDate, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar: " + err.Error()})
		return
	}

	updated.ID, _ = strconv.Atoi(id)
	c.JSON(http.StatusOK, updated)
}

// @Summary Eliminar partido
// @Tags Partidos
// @Produce json
// @Param id path int true "ID del partido"
// @Success 200 {object} object
// @Router /api/matches/{id} [delete]
func deleteMatch(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM matches WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partido eliminado correctamente"})
}
