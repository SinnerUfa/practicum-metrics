package hundlers

import (
	"log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	chi "github.com/go-chi/chi/v5"
)

func GetValue(log *slog.Logger, rep repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			name := chi.URLParam(r, "name")
			typ := chi.URLParam(r, "type")
			metr := &metrics.Metric{
				Name: name,
				Type: typ,
			}
			if name == "" {
				http.Error(w, codes.ErrGetValReqName.Error(), http.StatusBadRequest)
				log.Warn("", "err", codes.ErrGetValReqName)
				return
			}
			if typ == "" {
				http.Error(w, codes.ErrGetValReqType.Error(), http.StatusBadRequest)
				log.Warn("", "err", codes.ErrGetValReqType)
				return
			}

			switch rep.Get(metr) {
			case codes.ErrRepNotFound:
				http.Error(w, codes.ErrRepNotFound.Error(), http.StatusNotFound)
				log.Warn("", "err", codes.ErrRepNotFound)
				return
			case codes.ErrRepMetricNotSupported:
				http.Error(w, codes.ErrRepMetricNotSupported.Error(), http.StatusBadRequest)
				log.Warn("", "err", codes.ErrRepMetricNotSupported)
				return
			}

			w.Header().Set("Content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(metr.Value))
		})
}
