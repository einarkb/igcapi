package igcapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	igc "github.com/marni/goigc"
)

// GlobalStartTime contains the start time of the app
var GlobalStartTime time.Time

// RootHandler Responds with 404
func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// HandlerAPIMeta Replies with the API's metadata
func HandlerAPIMeta(w http.ResponseWriter, r *http.Request) {
	type MetaData struct {
		Uptime  string `json:"uptime"`
		Info    string `json:"info"`
		Version string `json:"version"`
	}

	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	encoder.Encode(MetaData{time.Since(GlobalStartTime).String(), "Service for IGC tracks.", "v1"})
}

// IgcHandler handles /igcinfo/api/igc/
func IgcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	if r.Method == "POST" {
		type Input struct {
			URL string `json:"url"`
		}
		type ID struct {
			ID int `json:"id"`
		}
		var input Input
		err := json.NewDecoder(r.Body).Decode(&input)
		if err == nil {
			if input.URL != "" {
				_, err2 := igc.ParseLocation(input.URL)
				if err2 != nil {
					http.Error(w, "could not get a track from url: "+input.URL, http.StatusNotFound)
					return
				}
				id, added := GlobalTracksDb.Add(input.URL)
				if added {
					json.NewEncoder(w).Encode(ID{id})
				} else {
					http.Error(w, "track already exists with id: "+strconv.Itoa(id), http.StatusBadRequest)
				}
			} else {
				http.Error(w, "Body does not contain an url", http.StatusBadRequest)
			}
		} else if err == io.EOF {
			http.Error(w, "POST body is empty", http.StatusBadRequest)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	} else if r.Method == "GET" {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) > 6 {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(parts[4])
		if err != nil {
			if parts[4] == "" {
				ReplyWithAllTrackIds(w, &GlobalTracksDb)
				return
			}
			http.Error(w, "Invalid url", http.StatusBadRequest)
			return
		}
		if len(parts) == 6 {
			ReplyWithSingleField(w, id, parts[5], &GlobalTracksDb)
			return
		}
		ReplyWithTrack(w, id, &GlobalTracksDb)
	}
}

// ReplyWithAllTrackIds replies with an array of all track ids
func ReplyWithAllTrackIds(w http.ResponseWriter, db *TrackURLsDB) {
	ids := db.GetIDs()
	json.NewEncoder(w).Encode(ids)
}

// ReplyWithTrack replies with the track of specified id
func ReplyWithTrack(w http.ResponseWriter, id int, db *TrackURLsDB) {
	trackURL, found := db.Get(id)
	if !found {
		http.Error(w, "No track found with id: "+strconv.Itoa(id), http.StatusNotFound)
		return
	}
	track, err := igc.ParseLocation(trackURL)
	if err != nil {
		http.Error(w, "Problem reading the track", http.StatusServiceUnavailable)
	} else {
		fmt.Fprintf(w, "Pilot: %s, gliderType: %s, date: %s",
			track.Pilot, track.GliderType, track.Date.String())
	}
}

// ReplyWithSingleField replies with the information found in the specified field of trakc with given id
func ReplyWithSingleField(w http.ResponseWriter, id int, field string, db *TrackURLsDB) {
	trackURL, found := db.Get(id)
	if !found {
		http.Error(w, "No track found with id: "+strconv.Itoa(id), http.StatusNotFound)
		return
	}
	track, err := igc.ParseLocation(trackURL)
	if err != nil {
		http.Error(w, "Problem reading the track", http.StatusServiceUnavailable)
	} else {
		switch field {
		case "pilot":
			fmt.Fprintf(w, "Pilot: %s", track.Pilot)
		case "glider":
			fmt.Fprintf(w, "glider: %s", track.GliderType)
		case "glider_id":
			fmt.Fprintf(w, "glider_id: %s", track.GliderID)
		case "H_date":
			fmt.Fprintf(w, "H_date: %s", track.Date.String())
		case "track_length":
			d := 0.0
			for i := 0; i < len(track.Points)-1; i++ {
				d += track.Points[i].Distance(track.Points[i+1])
			}
			fmt.Fprintf(w, "distance: %f", d)
		default:
			http.Error(w, "invalid field specified", http.StatusBadRequest)
		}
	}
}
