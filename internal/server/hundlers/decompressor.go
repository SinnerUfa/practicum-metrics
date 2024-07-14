// файл реализации декомпрессора gzip, не используется, потому что работает без него (за счет встроенной функции роутора)
package hundlers

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
)

type gzReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (с *gzReader) Read(p []byte) (n int, err error) {
	return с.zr.Read(p)
}

func (c *gzReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func Decompressor(log *slog.Logger) func(http.Handler) http.Handler {
	var buf bytes.Buffer

	gzw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	gzw.Write([]byte(" "))
	gz, err := gzip.NewReader(&buf)

	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			if ct := r.Header.Get("Content-Type"); (ct == "application/json" || ct == "text/html") &&
				strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				if err != nil {
					http.Error(w, codes.ErrDecompressor.Error(), http.StatusInternalServerError)
					log.Warn("", "err", codes.ErrDecompressor, "gz_err", err)
				}
				gz.Reset(r.Body)
				rr := &gzReader{
					r:  r.Body,
					zr: gz,
				}
				r.Body = rr
				defer rr.Close()
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
