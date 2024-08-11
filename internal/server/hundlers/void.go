package hundlers

import (
	"net/http"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Void(rep repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "text/plain ")
			w.WriteHeader(http.StatusOK)
		})
}
