package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"urlshorter/internal/config"
	"urlshorter/internal/http-server/handlers/redirect"
	delete2 "urlshorter/internal/http-server/handlers/url/delete"
	"urlshorter/internal/http-server/handlers/url/save"
	"urlshorter/internal/http-server/middleware/logger"
	"urlshorter/internal/lib/logger/sl"
	"urlshorter/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.LoadConfig()

	log := setupLogger(cfg.Env)
	log.Info("start url-shorter", "env", slog.String("env", cfg.Env))

	storage, err := postgres.New(cfg.StorageConfig)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)

	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortner", map[string]string{
			cfg.UserConfig.Name: cfg.UserConfig.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete2.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))
	log.Info("starting server", slog.String("host", cfg.ServerConfig.Host), slog.String("port", cfg.ServerConfig.Port))

	addr := cfg.ServerConfig.Host + ":" + cfg.ServerConfig.Port
	srv := &http.Server{
		Addr:                         addr,
		Handler:                      router,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  cfg.ServerConfig.Timeout,
		WriteTimeout:                 cfg.ServerConfig.Timeout,
		IdleTimeout:                  cfg.ServerConfig.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
