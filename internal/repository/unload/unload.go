package unload

import (
	"context"
	"encoding/json"
	// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20
	slog "golang.org/x/exp/slog" // slog "log/slog"
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
	ticker *ticker.Ticker
}

func New(ctx context.Context, file string, interval uint) (*Unload, error) {
	wd, _ := os.Getwd()
	wd = filepath.Join(wd, filepath.Dir(file))
	file = filepath.Join(wd, filepath.Base(file))

	if err := os.MkdirAll(wd, os.ModePerm); err != nil {
		return nil, err
	}

	mem := memory.New()
	buf, err := os.ReadFile(file)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if len(buf) != 0 {
		out := make([]metrics.Metric, 0)
		if err := json.Unmarshal(buf, &out); err != nil {
			return nil, err
		}
		if err := mem.SetList(out); err != nil {
			return nil, err
		}
	}

	if interval == 0 {
		return &Unload{Memory: mem, always: true, file: file, ticker: nil}, nil
	}
	u := &Unload{Memory: mem, always: false, file: file}
	u.ticker = ticker.NewAndRun(ctx, interval, u)
	return u, nil
}

func (u *Unload) Set(m metrics.Metric) error {
	if u.always {
		if err := u.Memory.Set(m); err != nil {
			return err
		}
		m, _ := u.Memory.GetList()
		return unload(u.file, m)
	}
	return u.Memory.Set(m)
}

func (u *Unload) Tick() {
	slog.Debug("repository unload start")
	m, _ := u.Memory.GetList()
	if err := unload(u.file, m); err != nil {
		slog.Warn("repository unload error", "err", err)
		return
	}
	slog.Debug("repository unload end")
}

func (u *Unload) Close() error {
	u.ticker.Close()
	return nil
}

func unload(file string, out []metrics.Metric) error {
	data, err := json.MarshalIndent(out, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0666)
}
