package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
	fmt.Fprintf(w, "igchandler")
	switch r.Method {
	case "POST":
		fmt.Fprintf(w, "post")
		type InputURL struct {
			URL *string `json:"url"`
		}
		inputURL := InputURL{}
		err := json.NewDecoder(r.Body).Decode(inputURL)
		switch err {
		case io.EOF:
			fmt.Fprintf(w, "empty body")
		case nil:
			fmt.Fprintf(w, "other error")
		default:
			fmt.Fprintf(w, "has body")
		}
	default:
		fmt.Fprintf(w, "not post")
	}

}

func main() {
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
