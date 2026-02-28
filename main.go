package main

import "fmt"

// This is your Photographer blueprint.
// Every photographer on thegambar will be one of these.
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

func main() {
	// A slice — think of it as your directory listing.
	// Right now it's fake data. Later this comes from a real database.
	photographers := []Photographer{
		{ID: 1, Name: "Amir Hamzah", Specialty: "Wedding", City: "Kuala Lumpur", Email: "amir@email.com", WhatsApp: "+60123456789"},
		{ID: 2, Name: "Sara Lim", Specialty: "Street", City: "Penang", Email: "sara@email.com", Website: "saralim.com"},
		{ID: 3, Name: "Razif Osman", Specialty: "Commercial", City: "Johor Bahru", WhatsApp: "+60198765432"},
	}

	// Loop through every photographer and describe them
	for _, p := range photographers {
		describe(p)
		fmt.Println("---")
	}
}

// describe takes one Photographer and prints their summary.
// Notice Razif has no email — that's fine, it'll just print empty for now.
func describe(p Photographer) {
	fmt.Printf("%s | %s | %s\n", p.Name, p.Specialty, p.City)
	printContact(p)
}

func printContact(p Photographer) {
	if p.Email != "" {
		fmt.Printf("Email: %s\n", p.Email)
	}
	if p.WhatsApp != "" {
		fmt.Printf("WhatsApp: %s\n", p.WhatsApp)
	}
	if p.Website != "" {
		fmt.Printf("Website: %s\n", p.Website)
	}
}
