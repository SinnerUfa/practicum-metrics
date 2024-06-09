package server

import (
	"net/http"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	hundlers "github.com/SinnerUfa/practicum-metric/internal/server/hundlers"
	chi "github.com/go-chi/chi/v5"
)

func Routes(log mlog.Logger, cfg Config, rep repository.Repository) http.Handler {

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", hundlers.Void(log, rep))
		r.Post("/update/{type}/{name}/{value}", hundlers.PostValue(log, rep))
		r.Get("/value/{type}/{name}", hundlers.GetValue(log, rep))
	})
	return r
}
