package server

import (
	"compress/gzip"
	"log/slog"
	"net/http"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	hundlers "github.com/SinnerUfa/practicum-metric/internal/server/hundlers"
	chi "github.com/go-chi/chi/v5"
)

func Routes(log *slog.Logger, rep repository.Repository, gzr *gzip.Reader, gzw *gzip.Writer) http.Handler {
	r := chi.NewRouter()
	if gzr != nil {
		r.Use(hundlers.Decompressor(log, gzr))
	}
	if gzw != nil {
		r.Use(hundlers.Compressor(log, gzw))
	}
	r.Use(hundlers.Logger(log))
	r.Route("/", func(r chi.Router) {
		r.Get("/", hundlers.GetList(log, rep))
		r.Post("/update/{type}/{name}/{value}", hundlers.PostUpdate(log, rep))
		r.Get("/value/{type}/{name}", hundlers.GetValue(log, rep))
		r.Post("/update/", hundlers.PostJSONUpdate(log, rep))
		r.Post("/value/", hundlers.PostJSONValue(log, rep))
	})

	return r
}
