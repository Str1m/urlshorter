package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/geturl", app.GetURL)
	mux.HandleFunc("/short/", app.SendShortURL)

	return mux
}
