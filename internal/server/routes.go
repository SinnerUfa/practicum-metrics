package server

import (
	"log/slog"
	"net/http"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	hundlers "github.com/SinnerUfa/practicum-metric/internal/server/hundlers"
	chi "github.com/go-chi/chi/v5"
)

func Routes(log *slog.Logger, rep repository.Repository) http.Handler {

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", hundlers.GetList(log, rep))
		r.Post("/update/{type}/{name}/{value}", hundlers.PostValue(log, rep))
		r.Get("/value/{type}/{name}", hundlers.GetValue(log, rep))
	})
	return r
}
