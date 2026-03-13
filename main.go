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
	http.HandleFunc("/api/genres/", genresHandler)

	log.Println("Server running on :24347")
	log.Fatal(http.ListenAndServe(":24347", nil))
}

// Función para conectar a la base de datos
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

// Extrae el ID del path: /api/genres/1
func extractID(r *http.Request, prefix string) (int, bool) {
	part := strings.TrimPrefix(r.URL.Path, prefix)
	part = strings.Trim(part, "/")
	if part == "" {
		return 0, false
	}
	id, err := strconv.Atoi(part)
	return id, err == nil && id > 0
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Message{Message: "pong"})
}

// Handler de generos, maneja GET, POST, PUT, PATCH, DELETE
func genresHandler(w http.ResponseWriter, r *http.Request) {

	id, hasID := extractID(r, "/api/genres/")

	if hasID {
		switch r.Method {
		case http.MethodGet:
			handleGetGenreByID(w, id)
		case http.MethodPut:
			handleUpdateGenre(w, r, id)
		case http.MethodPatch:
			handlePatchGenre(w, r, id)
		case http.MethodDelete:
			handleDeleteGenre(w, id)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetGenres(w, r)
	case http.MethodPost:
		handleCreateGenre(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GET 
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

// GET 
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

// POST 
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

// PUT 
func handleUpdateGenre(w http.ResponseWriter, r *http.Request, id int) {
	var g Genre
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(g.Name) == "" || strings.TrimSpace(g.Origin) == "" ||
		strings.TrimSpace(g.Mood) == "" || strings.TrimSpace(g.Tempo) == "" ||
		strings.TrimSpace(g.Description) == "" || g.Decade == 0 {
		http.Error(w, "All fields are required for PUT: name, origin, decade, mood, tempo, description", http.StatusBadRequest)
		return
	}

	res, err := db.Exec(
		"UPDATE genres SET name=?, origin=?, decade=?, mood=?, tempo=?, description=? WHERE id=?",
		g.Name, g.Origin, g.Decade, g.Mood, g.Tempo, g.Description, id,
	)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		http.Error(w, "Genre not found", http.StatusNotFound)
		return
	}
	g.ID = id
	writeJSON(w, http.StatusOK, g)
}

// PATCH 
func handlePatchGenre(w http.ResponseWriter, r *http.Request, id int) {
	var fields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if len(fields) == 0 {
		http.Error(w, "No fields provided", http.StatusBadRequest)
		return
	}

	allowed := map[string]bool{"name": true, "origin": true, "decade": true, "mood": true, "tempo": true, "description": true}
	var sets []string
	var args []interface{}
	for k, v := range fields {
		if !allowed[k] {
			http.Error(w, "Field '"+k+"' cannot be patched", http.StatusBadRequest)
			return
		}
		sets = append(sets, k+" = ?")
		args = append(args, v)
	}
	args = append(args, id)

	res, err := db.Exec("UPDATE genres SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		http.Error(w, "Genre not found", http.StatusNotFound)
		return
	}
	handleGetGenreByID(w, id)
}

// DELETE 
func handleDeleteGenre(w http.ResponseWriter, id int) {
	res, _ := db.Exec("DELETE FROM genres WHERE id = ?", id)
	n, _ := res.RowsAffected()
	if n == 0 {
		http.Error(w, "Genre not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, Message{Message: "Genre deleted"})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}