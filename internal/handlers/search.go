package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/data"
	"groupie-tracker/internal/utils"
)

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"replaceSpaces":  utils.ReplaceSpaces,
	"cleanDate":      utils.CleanDate,
	"formatLocation": utils.FormatLocation,
}).ParseGlob("templates/*.html"))

type SearchResult struct {
	Value  string `json:"Value"`
	Type   string `json:"Type"`
	Artist string `json:"Artist"`
	ID     int    `json:"ID"`
}

func addResult(results *[]SearchResult, seen map[string]bool, key, value, rtype, artist string, id int) {
	if !seen[key] {
		*results = append(*results, SearchResult{
			Value:  value,
			Type:   rtype,
			Artist: artist,
			ID:     id,
		})
		seen[key] = true
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("query"))
	query = utils.NormalizeQuery(query)
	if query == "" {
		json.NewEncoder(w).Encode([]SearchResult{})
		return
	}

	var results []SearchResult
	seen := make(map[string]bool)
	normalizedQuery := utils.NormalizeString(query)

	for _, artist := range data.AllArtists {
		// Check artist name match.
		if strings.Contains(utils.NormalizeString(artist.Name), normalizedQuery) {
			key := artist.Name + "_Artist"
			addResult(&results, seen, key, artist.Name+" — Artist/Band", "artist/band", artist.Name, artist.ID)
		}

		// Check members.
		for _, member := range artist.Members {
			if strings.Contains(utils.NormalizeString(member), normalizedQuery) {
				memberKey := member + "_Member"
				addResult(&results, seen, memberKey, member+" — Member", "member", artist.Name, artist.ID)
				// Also add the artist once.
				artistKey := artist.Name + "_Artist"
				addResult(&results, seen, artistKey, artist.Name+" — Artist/Band", "artist/band", artist.Name, artist.ID)
			}
		}

		// Check creation date.
		creationStr := strconv.Itoa(artist.CreationDate)
		if strings.HasPrefix(creationStr, query) {
			key := creationStr + "_Creation"
			addResult(&results, seen, key, creationStr+" — Creation Date", "creation date", artist.Name, artist.ID)
			artistKey := artist.Name + "_Artist_Creation"
			addResult(&results, seen, artistKey, artist.Name+" — Artist/Band", "artist/band", artist.Name, artist.ID)
		}

		// Check first album.
		cleanedAlbumDate := utils.CleanDate(artist.FirstAlbum)
		if strings.Contains(strings.ToLower(cleanedAlbumDate), query) {
			key := cleanedAlbumDate + "_First_Album"
			addResult(&results, seen, key, cleanedAlbumDate+" — First Album Date", "first album date", artist.Name, artist.ID)
			artistKey := artist.Name + "_Artist_FirstAlbum"
			addResult(&results, seen, artistKey, artist.Name+" — Artist/Band", "artist/band", artist.Name, artist.ID)
		}
	}

	// Check concert locations and dates.
	for _, rel := range data.AllRelations.Index {
		for location := range rel.DatesLocations {
			if len(query) >= 3 && strings.Contains(strings.ToLower(location), query) {
				locKey := location + "_Location"
				artistPtr := findArtistByID(rel.ID)
				if artistPtr != nil {
					addResult(&results, seen, locKey, utils.FormatLocation(location)+" — Location", "location", artistPtr.Name, artistPtr.ID)
					// Also add artist suggestion.
					artistKey := artistPtr.Name + "_Artist_Location"
					addResult(&results, seen, artistKey, artistPtr.Name+" — Artist/Band", "artist/band", artistPtr.Name, artistPtr.ID)

					// Check concert dates within this location.
					for _, date := range rel.DatesLocations[location] {
						cleanedDate := utils.CleanDate(date)
						dateKey := location + cleanedDate + "_ConcertDate"
						if strings.Contains(strings.ToLower(cleanedDate), query) {
							addResult(&results, seen, dateKey, cleanedDate+" — Concert Date", "concert date", artistPtr.Name, artistPtr.ID)
						}
					}
				}
			}
		}
	}

	json.NewEncoder(w).Encode(results)
}

// findArtistByID returns a pointer to the artist with the provided ID.
func findArtistByID(id int) *data.Artist {
	for _, artist := range data.AllArtists {
		if artist.ID == id {
			return &artist
		}
	}
	return nil
}

func ResultsPageHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("query"))
	query = utils.NormalizeQuery(query)
	normalizedQuery := utils.NormalizeString(query)

	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var matchedArtists []data.Artist
	seen := make(map[int]bool)

	for _, artist := range data.AllArtists {
		// Match artist name.
		if strings.Contains(utils.NormalizeString(artist.Name), normalizedQuery) {
			matchedArtists = append(matchedArtists, artist)
			seen[artist.ID] = true
			continue
		}
		// Match members.
		for _, member := range artist.Members {
			if strings.Contains(utils.NormalizeString(member), normalizedQuery) {
				if !seen[artist.ID] {
					matchedArtists = append(matchedArtists, artist)
					seen[artist.ID] = true
				}
				break
			}
		}
		// Match creation date.
		if strings.Contains(strconv.Itoa(artist.CreationDate), normalizedQuery) {
			if !seen[artist.ID] {
				matchedArtists = append(matchedArtists, artist)
				seen[artist.ID] = true
			}
		}
		// Match first album.
		if strings.Contains(utils.NormalizeString(utils.CleanDate(artist.FirstAlbum)), normalizedQuery) {
			if !seen[artist.ID] {
				matchedArtists = append(matchedArtists, artist)
				seen[artist.ID] = true
			}
		}
	}

	// Match locations.
	for _, rel := range data.AllRelations.Index {
		artistPtr := findArtistByID(rel.ID)
		if artistPtr == nil || seen[artistPtr.ID] {
			continue
		}
		for location := range rel.DatesLocations {
			if strings.Contains(utils.NormalizeString(location), normalizedQuery) {
				matchedArtists = append(matchedArtists, *artistPtr)
				seen[artistPtr.ID] = true
				break
			}
		}
	}

	err := templates.ExecuteTemplate(w, "results.html", matchedArtists)
	if err != nil {
		http.Error(w, "Failed to load results", http.StatusInternalServerError)
	}
}
