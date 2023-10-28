package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"os"
	"time"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/gin-gonic/gin"
)

type Note struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Important bool      `json:"important"`
	Date      time.Time `json:"date"`
}

var notes = []Note{
	{ID: 1, Content: "🥸 HTML is easy", Important: false, Date: time.Now()},
	{ID: 2, Content: "🌎 Browser can execute only JavaScript", Important: false, Date: time.Now().Add(-24 * time.Hour)},
	{ID: 3, Content: "🌤️ GET and POST are the most important methods of HTTP protocol", Important: true, Date: time.Now().Add(-48 * time.Hour)},
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "🏠 Página principal!!")
}

func Hello(ctx *gin.Context) {
	name := ctx.Param("name")
	ctx.JSON(200, gin.H{"message": fmt.Sprintf("Hello %s!", name)})
}
func Engine() *gin.Engine {
	engine := gin.Default()
	engine.GET("/hello/:name", Hello)
	return engine
}

func GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(notes)
}

func GetNoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for _, note := range notes {
		if note.ID == id {
			json.NewEncoder(w).Encode(note)
			return
		}
	}
	http.NotFound(w, r)
}

func UpdateNoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.NotFound(w, r)
		return
	}

	var note *Note
	for i := range notes {
		if notes[i].ID == id {
			note = &notes[i]
			break
		}
	}

	log.Println(note)

	if note == nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(notes)
}

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(note)

	if note.Content == "" {
		http.Error(w, "Note cannot be empty", http.StatusBadRequest)
		return
	}

	var maxID int = 0
	for _, note := range notes {
		if note.ID > maxID {
			maxID = note.ID
		}
	}

	newID := maxID + 1
	note.ID = newID

	timeNow := time.Now()
	layout := "02/01/2006 15:04:05"
	formattedTime := timeNow.Format(layout)

	date, err := time.Parse(layout, formattedTime)
	note.Date = date

	notes = append(notes, note)
	log.Println(notes)

	json.NewEncoder(w).Encode(note)
}

var PORT = os.Getenv("PORT")

func main() {
	router := mux.NewRouter()

	// Endpoints
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/api/notes", GetAllNotesHandler).Methods("GET")
	router.HandleFunc("/api/notes/{id}", GetNoteByIDHandler).Methods("GET")
	router.HandleFunc("/api/notes/{id}", UpdateNoteByIDHandler).Methods("PUT")
	router.HandleFunc("/api/notes", CreateNoteHandler).Methods("POST")

	// Cors
	c := cors.AllowAll()

	// Handler
	handler := c.Handler(router)

	// Start server
	if PORT == "" {
		PORT = "5000"
	}
	log.Println("Service running on port :" + PORT)
	err := http.ListenAndServe(":"+PORT, handler)
	if err != nil {
		log.Fatal(err)
	}
}
