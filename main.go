package igcapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// RootHandler Responds with 404
func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// Replies with the API's metadata
func handlerAPIMeta(w http.ResponseWriter, r *http.Request) {
	type MetaData struct {
		Uptime  string `json:"uptime"`
		Info    string `json:"info"`
		Version string `json:"version"`
	}

	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	encoder.Encode(MetaData{time.Since(globalStartTime).String(), "Service for IGC tracks.", "v1"})
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
				/*_, err2 := igc.ParseLocation(input.URL)
				if err2 != nil {
					http.Error(w, "could not get a track from url: "+input.URL, http.StatusNotFound)
					return
				}*/
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
			ReplyWithSingleField(w, id, parts[5])
			return
		}
		ReplyWithTrack(w, id)
	}
}

var globalTracksDb TrackURLsDB
var globalStartTime time.Time

func main() {
	globalTracksDb = TrackURLsDB{}
	globalTracksDb.Init()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT is not set")
	}

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/igcinfo/api/", handlerAPIMeta)
	http.HandleFunc("/igcinfo/api/igc/", IgcHandler)

	globalStartTime = time.Now()
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
