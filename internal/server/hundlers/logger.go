package hundlers

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type logWriter struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (w *logWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *logWriter) Write(b []byte) (int, error) {
	mw := io.MultiWriter(w.ResponseWriter, &(w.buf))
	return mw.Write(b)
}

func (w *logWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			lw := &logWriter{
				ResponseWriter: w,
			}
			start := time.Now()
			body := ""
			if r.Body != nil {
				b, _ := io.ReadAll(r.Body)
				r.Body.Close()
				buf := bytes.NewBuffer(b)
				body = buf.String()
				r.Body = io.NopCloser(buf)
			}

			h.ServeHTTP(lw, r)
			log.Info("",
				slog.Group("request",
					slog.String("method", r.Method),
					slog.String("url", r.RequestURI),
					slog.String("content-type", strings.Join(r.Header.Values("Content-Type"), ";")),
					slog.String("content-encoding", strings.Join(r.Header.Values("Content-Encoding"), ";")),
					slog.String("accept-encoding", strings.Join(r.Header.Values("Accept-Encoding"), ";")),
					slog.String("body", body)),
				slog.Group("response",
					slog.Int("status", lw.status),
					slog.Int("size", lw.buf.Len()),
					slog.String("body", lw.buf.String())),
				slog.Duration("duration", time.Since(start)))
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
