package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/marni/goigc"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
	track, err := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	if err != nil {
		fmt.Errorf("Problem reading the track", err)
	} else {
		fmt.Fprintf(w, "Pilot: %s, gliderType: %s, date: %s",
			track.Pilot, track.GliderType, track.Date.String())
	}
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
				id, added := globalTracksDb.Add(input.URL)
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
				ReplyWithAllTrackIds(w)
				return
			}
			http.Error(w, "Invalid url", http.StatusBadRequest)
			return
		}
		if len(parts) == 6 {
			ReplyWithSingleField(w, id, parts[6])
			return
		}
		ReplyWithTrack(w, id)
	}
}

// ReplyWithTrack replies with the track of specified id
func ReplyWithTrack(w http.ResponseWriter, id int) {
	trackURL, found := globalTracksDb.Get(id)
	if !found {
		http.Error(w, "No track found with id: "+strconv.Itoa(id), http.StatusBadRequest)
		return
	}
	track, err := igc.ParseLocation(trackURL)
	if err != nil {
		fmt.Errorf("Problem reading the track", err)
	} else {
		fmt.Fprintf(w, "Pilot: %s, gliderType: %s, date: %s",
			track.Pilot, track.GliderType, track.Date.String())
	}
}

// ReplyWithAllTrackIds replies with an array of all track ids
func ReplyWithAllTrackIds(w http.ResponseWriter) {
	ids := globalTracksDb.GetIDs()
	json.NewEncoder(w).Encode(ids)
}

// ReplyWithSingleField replies with the information found in the specified field of trakc with given id
func ReplyWithSingleField(w http.ResponseWriter, id int, field string) {
	trackURL, found := globalTracksDb.Get(id)
	if !found {
		http.Error(w, "No track found with id: "+strconv.Itoa(id), http.StatusBadRequest)
		return
	}
	track, err := igc.ParseLocation(trackURL)
	if err != nil {
		fmt.Errorf("Problem reading the track", err)
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
			// calc
		default:
			http.Error(w, "invalid field specified", http.StatusBadRequest)
		}
	}
}

var globalTracksDb TrackURLsDB

func main() {
	globalTracksDb = TrackURLsDB{}
	globalTracksDb.Init()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT is not set")
	}

	http.HandleFunc("/", hello)
	http.HandleFunc("/igcinfo/api/igc/", IgcHandler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
