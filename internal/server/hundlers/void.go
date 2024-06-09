package hundlers

import (
	// "fmt"
	// "html/template"

	"net/http"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Void(log mlog.Logger, rep repository.Repository) http.HandlerFunc {
	// t, err := template.ParseFiles("index.html", "body_list.html")
	// if err != nil {
	// 	log.Error("Loaded template incorrect")
	// }
	// log.Info("Loaded template")
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Info(r.URL.Path)
			// metrs := rep.List()
			w.Header().Set("Content-type", "text/plain ")
			w.WriteHeader(http.StatusOK)
			// log.Info(metrs)
			// body := ""
			// for _, v := range metrs {
			// 	body += fmt.Sprintln(v)
			// }
			// t.Execute(w, metrs)
		})
}
