import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
 
	_ "github.com/mattn/go-sqlite3"
)

// Estructuras para representar los datos de géneros, artistas, álbumes y canciones
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

// Main
func main() {
	loadMusic()

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