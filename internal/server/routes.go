package server

import (
	"net/http"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	hundlers "github.com/SinnerUfa/practicum-metric/internal/server/hundlers"
	chi "github.com/go-chi/chi/v5"
)

func Routes(rep *repository.Repository, key string) http.Handler {
	r := chi.NewRouter()
	// r.Use(hundlers.Decompressor())
	r.Use(hundlers.Compressor())
	r.Use(hundlers.Hasher(key))
	r.Use(hundlers.Logger())

	r.Route("/", func(r chi.Router) {
		s := rep.Storage()
		r.Get("/", hundlers.GetList(s))
		r.Get("/ping", hundlers.GetPing(rep))
		r.Post("/update/{type}/{name}/{value}", hundlers.PostUpdate(s))
		r.Get("/value/{type}/{name}", hundlers.GetValue(s))
		r.Post("/update/", hundlers.PostJSONUpdate(s))
		r.Post("/value/", hundlers.PostJSONValue(s))
		r.Post("/updates/", hundlers.PostUpdates(s))
	})

	return r
}
