package hundlers

import (
	"compress/gzip"
	"log/slog"
	"net/http"
	"strings"
)

// type logWriter struct {
// 	http.ResponseWriter
// 	status int
// 	buf    *bytes.Buffer
// }

// func (w *logWriter) Header() http.Header {
// 	return w.ResponseWriter.Header()
// }

// func (w *logWriter) Write(b []byte) (int, error) {
// 	mw := io.MultiWriter(w.ResponseWriter, w.buf)
// 	return mw.Write(b)
// }

// func (w *logWriter) WriteHeader(statusCode int) {
// 	w.ResponseWriter.WriteHeader(statusCode)
// 	w.status = statusCode
// }

func Decompressor(log *slog.Logger, gz *gzip.Reader) func(http.Handler) http.Handler {
	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			if ct := r.Header.Get("Content-Type"); (ct == "application/json" || ct == "text/html") &&
				strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				gz.Reset(r.Body)
				r.Body = gz
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
