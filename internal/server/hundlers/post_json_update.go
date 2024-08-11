package hundlers

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func PostJSONUpdate(rep repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				http.Error(w, codes.ErrPostNotJSON.Error(), http.StatusBadRequest)
				slog.Warn("", "err", codes.ErrPostNotJSON)
				return
			}

			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, codes.ErrPostBadBody.Error(), http.StatusBadRequest)
				slog.Warn("", "err", codes.ErrPostBadBody)
				return
			}
			metr := &metrics.Metric{}
			if err = json.Unmarshal(buf.Bytes(), metr); err != nil {
				http.Error(w, codes.ErrPostUnmarshal.Error(), http.StatusBadRequest)
				slog.Warn("", "err", codes.ErrPostUnmarshal)
				return
			}

			switch err := rep.Set(*metr); err {
			case codes.ErrRepParseInt, codes.ErrRepParseFloat, codes.ErrRepMetricNotSupported:
				http.Error(w, err.Error(), http.StatusBadRequest)
				slog.Warn("", "err", err)
				return
			}

			resp, err := json.Marshal(metr)
			if err != nil {
				http.Error(w, codes.ErrPostMarshal.Error(), http.StatusInternalServerError)
				slog.Warn("", "err", codes.ErrPostMarshal)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		})
}
