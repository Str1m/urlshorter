package save

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"urlshorter/internal/lib/api/response"
	"urlshorter/internal/lib/logger/sl"
	"urlshorter/internal/lib/random"
	"urlshorter/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-playground/validator/v10"
)

// TODO Move to cfg
const AliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			json.NewEncoder(w).Encode(response.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			var validationErr validator.ValidationErrors
			errors.As(err, &validationErr)

			log.Error("invalid request", sl.Err(err))

			json.NewEncoder(w).Encode(response.ValidationError(validationErr))
			return
		}

		alias := req.Alias

		// TODO Обработка ошибок если они повторяются
		if alias == "" {
			alias = random.NewRandomString(AliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))

				json.NewEncoder(w).Encode(response.Error("url already exists"))
				return
			}
			log.Error("failed to save url", sl.Err(err))

			json.NewEncoder(w).Encode(response.Error("failed to save url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOk(w, r, alias)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	json.NewEncoder(w).Encode(Response{
		Response: response.Ok(),
		Alias:    alias,
	})
}
