package hundlers

import (
	"cmp"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
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

func GetList(getter metrics.ListGetter) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t := template.New("list")
			var funcMap = template.FuncMap{
				"valueString": ValueString,
			}
			t.Funcs(funcMap)
			t, err := t.Parse(tpi)
			if err != nil {
				http.Error(w, codes.ErrRepNotFound.Error(), http.StatusInternalServerError)
				slog.Warn("", "err", codes.ErrGetLstParse)
				return
			}
			metrs, err := getter.GetList()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				slog.Warn("", "err", err)
				return
			}
			slices.SortFunc(metrs, func(a, b metrics.Metric) int {
				return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
			})
			w.Header().Set("Content-type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			t.Execute(w, metrs)
		})
}
