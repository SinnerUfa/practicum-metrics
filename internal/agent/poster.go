package agent

import (
	"context"
	"log/slog"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	resty "github.com/go-resty/resty/v2"
)

type MetricPost struct {
	log     *slog.Logger
	rep     repository.Repository
	ctx     context.Context
	adress  string
	counter uint
}

func NewPoster(ctx context.Context, log *slog.Logger, rep repository.Repository, adress string) *MetricPost {
	return &MetricPost{ctx: ctx, log: log, rep: rep, adress: adress}
}

func (m *MetricPost) Post() error {
	l := m.rep.List()
	client := resty.New()

	endpoint := "http://" + m.adress + "/update/"

	for _, v := range l {
		_, err := client.R().SetHeader("Content-Type", "application/json").SetBody(v).Post(endpoint)
		if err != nil {
			m.log.Warn("", "err", err)
			return err
		}

	}
	m.log.Info("Post metrics", "increment", m.counter)
	return nil
}

func (m *MetricPost) Tick() error {
	return m.Post()
}
