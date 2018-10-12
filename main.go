package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
		json.NewDecoder(r.Body).Decode(&input)
		if input.URL != "" {
			id, added := globalTracksDb.Add(input.URL)
			if added {
				json.NewEncoder(w).Encode(ID{id})
			} else {
				http.Error(w, "track already exists with id: "+strconv.Itoa(id), http.StatusBadRequest)
			}
		}

		//}
		//fmt.Fprintf(w, "The URL already exists")
	}
	/*switch r.Method {
	case "POST":
		type Input struct {
			URL string `json:"url"`
		}
		var input Input
		err := json.NewDecoder(r.Body).Decode(&input)
		switch err {
		case nil:
			if input.URL == "" {
				http.Error(w, "No url in body", http.StatusBadRequest)
			} else {
				//check if valid igc url
				if id := globalTracksDb.Add(input.URL); id > 0 {
					fmt.Fprintf(w, "{\"id\": "+strconv.Itoa(id)+"}")
					return
				}
				fmt.Fprintf(w, "The URL already exists")
			}
			fmt.Fprintf(w, "heeei")
		case io.EOF:
			http.Error(w, "POST body is empty", http.StatusBadRequest)
		default:
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
		}

	case "GET":
		type IDStruct struct {
			ID int `json:"id"`
		}
		for _, v := range globalTracksDb.ids {
			idStruct := IDStruct{v}
			json.NewEncoder(w).Encode(idStruct)
		}
	default:
		fmt.Fprintf(w, "not post")
	}*/

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
