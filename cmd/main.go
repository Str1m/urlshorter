package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	storage  *RedisStore
	ctx      context.Context
}

func main() {
	addr := *flag.String("addr", ":8080", "HTTP network address")
	redisAddr := *flag.String("redis", "localhost:6379", "Redis addres")
	flag.Parse()

	var ctx = context.Background()
	storage := NewRedisStore(redisAddr)

	infoLog := log.New(os.Stdout, "INFO:\t ", log.Ldate)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Lshortfile)

	app := &application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		storage:  storage,
		ctx:      ctx,
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
