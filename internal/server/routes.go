package server

import (
	"net/http"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	hundlers "github.com/SinnerUfa/practicum-metric/internal/server/hundlers"
	chi "github.com/go-chi/chi/v5"
)

func Routes(rep repository.Repository) http.Handler {
	r := chi.NewRouter()
	// r.Use(hundlers.Decompressor())
	r.Use(hundlers.Compressor())
	r.Use(hundlers.Logger())
	r.Route("/", func(r chi.Router) {
		r.Get("/", hundlers.GetList(rep))
		r.Get("/ping", hundlers.GetPing(rep))
		r.Post("/update/{type}/{name}/{value}", hundlers.PostUpdate(rep))
		r.Get("/value/{type}/{name}", hundlers.GetValue(rep))
		r.Post("/update/", hundlers.PostJSONUpdate(rep))
		r.Post("/value/", hundlers.PostJSONValue(rep))
	})

	return r
}
