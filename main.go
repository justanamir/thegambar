package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"thegambar/internal/db"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var templates = template.Must(template.ParseGlob("web/templates/*.html"))
var queries *db.Queries

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer conn.Close()

	if err = conn.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Connected to database ✓")

	// Hand the connection to sqlc
	queries = db.New(conn)

	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/photographer/", profileHandler)

	fmt.Println("thegambar running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	photographers, err := queries.ListPhotographers(r.Context())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("homepageHandler error:", err)
		return
	}

	templates.ExecuteTemplate(w, "home.html", photographers)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/photographer/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid photographer ID", http.StatusBadRequest)
		return
	}

	photographer, err := queries.GetPhotographer(r.Context(), int32(id))
	if err == sql.ErrNoRows {
		http.Error(w, "Photographer not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("profileHandler error:", err)
		return
	}

	templates.ExecuteTemplate(w, "profile.html", photographer)
}
