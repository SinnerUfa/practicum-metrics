package hundlers

import (
	// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20
	// "cmp"
	"html/template"
	// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20
	slog "golang.org/x/exp/slog" // slog "log/slog"
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
			metrs, err := getter.GetList(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				slog.Warn("", "err", err)
				return
			}
			slices.SortFunc(metrs, func(a, b metrics.Metric) int {
				if strings.ToLower(a.Name) == strings.ToLower(b.Name) {
					return 1
				}
				return 0
				// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20
				// return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
			})
			w.Header().Set("Content-type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			t.Execute(w, metrs)
		})
}
