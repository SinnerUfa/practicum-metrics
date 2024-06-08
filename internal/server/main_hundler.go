package server

import (
	"net/http"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func NewMainHundler(log mlog.Logger, cfg Config, rep repository.Repository) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		log,
		rep,
	)
	// // space for middleware, not necessary at this stage
	// hundler = any0(mux)
	// hundler = any0(hundler)
	// hundler = any0(hundler)
	// // ... etc
	// return hundler
	return mux
}
