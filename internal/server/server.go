package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Run(ctx context.Context, log *slog.Logger, cfg Config) error {
	rep := repository.New()

	var buf bytes.Buffer
	gzw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	_, err := gzw.Write([]byte(" "))
	if err != nil {
		log.Warn("", "err", codes.ErrCompressor, "gzerr", err, "gzw", gzw)
	}

	gzr, err := gzip.NewReader(&buf)
	if err != nil {
		log.Warn("", "err", codes.ErrDecompressor, "gzerr", err, "gzr", gzr)
	}

	httpServer := &http.Server{
		Addr:    cfg.Adress,
		Handler: Routes(log, rep, gzr, gzw),
	}
	errChan := make(chan error)
	go func(log *slog.Logger, ch chan error) {
		log.Info("start server on adress", "Addr:", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn("", "err", codes.ErrSrvListen, "server err", err)
			ch <- codes.ErrSrvListen
		}
	}(log, errChan)

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Warn("", "err", codes.ErrSrvShutdown, "server err:", err)
			return codes.ErrSrvShutdown
		}
	}
	return nil
}
