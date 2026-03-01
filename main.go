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
	"thegambar/internal/storage"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var templates = template.Must(template.ParseGlob("web/templates/*.html"))
var queries *db.Queries
var r2 *storage.R2Client

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r2 = storage.NewR2Client(
		os.Getenv("R2_ACCOUNT_ID"),
		os.Getenv("R2_ACCESS_KEY_ID"),
		os.Getenv("R2_SECRET_ACCESS_KEY"),
		os.Getenv("R2_BUCKET_NAME"),
		os.Getenv("R2_PUBLIC_URL"),
	)
	fmt.Println("R2 storage client initialised ✓")

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
	http.HandleFunc("/photos/", photoUploadHandler)
	http.HandleFunc("/join", joinHandler)

	fmt.Println("thegambar running on http://localhost:8080")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
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

func joinHandler(w http.ResponseWriter, r *http.Request) {
	// GET → show the empty form
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "join.html", nil)
		return
	}

	// POST → process the submission
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ParseForm makes r.FormValue() work.
	// Without this, form fields come back as empty strings.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Pull values from the form.
	// r.FormValue trims nothing — you get exactly what the user typed.
	name := strings.TrimSpace(r.FormValue("name"))
	specialty := strings.TrimSpace(r.FormValue("specialty"))
	city := strings.TrimSpace(r.FormValue("city"))
	bio := strings.TrimSpace(r.FormValue("bio"))
	email := strings.TrimSpace(r.FormValue("email"))
	whatsapp := strings.TrimSpace(r.FormValue("whatsapp"))
	website := strings.TrimSpace(r.FormValue("website"))

	// --- Validation ---
	// Collect all errors so we can show them all at once,
	// not one-by-one (nothing more frustrating than whack-a-mole form errors).
	var formErrors []string

	if name == "" {
		formErrors = append(formErrors, "Name is required.")
	}
	if specialty == "" {
		formErrors = append(formErrors, "Specialty is required.")
	}
	if city == "" {
		formErrors = append(formErrors, "City is required.")
	}
	if email == "" && whatsapp == "" {
		formErrors = append(formErrors, "At least one contact method (email or WhatsApp) is required.")
	}

	// If there are errors, re-render the form with the errors AND the values
	// they already typed — never make users retype a full form.
	if len(formErrors) > 0 {
		templates.ExecuteTemplate(w, "join.html", map[string]any{
			"Errors":    formErrors,
			"Name":      name,
			"Specialty": specialty,
			"City":      city,
			"Bio":       bio,
			"Email":     email,
			"Whatsapp":  whatsapp,
			"Website":   website,
		})
		return
	}

	// --- Database insert ---
	// sqlc generated this function. It expects a params struct.
	// Nullable fields use sql.NullString — that's the Go way of saying "this might be empty".
	photographer, err := queries.InsertPhotographer(r.Context(), db.InsertPhotographerParams{
		Name:      name,
		Specialty: specialty,
		City:      city,
		Bio:       bio,
		Email:     email,
		Whatsapp:  whatsapp,
		Website:   website,
	})
	if err != nil {
		http.Error(w, "Failed to save photographer", http.StatusInternalServerError)
		log.Println("joinHandler insert error:", err)
		return
	}

	// Success — redirect to their new profile page.
	// 303 See Other is the correct redirect after a POST.
	// (Not 301 or 302 — those can cause browsers to re-POST on back button.)
	http.Redirect(w, r, fmt.Sprintf("/photographer/%d", photographer.ID), http.StatusSeeOther)
}

func photoUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from /photos/3
	idStr := strings.TrimPrefix(r.URL.Path, "/photos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// 10MB max — generous for two photos, protects your server
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	var avatarURL, coverURL string

	// Upload avatar if provided
	if avatarFile, avatarHeader, err := r.FormFile("avatar"); err == nil {
		defer avatarFile.Close()
		url, err := r2.UploadFile(r.Context(), avatarFile, avatarHeader)
		if err != nil {
			http.Error(w, "Avatar upload failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		avatarURL = url
	}

	// Upload cover if provided
	if coverFile, coverHeader, err := r.FormFile("cover"); err == nil {
		defer coverFile.Close()
		url, err := r2.UploadFile(r.Context(), coverFile, coverHeader)
		if err != nil {
			http.Error(w, "Cover upload failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		coverURL = url
	}

	// Save URLs to database
	_, err = queries.UpdatePhotographerPhotos(r.Context(), db.UpdatePhotographerPhotosParams{
		ID:        int32(id),
		AvatarUrl: avatarURL,
		CoverUrl:  coverURL,
	})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("photoUploadHandler db error:", err)
		return
	}

	// Back to their profile
	http.Redirect(w, r, fmt.Sprintf("/photographer/%d", id), http.StatusSeeOther)
}
