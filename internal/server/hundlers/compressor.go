package hundlers

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
func Compressor(log *slog.Logger, gz *gzip.Writer) func(http.Handler) http.Handler {
	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			if ct := r.Header.Get("Content-Type"); ct != "application/json" && ct != "text/html" {
				h.ServeHTTP(w, r)
				return
			}
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				h.ServeHTTP(w, r)
				return
			}
			gz.Reset(w)
			w.Header().Set("Content-Encoding", "gzip")
			h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
