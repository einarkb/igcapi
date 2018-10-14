package main

import (
	"IGCApp/igcapi"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	igcapi.GlobalTracksDb = igcapi.TrackURLsDB{}
	igcapi.GlobalTracksDb.Init()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT is not set")
	}

	http.HandleFunc("/", igcapi.RootHandler)
	http.HandleFunc("/igcinfo/api/", igcapi.HandlerAPIMeta)
	http.HandleFunc("/igcinfo/api/igc/", igcapi.IgcHandler)

	igcapi.GlobalStartTime = time.Now()
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
