package routes

import (
	"html/template"
	"net/http"

	"groupie-tracker/internal/handlers"
)

// NewRouter sets up the mux with all the necessary routes.
func NewRouter(tpl *template.Template) *http.ServeMux {
	mux := http.NewServeMux()

	// Home route with exact path checking.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(tpl)(w, r)
	})

	// Artist detail route.
	mux.HandleFunc("/artist/", handlers.DetailHandler(tpl))

	// API Route for fetching artists
	mux.HandleFunc("/api/artists", handlers.GetArtists)

	// About page section
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "about.html", nil)
	})

	// Trigger event route.
	mux.HandleFunc("/trigger-event", handlers.TriggerEventHandler)

	// Serve static files.
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
