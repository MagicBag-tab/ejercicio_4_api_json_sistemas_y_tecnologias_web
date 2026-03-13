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
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
	Decade int    `json:"decade"`
	Mood   string `json:"mood"`
	Tempo  string `json:"tempo"`
}

type Message struct {
	Message string `json:"message"`
}

var genres []Genre

var db *sql.DB


// Main
func main() {
	loadMusic()
	defer db.Close()

	http.HandleFunc("/api/ping", pingHandler)
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
			handleGetGenre(w, r)
		case http.MethodPut:
			handleUpdateGenre(w, r)
		case http.MethodPatch:
			handlePatchGenre(w, r)
		case http.MethodDelete:
			handleDeleteGenre(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}