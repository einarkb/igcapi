package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marni/goigc"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT is not set")
	}

	track, err := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	if err != nil {
		fmt.Errorf("Problem reading the track", err)
	} else {
		fmt.Printf("Pilot: %s, gliderType: %s, date: %s",
			track.Pilot, track.GliderType, track.Date.String())
	}

	http.HandleFunc("/", hello)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
