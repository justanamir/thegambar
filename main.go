package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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

var photographers = []Photographer{
	{ID: 1, Name: "Amir Hamzah", Specialty: "Wedding", City: "Kuala Lumpur", Email: "amir@email.com", WhatsApp: "+60123456789"},
	{ID: 2, Name: "Sara Lim", Specialty: "Street", City: "Penang", Email: "sara@email.com", Website: "saralim.com"},
	{ID: 3, Name: "Razif Osman", Specialty: "Commercial", City: "Johor Bahru", WhatsApp: "+60198765432"},
}

// Parse all templates once at startup
var templates = template.Must(template.ParseGlob("web/templates/*.html"))

func main() {
	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/photographer/", profileHandler)

	fmt.Println("thegambar running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.html", photographers)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/photographer/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid photographer ID", http.StatusBadRequest)
		return
	}

	var found *Photographer
	for i := range photographers {
		if photographers[i].ID == id {
			found = &photographers[i]
			break
		}
	}

	if found == nil {
		http.Error(w, "Photographer not found", http.StatusNotFound)
		return
	}

	templates.ExecuteTemplate(w, "profile.html", found)
}
