package hundlers

import (
	"log/slog"
	"net/http"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func GetPing(log *slog.Logger, rep repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
}
