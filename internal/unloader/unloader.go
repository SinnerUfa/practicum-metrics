package unloader

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"log/slog"

	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	ticker "github.com/SinnerUfa/practicum-metric/internal/ticker"
)

func Load(file string, log *slog.Logger, rep repository.Repository) {
	var out []metrics.Metric
	b, err := os.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &out)
	if err != nil {
		log.Info("fail1", "err", err)
		return
	}
	if len(out) != 0 {
		log.Info("fail2")
		rep.SetList(out)
	}
	log.Info("loaded", "mem", rep.List())
}

type repSlave struct {
	repository.Repository
	file string
}

func (rp repSlave) Set(m metrics.Metric) error {
	err := rp.Repository.Set(m)
	mem := rp.Repository.List()
	if len(mem) != 0 {
		ship(rp.file, mem)
	}
	return err
}

func ship(file string, out []metrics.Metric) {
	data, err := json.MarshalIndent(out, "", "   ")
	if err != nil {
		return
	}
	os.MkdirAll(path.Dir(file), 0777)
	os.WriteFile(file, data, 0666)
}

type repTicker struct {
	rep  repository.Repository
	file string
	log  *slog.Logger
}

func (rp repTicker) Tick() error {
	mem := rp.rep.List()
	if len(mem) != 0 {
		ship(rp.file, mem)
	}
	rp.log.Info("Tick")
	return nil
}

func Save(ctx context.Context, file string, intrv uint, log *slog.Logger, rep repository.Repository) repository.Repository {
	if intrv == 0 {
		return repSlave{rep, file}
	}
	t := &repTicker{
		rep,
		file,
		log,
	}
	// ticker.NewAndRun(ctx, intrv>>1, t)
	ticker.NewAndRun(ctx, intrv, t)
	return rep
}
