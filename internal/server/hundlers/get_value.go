package hundlers

import (
	"net/http"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func GetValue(log mlog.Logger, rep repository.Repository) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			metr, err := SplitURL(r.URL.Path)
			log.Info(r.URL.Path)
			if err != nil {
				if err == ExBadReqStringType {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Warning(err)
					return
				}
				if err == ExBadReqStringName {
					http.Error(w, err.Error(), http.StatusNotFound)
					log.Warning(err)
					return
				}
			}
			if err := rep.Get(metr); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Warning(err)
				return
			}
			w.Header().Set("Content-type", "text/plain ")
			w.WriteHeader(http.StatusOK)
			log.Info(*metr)
			w.Write([]byte(metr.Value))
		})
}
