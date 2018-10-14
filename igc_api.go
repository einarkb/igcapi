package igcapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	igc "github.com/marni/goigc"
)

// ReplyWithAllTrackIds replies with an array of all track ids
func ReplyWithAllTrackIds(w http.ResponseWriter) {
	ids := globalTracksDb.GetIDs()
	json.NewEncoder(w).Encode(ids)
}

// ReplyWithTrack replies with the track of specified id
func ReplyWithTrack(w http.ResponseWriter, id int) {
	trackURL, found := globalTracksDb.Get(id)
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
func ReplyWithSingleField(w http.ResponseWriter, id int, field string) {
	trackURL, found := globalTracksDb.Get(id)
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
