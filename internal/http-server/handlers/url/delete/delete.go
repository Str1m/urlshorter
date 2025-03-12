package delete

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"urlshorter/internal/lib/api/response"
	"urlshorter/internal/lib/logger/sl"
	"urlshorter/internal/storage"
)

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLDelete interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			json.NewEncoder(w).Encode(response.Error("invalid request"))
			return
		}

		err := urlDelete.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", slog.String("alias", alias))
				json.NewEncoder(w).Encode(response.Error("not found"))
				return
			}
			log.Error("failed to delete url", sl.Err(err))
			json.NewEncoder(w).Encode(response.Error("internal error"))

			return
		}

		responseOk(w, r, alias)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	json.NewEncoder(w).Encode(Response{
		Response: response.Ok(),
		Alias:    alias,
	})
}
