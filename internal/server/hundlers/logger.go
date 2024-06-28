package hundlers

import (
	"log/slog"
	"net/http"
)

func Logger(log *slog.Logger) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Info("", "request", r, "response", w)
		})
}
