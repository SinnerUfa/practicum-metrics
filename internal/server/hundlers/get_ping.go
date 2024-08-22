package hundlers

import (
	"context"
	"net/http"
)

type Pinger interface {
	Ping(ctx context.Context) bool
}

func GetPing(p Pinger) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !p.Ping(r.Context()) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
}
