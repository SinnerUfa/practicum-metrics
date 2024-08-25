package hundlers

import (
	// "bytes"
	// "io"

	// codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	// hash "github.com/SinnerUfa/practicum-metric/internal/hash"
	// slog "log/slog"
	"net/http"
)

type hashWriter struct {
	http.ResponseWriter
	key string
}

func (w *hashWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *hashWriter) Write(b []byte) (int, error) {
	// w.ResponseWriter.Header().Set("HashSHA256", hash.Hash(b, w.key))
	return w.ResponseWriter.Write(b)
}

func (w *hashWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func Hasher(key string) func(http.Handler) http.Handler {
	mid := func(h http.Handler) http.Handler {
		hundler := func(w http.ResponseWriter, r *http.Request) {
			// if key == "" {
				h.ServeHTTP(w, r)
				// return
			// }
			// if r.Body == nil {
			// 	http.Error(w, codes.ErrHashNotBody.Error(), http.StatusBadRequest)
			// 	slog.Warn("", "err", codes.ErrHashNotBody)
			// 	return
			// }
			// headerHash := r.Header.Get("HashSHA256")
			// if headerHash == "" {
			// 	http.Error(w, codes.ErrHashNilHeader.Error(), http.StatusBadRequest)
			// 	slog.Warn("", "err", codes.ErrHashNilHeader)
			// 	return
			// }
			// b, _ := io.ReadAll(r.Body)
			// r.Body.Close()
			// r.Body = io.NopCloser(bytes.NewBuffer(b))
			// bodyHash := hash.Hash(b, key)
			// if headerHash != bodyHash {
			// 	http.Error(w, codes.ErrHashNotCorrect.Error(), http.StatusBadRequest)
			// 	slog.Warn("", "err", codes.ErrHashNotCorrect)
			// 	return
			// }
			// hw := &hashWriter{ResponseWriter: w, key: key}
			// h.ServeHTTP(hw, r)
		}
		return http.HandlerFunc(hundler)
	}
	return mid
}
