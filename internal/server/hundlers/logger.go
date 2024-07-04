package hundlers

import (
	"log/slog"
	"net/http"
	"time"
)

type responseData struct {
	status int
	size   int
}
type logWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (w *logWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *logWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.responseData.size += size
	return size, err
}

func (w *logWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.responseData.status = statusCode
}

func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			lw := &logWriter{
				ResponseWriter: w,
				responseData:   &responseData{},
			}
			start := time.Now()
			h.ServeHTTP(lw, r)
			log.Info("",
				slog.Group("request", slog.String("method", r.Method), slog.String("url", r.RequestURI)),
				slog.Group("response", slog.Int("status", lw.responseData.status), slog.Int("size", lw.responseData.size)),
				slog.Duration("duration", time.Since(start)),
			)
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
