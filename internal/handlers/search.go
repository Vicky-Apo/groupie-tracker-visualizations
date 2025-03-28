package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker-search-bar/internal/data"
)

type SearchResult struct {
	Value  string `json:"Value"`
	Type   string `json:"Type"`
	Artist string `json:"Artist"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("query"))
	filter := r.URL.Query().Get("filter") // all, artist, member, etc.

	if query == "" {
		json.NewEncoder(w).Encode([]SearchResult{})
		return
	}

	var results []SearchResult

	for _, a := range data.AllArtists {
		switch filter {
		case "artist", "all":
			if strings.Contains(strings.ToLower(a.Name), query) {
				results = append(results, SearchResult{Value: a.Name, Type: "Artist", Artist: a.Name})
			}
		}
		switch filter {
		case "member", "all":
			for _, m := range a.Members {
				if strings.Contains(strings.ToLower(m), query) {
					results = append(results, SearchResult{Value: m, Type: "Member", Artist: a.Name})
				}
			}
		}
		switch filter {
		case "location", "all":
			for _, loc := range data.AllLocations.Index {
				if loc.ID == a.ID {
					for _, l := range loc.Locations {
						if strings.Contains(strings.ToLower(l), query) {
							results = append(results, SearchResult{Value: l, Type: "Location", Artist: a.Name})
						}
					}
					break
				}
			}
		}

		switch filter {
		case "album", "all":
			if strings.Contains(strings.ToLower(a.FirstAlbum), query) {
				results = append(results, SearchResult{Value: a.FirstAlbum, Type: "First Album", Artist: a.Name})
			}
		}
		switch filter {
		case "date", "all":
			year := strconv.Itoa(a.CreationDate)
			if strings.Contains(year, query) {
				results = append(results, SearchResult{Value: year, Type: "Creation Date", Artist: a.Name})
			}
		}
	}

	json.NewEncoder(w).Encode(results)
}
