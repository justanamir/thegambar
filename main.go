package main

import (
	"fmt"
	"net/http"
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

// Your fake data lives here for now
var photographers = []Photographer{
	{ID: 1, Name: "Amir Hamzah", Specialty: "Wedding", City: "Kuala Lumpur", Email: "amir@email.com", WhatsApp: "+60123456789"},
	{ID: 2, Name: "Sara Lim", Specialty: "Street", City: "Penang", Email: "sara@email.com", Website: "saralim.com"},
	{ID: 3, Name: "Razif Osman", Specialty: "Commercial", City: "Johor Bahru", WhatsApp: "+60198765432"},
}

func main() {
	// When someone visits /, run the homepageHandler function
	http.HandleFunc("/", homepageHandler)

	// Start listening on port 8080
	fmt.Println("thegambar running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Welcome to thegambar</h1>")
	fmt.Fprintf(w, "<p>%d photographers and counting.</p>", len(photographers))
	fmt.Fprintf(w, "<hr>")

	for _, p := range photographers {
		fmt.Fprintf(w, "<div style='margin-bottom:20px'>")
		fmt.Fprintf(w, "<h2>%s</h2>", p.Name)
		fmt.Fprintf(w, "<p>%s · %s</p>", p.Specialty, p.City)
		printContactHTML(w, p)
		fmt.Fprintf(w, "</div>")
	}
}

func printContactHTML(w http.ResponseWriter, p Photographer) {
	if p.Email != "" {
		fmt.Fprintf(w, "<p>Email: %s</p>", p.Email)
	}
	if p.WhatsApp != "" {
		fmt.Fprintf(w, "<p>WhatsApp: %s</p>", p.WhatsApp)
	}
	if p.Website != "" {
		fmt.Fprintf(w, "<p>Website: <a href='https://%s'>%s</a></p>", p.Website, p.Website)
	}
}
