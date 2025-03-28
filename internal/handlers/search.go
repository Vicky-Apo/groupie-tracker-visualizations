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
	ID     int    `json:"ID"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("query"))
	if query == "" {
		json.NewEncoder(w).Encode([]SearchResult{})
		return
	}

	var results []SearchResult
	seen := make(map[string]bool)

	for _, a := range data.AllArtists {
		// Artist/Band Name
		if strings.Contains(strings.ToLower(a.Name), query) {
			key := a.Name + "_artist"
			if !seen[key] {
				results = append(results, SearchResult{
					Value:  a.Name + " — artist/band",
					Type:   "artist/band",
					Artist: a.Name,
					ID:     a.ID,
				})
				seen[key] = true
			}
		}

		// Members
		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), query) {
				key := m + "_member"
				if !seen[key] {
					results = append(results, SearchResult{
						Value:  m + " — member",
						Type:   "member",
						Artist: a.Name,
						ID:     a.ID,
					})
					seen[key] = true
				}
			}
		}

		// First Album Date
		if strings.Contains(strings.ToLower(a.FirstAlbum), query) {
			key := a.FirstAlbum + "_firstalbum"
			if !seen[key] {
				results = append(results, SearchResult{
					Value:  a.FirstAlbum + " — first album date",
					Type:   "first album date",
					Artist: a.Name,
					ID:     a.ID,
				})
				seen[key] = true
			}
		}

		// Creation Date
		creationStr := strconv.Itoa(a.CreationDate)
		if strings.Contains(creationStr, query) {
			key := creationStr + "_creation"
			if !seen[key] {
				results = append(results, SearchResult{
					Value:  creationStr + " — creation date",
					Type:   "creation date",
					Artist: a.Name,
					ID:     a.ID,
				})
				seen[key] = true
			}
		}

		// Location (using RelationIndex)
		for _, rel := range data.AllRelations.Index {
			if rel.ID == a.ID {
				for location := range rel.DatesLocations {
					if strings.Contains(strings.ToLower(location), query) {
						key := location + "_location"
						if !seen[key] {
							results = append(results, SearchResult{
								Value:  location + " — location",
								Type:   "location",
								Artist: a.Name,
								ID:     a.ID,
							})
							seen[key] = true
						}
					}
				}
				break
			}
		}
	}

	json.NewEncoder(w).Encode(results)
}
