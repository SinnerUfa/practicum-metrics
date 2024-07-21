package unload

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/SinnerUfa/practicum-metric/internal/ticker"

	"github.com/SinnerUfa/practicum-metric/internal/repository/memory"

	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
)

type Unload struct {
	*memory.Memory
	always bool
	file   string
	log    *slog.Logger
}

func New(ctx context.Context, file string, interval uint, log *slog.Logger) (*Unload, error) {
	wd, _ := os.Getwd()
	wd = filepath.Join(wd, filepath.Dir(file))
	ex, _ := os.Executable()
	file = filepath.Join(wd, filepath.Base(file))
	log.Info("path info", "executable", ex, "request file", file)

	if err := os.MkdirAll(wd, os.ModePerm); err != nil {
		return nil, err
	}

	mem := memory.New()
	buf, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("file storage not found")
		} else {
			return nil, err
		}
	}
	if len(buf) == 0 {
		log.Info("len == 0")
	} else {
		out := make([]metrics.Metric, 0)
		if err := json.Unmarshal(buf, &out); err != nil {
			return nil, err
		}
		if err := mem.SetList(out); err != nil {
			return nil, err
		}
	}

	if interval == 0 {
		return &Unload{Memory: mem, always: true, file: file, log: log}, nil
	}
	u := &Unload{Memory: mem, always: true, file: file, log: log}
	ticker.NewAndRun(ctx, interval, u)
	return u, nil
}

func (u *Unload) Set(m metrics.Metric) error {
	if u.always {
		if err := u.Memory.Set(m); err != nil {
			return err
		}
		return ship(u.file, u.Memory.List())
	}
	return u.Memory.Set(m)
}

func (u *Unload) Tick() {
	if err := ship(u.file, u.Memory.List()); err != nil {
		u.log.Info("Tick error", "err", err)
		return
	}
	u.log.Info("Tick good")
}

func ship(file string, out []metrics.Metric) error {
	data, err := json.MarshalIndent(out, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0666)
}
