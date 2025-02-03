package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	str := []byte("Hello")
	w.Write(str)
}

type OriginURL struct {
	URL string `json:"url"`
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var url OriginURL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("URL revieved: %s", url)))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/geturl", GetURL)
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
