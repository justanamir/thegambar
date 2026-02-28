package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Photographer struct {
	ID        int
	Name      string
	Specialty string
	City      string
	Bio       string
	Email     string
	WhatsApp  string
	Website   string
}

var templates = template.Must(template.ParseGlob("web/templates/*.html"))
var db *sql.DB

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Verify connection is actually alive
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Connected to database ✓")

	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/photographer/", profileHandler)

	fmt.Println("thegambar running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, specialty, city FROM photographers ORDER BY id")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("homepageHandler query error:", err)
		return
	}
	defer rows.Close()

	var photographers []Photographer
	for rows.Next() {
		var p Photographer
		if err := rows.Scan(&p.ID, &p.Name, &p.Specialty, &p.City); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Println("homepageHandler scan error:", err)
			return
		}
		photographers = append(photographers, p)
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

	var p Photographer
	err = db.QueryRow(
		"SELECT id, name, specialty, city, bio, email, whatsapp, website FROM photographers WHERE id = $1",
		id,
	).Scan(&p.ID, &p.Name, &p.Specialty, &p.City, &p.Bio, &p.Email, &p.WhatsApp, &p.Website)

	if err == sql.ErrNoRows {
		http.Error(w, "Photographer not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("profileHandler query error:", err)
		return
	}

	templates.ExecuteTemplate(w, "profile.html", p)
}
