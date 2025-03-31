package handlers

import (
	"encoding/json"
	"groupie-tracker/internal/data"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// HomeHandler handles the "/" route and renders the home page.
func HomeHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/" {
			handler404(tpl, w)
			return
		}

		if len(data.AllArtists) == 0 {
			log.Println("ERROR: No artist data available")
			http.Error(w, "No artist data available", http.StatusInternalServerError)
			return
		}

		// Sort artists by name, case insensitive.
		sortedArtists := make([]data.Artist, len(data.AllArtists))
		copy(sortedArtists, data.AllArtists)
		sort.Slice(sortedArtists, func(i, j int) bool {
			return strings.ToLower(sortedArtists[i].Name) < strings.ToLower(sortedArtists[j].Name)
		})

		// Only pass the initial 10 artists.
		initialCount := 10
		var initialArtists []data.Artist
		if len(sortedArtists) > initialCount {
			initialArtists = sortedArtists[:initialCount]
		} else {
			initialArtists = sortedArtists
		}

		// Render the homepage template with initialArtists.
		err := tpl.ExecuteTemplate(w, "home.html", initialArtists)
		if err != nil {
			log.Println("ERROR rendering template:", err)
			http.Error(w, "Internal Server Error while rendering index", http.StatusInternalServerError)
		}
	}
}

// render404 renders a custom 404 page.
func handler404(tpl *template.Template, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	if err := tpl.ExecuteTemplate(w, "404.html", nil); err != nil {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
	}
}

// GetArtists handles API requests for fetching artists with pagination.
func GetArtists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(data.AllArtists) == 0 {
		log.Println("ERROR: No artist data available")
		http.Error(w, "No artist data available", http.StatusInternalServerError)
		return
	}

	// Sort artists by name, case insensitive.
	sortedArtists := make([]data.Artist, len(data.AllArtists))
	copy(sortedArtists, data.AllArtists)
	sort.Slice(sortedArtists, func(i, j int) bool {
		return strings.ToLower(sortedArtists[i].Name) < strings.ToLower(sortedArtists[j].Name)
	})

	// Default values.
	offset := 0
	limit := 10

	// Parse offset if provided.
	if offStr := r.URL.Query().Get("offset"); offStr != "" {
		if off, err := strconv.Atoi(offStr); err == nil {
			offset = off
		}
	}

	// Parse limit if provided.
	if limStr := r.URL.Query().Get("limit"); limStr != "" {
		if lim, err := strconv.Atoi(limStr); err == nil {
			limit = lim
		}
	}

	// If offset is beyond available data, return an empty list.
	if offset >= len(sortedArtists) {
		json.NewEncoder(w).Encode([]data.Artist{})
		return
	}

	// Calculate the end index.
	end := offset + limit
	if end > len(sortedArtists) {
		end = len(sortedArtists)
	}

	paginatedArtists := sortedArtists[offset:end]
	err := json.NewEncoder(w).Encode(paginatedArtists)
	if err != nil {
		log.Println("ERROR encoding JSON:", err)
		http.Error(w, "Internal Server Error while encoding JSON", http.StatusInternalServerError)
	}
}
