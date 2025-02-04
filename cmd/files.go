package main

import (
	"encoding/json"
	"log"
	"os"
)

func (app *application) loadFile() {
	f, err := os.Open(app.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&app.storage); err != nil {
		log.Fatal(err)
	}
}

func (app *application) saveFile() {
	f, err := os.Create(app.storagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&app.storage); err != nil {
		log.Fatal(err)
	}
}
