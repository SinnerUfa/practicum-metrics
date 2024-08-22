package hundlers

import (
	"log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	chi "github.com/go-chi/chi/v5"
)

func GetValue(getter metrics.Getter) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			name := chi.URLParam(r, "name")
			typ := chi.URLParam(r, "type")
			metr := &metrics.Metric{
				Name: name,
				Type: metrics.MetricType(typ),
			}
			if name == "" {
				http.Error(w, codes.ErrGetValReqName.Error(), http.StatusBadRequest)
				slog.Warn("", "err", codes.ErrGetValReqName)
				return
			}
			if typ == "" {
				http.Error(w, codes.ErrGetValReqType.Error(), http.StatusBadRequest)
				slog.Warn("", "err", codes.ErrGetValReqType)
				return
			}

			switch err := getter.Get(r.Context(), metr); err {
			case nil:
			case codes.ErrRepNotFound:
				http.Error(w, codes.ErrRepNotFound.Error(), http.StatusNotFound)
				slog.Warn("", "err", codes.ErrRepNotFound)
				return
			case codes.ErrRepMetricNotSupported:
				http.Error(w, codes.ErrRepMetricNotSupported.Error(), http.StatusBadRequest)
				slog.Warn("", "err", codes.ErrRepMetricNotSupported)
				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				slog.Warn("", "err", err)
				return
			}

			w.Header().Set("Content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(metr.Value.String()))
		})
}
