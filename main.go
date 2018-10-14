package main

import (
	"IGCApp/igcapi"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	igcapi.globalTracksDb = igcapi.TrackURLsDB{}
	igcapi.globalTracksDb.Init()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT is not set")
	}

	http.HandleFunc("/", igcapi.RootHandler)
	http.HandleFunc("/igcinfo/api/", igcapi.handlerAPIMeta)
	http.HandleFunc("/igcinfo/api/igc/", igcapi.IgcHandler)

	igcapi.globalStartTime = time.Now()
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
