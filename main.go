package main

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

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
	shortURL := createShortURL(url)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(shortURL); err != nil {
		log.Fatal(err)
	}
}

type ShortURL struct {
	URL string `json:"short_url"`
}

type OriginURL struct {
	URL string `json:"url"`
}

func createShortURL(URL OriginURL) ShortURL {
	url := URL.URL + strconv.Itoa(rand.Int())
	h := crypto.SHA3_256.New()
	h.Write([]byte(url))
	hash := h.Sum(nil)
	encoded := base64.RawURLEncoding.EncodeToString(hash)[:6]
	storage[URL.URL] = encoded
	saveFile()
	return ShortURL{URL: encoded}
}

const fp = "storage.json"

var storage map[string]string

func loadFile() {
	f, err := os.Open(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}
	defer f.Close()

	if storage == nil {
		storage = make(map[string]string)
	}
	if err := json.NewDecoder(f).Decode(&storage); err != nil {
		log.Fatal(err)
	}
}

func saveFile() {
	f, err := os.Create(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&storage); err != nil {
		log.Fatal(err)
	}
}
func main() {
	loadFile()
	fmt.Println(storage)
	mux := http.NewServeMux()
	mux.HandleFunc("/geturl", GetURL)
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)

}
