package hundlers

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type gzReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

func (g *gzReader) Read(p []byte) (n int, err error) {
	return g.zr.Read(p)
}

func Decompressor(log *slog.Logger, gz *gzip.Reader) func(http.Handler) http.Handler {
	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			if ct := r.Header.Get("Content-Type"); (ct == "application/json" || ct == "text/html") &&
				strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				gz.Reset(r.Body)
				r.Body = &gzReader{
					ReadCloser: r.Body,
					zr:         gz,
				}
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
