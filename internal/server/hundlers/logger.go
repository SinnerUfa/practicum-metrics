package hundlers

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type logWriter struct {
	http.ResponseWriter
	status int
	buf    *bytes.Buffer
}

func (w *logWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *logWriter) Write(b []byte) (int, error) {
	mw := io.MultiWriter(w.ResponseWriter, w.buf)
	return mw.Write(b)
}

func (w *logWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.status = statusCode
}

func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			lw := &logWriter{
				ResponseWriter: w,
				buf:            bytes.NewBuffer(make([]byte, 512)),
			}
			start := time.Now()
			h.ServeHTTP(lw, r)

			var buf bytes.Buffer

			if r.Body != nil {
				buf.ReadFrom(r.Body)
			}
			log.Info("",
				slog.Group("request",
					slog.String("method", r.Method),
					slog.String("url", r.RequestURI),
					slog.String("body", buf.String())),
				slog.Group("response",
					slog.Int("status", lw.status),
					slog.Int("size", lw.buf.Len()),
					/*slog.String("body", lw.buf.String())*/),
				slog.Duration("duration", time.Since(start)))
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
