package hundlers

import (
	"bytes"
	"compress/gzip"
	"io"
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
func Compressor() func(http.Handler) http.Handler {
	var buf bytes.Buffer

	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)

	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				h.ServeHTTP(w, r)
				return
			}
			gz.Reset(w)
			w.Header().Set("Content-Encoding", "gzip")
			h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
			gz.Close()
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
