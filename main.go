import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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
 
type Artist struct {
	ID      int    `json:"id"`
	GenreID int    `json:"genre_id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Formed  int    `json:"formed"`
	Active  bool   `json:"active"`
	Bio     string `json:"bio"`
}
 
type Album struct {
	ID       int    `json:"id"`
	ArtistID int    `json:"artist_id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Label    string `json:"label"`
	Tracks   int    `json:"tracks"`
}
 
type Song struct {
	ID        int     `json:"id"`
	AlbumID   int     `json:"album_id"`
	Title     string  `json:"title"`
	Duration  string  `json:"duration"`
	TrackNo   int     `json:"track_no"`
	Featuring *string `json:"featuring,omitempty"` //omitempty en go representa un campo que puede ser vacío
	Explicit  bool    `json:"explicit"`
}

type Message struct {
	Message string `json:"message"`
}