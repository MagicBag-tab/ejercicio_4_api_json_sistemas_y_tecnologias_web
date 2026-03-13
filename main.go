package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Estructuras para representar los datos de géneros
type Genre struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Origin      string `json:"origin"`
	Decade      int    `json:"decade"`
	Mood        string `json:"mood"`
	Tempo       string `json:"tempo"`
	Description string `json:"description"`
}

type Message struct {
	Message string `json:"message"`
}

var db *sql.DB

// Main
func main() {
	loadMusic()
	defer db.Close()

	http.HandleFunc("/api/ping", pingHandler)
	http.HandleFunc("/api/genres", genresHandler)

	log.Println("Server running on :24347")
	log.Fatal(http.ListenAndServe(":80", nil))
}

// Función para obtener la data del db de Music
func loadMusic() {
	var err error
	db, err = sql.Open("sqlite3", "./data/music.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	log.Println("Connected to music.db")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Message{Message: "pong"})
}

// Handler para /api/genres
func genresHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetGenres(w, r)
	case http.MethodPost:
		handleCreateGenre(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Función para manejar GET /api/genres y GET /api/genres?id=1
func handleGetGenres(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}
		handleGetGenreByID(w, id)
		return
	}

	rows, err := db.Query("SELECT id, name, origin, decade, mood, tempo, description FROM genres ORDER BY id")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var genres []Genre
	for rows.Next() {
		var g Genre
		rows.Scan(&g.ID, &g.Name, &g.Origin, &g.Decade, &g.Mood, &g.Tempo, &g.Description)
		genres = append(genres, g)
	}
	writeJSON(w, http.StatusOK, genres)
}

// Función para manejar GET /api/genres?id=1 (buscar por ID)
func handleGetGenreByID(w http.ResponseWriter, id int) {
	var g Genre
	err := db.QueryRow("SELECT id, name, origin, decade, mood, tempo, description FROM genres WHERE id = ?", id).
		Scan(&g.ID, &g.Name, &g.Origin, &g.Decade, &g.Mood, &g.Tempo, &g.Description)
	if err == sql.ErrNoRows {
		http.Error(w, "Genre not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

// Función para crear un nuevo género con POST /api/genres
func handleCreateGenre(w http.ResponseWriter, r *http.Request) {
	var g Genre
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(g.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(g.Origin) == "" {
		http.Error(w, "origin is required", http.StatusBadRequest)
		return
	}
	if g.Decade < 1600 || g.Decade > 2030 {
		http.Error(w, "decade must be between 1600 and 2030", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(g.Mood) == "" {
		http.Error(w, "mood is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(g.Tempo) == "" {
		http.Error(w, "tempo is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(g.Description) == "" {
		http.Error(w, "description is required", http.StatusBadRequest)
		return
	}

	res, err := db.Exec(
		"INSERT INTO genres (name, origin, decade, mood, tempo, description) VALUES (?, ?, ?, ?, ?, ?)",
		g.Name, g.Origin, g.Decade, g.Mood, g.Tempo, g.Description,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			http.Error(w, "Genre name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	g.ID = int(id)
	writeJSON(w, http.StatusCreated, g)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}