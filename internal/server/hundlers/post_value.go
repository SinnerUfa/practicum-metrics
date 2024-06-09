package hundlers

import (
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/err_codes"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	chi "github.com/go-chi/chi/v5"
)

func PostValue(log mlog.Logger, rep repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			name := chi.URLParam(r, "name")
			typ := chi.URLParam(r, "type")
			value := chi.URLParam(r, "value")
			metr := &metrics.Metric{
				Name:  name,
				Type:  typ,
				Value: value,
			}
			log.Info(r.URL.Path, "Name", name, "Type", typ, "Value", value)
			if name == "" {
				http.Error(w, codes.ErrPostValReqName.Error(), http.StatusBadRequest)
				log.Warning(codes.ErrPostValReqName)
				return
			}
			if typ == "" {
				http.Error(w, codes.ErrPostValReqType.Error(), http.StatusBadRequest)
				log.Warning(codes.ErrPostValReqType)
				return
			}
			if value == "" {
				http.Error(w, codes.ErrPostValReqValue.Error(), http.StatusBadRequest)
				log.Warning(codes.ErrGetValReqType)
				return
			}
			switch err := rep.Set(*metr); err {
			case codes.ErrRepParseInt, codes.ErrRepParseFloat, codes.ErrRepMetricNotSupported:
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Warning(err)
				return
			}

			w.Header().Set("Content-type", "text/plain ")
			w.WriteHeader(http.StatusOK)
		})
}
