package server

import (
	"net/http"

	hundlers "github.com/SinnerUfa/practicum-metric/internal/server/hundlers"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func addRoutes(mux *http.ServeMux, log mlog.Logger, rep repository.Repository) {
	mux.Handle("GET /{$}", hundlers.GetList(log, rep))
	mux.Handle("POST /update/", hundlers.PostValue(log, rep))
	mux.Handle("GET /value/", hundlers.GetValue(log, rep))
}
