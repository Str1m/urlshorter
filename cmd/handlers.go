package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) GetURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}
	var url OriginURL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	shortURL, err := app.SaveURL(url.URL)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		app.ErrorLog.Println(err)
		return
	}
	jsonURL := ShortURL{shortURL}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonURL); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
}

func (app *application) SendShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}
	code := getShortURL(r.URL.Path)
	if code == "" {
		http.Error(w, "Invalid short URL", http.StatusBadRequest)
		return
	}
	originURL, err := app.FindOriginalURL(code)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		app.ErrorLog.Println(err)
		return
	}
	//originURL, flag := app.storage[code]
	//if !flag {
	//	http.Error(w, "Short URL not found", http.StatusNotFound)
	//	return
	//}
	http.Redirect(w, r, originURL, http.StatusFound)
}
