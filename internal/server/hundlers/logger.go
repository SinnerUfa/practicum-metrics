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
	w.responseData.size = size
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

// Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
// Сведения об ответах должны содержать код статуса и размер содержимого ответа.

// func WithLogging(h http.Handler) http.Handler {
//     logFn := func(w http.ResponseWriter, r *http.Request) {
//         start := time.Now()

//         responseData := &responseData {
//             status: 0,
//             size: 0,
//         }
//         lw := loggingResponseWriter {
//             ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
//             responseData: responseData,
//         }
//         h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

//         duration := time.Since(start)

//         sugar.Infoln(
//             "uri", r.RequestURI,
//             "method", r.Method,
//             "status", responseData.status, // получаем перехваченный код статуса ответа
//             "duration", duration,
//             "size", responseData.size, // получаем перехваченный размер ответа
//         )
//     }
//     return http.HandlerFunc(logFn)
// }
