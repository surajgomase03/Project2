package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)

	// Serve static files (if needed)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Handler for home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Welcome to Go Website",
		"Name":  "Go Developer",
	}
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// Handler for about page
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "About Us",
		"Content": "This is a simple Go website for learning and practice.",
	}
	tmpl, err := template.ParseFiles("templates/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// Handler for contact page
func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle form submission
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		// Here you can process the form data
		fmt.Fprintf(w, "Thank you %s! Your message has been received.\nEmail: %s\nMessage: %s", name, email, message)
		return
	}

	tmpl, err := template.ParseFiles("templates/contact.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}