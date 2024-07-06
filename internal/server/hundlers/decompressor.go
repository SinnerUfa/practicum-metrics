package hundlers

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
)

type gzReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (g gzReader) Read(p []byte) (n int, err error) {
	return g.zr.Read(p)
}

func (c *gzReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func Decompressor(log *slog.Logger, gz *gzip.Reader) func(http.Handler) http.Handler {
	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			log.Info("", "Decompressor", 0)
			// if ct := r.Header.Get("Content-Type"); (ct == "application/json" || ct == "text/html") &&
			// 	strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// 	gz.Reset(r.Body)
			// 	rr := &gzReader{
			// 		r:  r.Body,
			// 		zr: gz,
			// 	}
			// 	r.Body = rr
			// 	defer rr.Close()
			// }
			h.ServeHTTP(w, r)

			log.Info("", "Decompressor", 1)
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
