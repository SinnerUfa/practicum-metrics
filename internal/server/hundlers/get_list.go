package hundlers

import (
	"html/template"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func GetList(log mlog.Logger, rep repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t, err := template.ParseFiles("index.html", "body_list.html")
			if err != nil {
				http.Error(w, codes.ErrRepNotFound.Error(), http.StatusInternalServerError)
				log.Warning(codes.ErrGetLstParse)
			}
			log.Info(r.URL.Path)
			metrs := rep.List()
			w.Header().Set("Content-type", "text/plain ")
			w.WriteHeader(http.StatusOK)
			log.Info(metrs)
			t.Execute(w, metrs)
		})
}
