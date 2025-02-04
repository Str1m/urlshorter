package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
	storage     map[string]string
	storagePath string
}

func main() {
	addr := *flag.String("addr", ":8080", "HTTP network address")
	fp := *flag.String("stor", "./storage.json", "JSON file path")
	flag.Parse()

	storage := make(map[string]string)

	infoLog := log.New(os.Stdout, "INFO:\t ", log.Ldate)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Lshortfile)

	app := &application{
		InfoLog:     infoLog,
		ErrorLog:    errorLog,
		storage:     storage,
		storagePath: fp,
	}

	server := http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}
	infoLog.Printf("Starting server on %s", addr)
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
