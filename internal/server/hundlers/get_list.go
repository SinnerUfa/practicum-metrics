package hundlers

import (
	"html/template"
	"log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func ValueString(m *metrics.Value) string {
	return m.String()
}

var tpi = `<!DOCTYPE html>
<html>
  <head>
      <meta content="width=device-width" charset="utf-8">
      <title>List</title>
  </head>
  <body>
    <div >
        <table >
            <thead>
                <tr><th>Name</th><th>Type</th><th>Value</th></tr>
            </thead>
            <tbody>
                {{ range . }}
                    <tr>
                        <td>{{ .Name }}</td>
                        <td>{{ .Type }}</td>
                        <td>{{ .Value | valueString}}</td>
                    </tr>
                {{ end }}
            </tbody>
        </table>
    </div>
  </body>
</html>`

func GetList(log *slog.Logger, rep repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t := template.New("list")
			var funcMap = template.FuncMap{
				"valueString": ValueString,
			}
			t.Funcs(funcMap)
			t, err := t.Parse(tpi)
			// t, err := t.ParseFiles("body_list.html", "index.html") // NOT WORK WITH FUNCS???
			if err != nil {
				http.Error(w, codes.ErrRepNotFound.Error(), http.StatusInternalServerError)
				log.Warn("", "err", codes.ErrGetLstParse)
				return
			}
			metrs := rep.List()
			w.Header().Set("Content-type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			t.Execute(w, metrs)
		})
}
