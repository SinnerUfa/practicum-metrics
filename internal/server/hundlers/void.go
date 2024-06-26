package hundlers

import (
	"log/slog"
	"net/http"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Void(log *slog.Logger, rep repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Info(r.URL.Path)
			w.Header().Set("Content-type", "text/plain ")
			w.WriteHeader(http.StatusOK)
		})
}
