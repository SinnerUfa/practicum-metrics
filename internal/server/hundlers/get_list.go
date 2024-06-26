package hundlers

import (
	"html/template"
	"log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func GetList(log *slog.Logger, rep repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t, err := template.ParseFiles("index.html", "body_list.html")
			if err != nil {
				http.Error(w, codes.ErrRepNotFound.Error(), http.StatusInternalServerError)
				log.Warn("", "err", codes.ErrGetLstParse)
				return
			}
			log.Info("get list request", "URL", r.URL.Path)
			metrs := rep.List()
			w.Header().Set("Content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
			log.Info("", "metrics", metrs)
			t.Execute(w, metrs)
		})
}
