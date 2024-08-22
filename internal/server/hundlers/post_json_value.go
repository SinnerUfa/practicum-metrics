package hundlers

import (
	"bytes"
	"encoding/json"
	// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20
	slog "golang.org/x/exp/slog" // slog "log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
)

func PostJSONValue(getter metrics.Getter) http.HandlerFunc {
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
