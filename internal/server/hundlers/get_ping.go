package hundlers

import (
	"net/http"

	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
)

func GetPing(getter metrics.Getter) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			metr := &metrics.Metric{
				Name:  "ping",
				Type:  "",
				Value: metrics.String("ping"),
			}
			if err := getter.Get(metr); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
}
