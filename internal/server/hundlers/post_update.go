package hundlers

import (
	"log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	chi "github.com/go-chi/chi/v5"
)

func PostUpdate(setter metrics.Setter) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			name := chi.URLParam(r, "name")
			typ := chi.URLParam(r, "type")
			value := chi.URLParam(r, "value")
			metr := &metrics.Metric{
				Name:  name,
				Type:  metrics.MetricType(typ),
				Value: metrics.String(value),
			}
			if name == "" {
				http.Error(w, codes.ErrPostValReqName.Error(), http.StatusBadRequest)
				slog.Warn("", "err", codes.ErrPostValReqName)
				return
			}
			if typ == "" {
				http.Error(w, codes.ErrPostValReqType.Error(), http.StatusBadRequest)
				slog.Warn("", "err", codes.ErrPostValReqType)
				return
			}
			if value == "" {
				http.Error(w, codes.ErrPostValReqValue.Error(), http.StatusBadRequest)
				slog.Warn("", "err", codes.ErrGetValReqType)
				return
			}
			switch err := setter.Set(r.Context(), *metr); err {
			case nil:
			case codes.ErrRepParseInt, codes.ErrRepParseFloat, codes.ErrRepMetricNotSupported:
				http.Error(w, err.Error(), http.StatusBadRequest)
				slog.Warn("", "err", err)
				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				slog.Warn("", "err", err)
				return
			}

			w.Header().Set("Content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
		})
}
