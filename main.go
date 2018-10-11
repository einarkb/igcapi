package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT is not set")
	}
	http.HandleFunc("/", hello)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
