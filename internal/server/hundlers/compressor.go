package hundlers

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
	log    *slog.Logger
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
func Compressor(log *slog.Logger) func(http.Handler) http.Handler {
	var buf bytes.Buffer

	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	// _, err := gz.Write([]byte(" "))

	// if err != nil {
	// 	log.Warn("", "err", codes.ErrCompressor, "gz_err", err)

	// }

	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				h.ServeHTTP(w, r)
				return
			}
			gz.Reset(w)
			w.Header().Set("Content-Encoding", "gzip")
			h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz, log: log}, r)
			gz.Close()
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
