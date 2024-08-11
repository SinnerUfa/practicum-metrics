package hundlers

import (
	"net/http"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func GetPing(rep repository.Storage) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
}
