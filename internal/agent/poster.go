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

func (m *MetricPost) Post() (err error) {
	l := m.rep.List()
	client := resty.New()

	endpoint := "http://" + m.adress + "/update/"
	m.log.Info("", "l", l)
	for i, v := range l {
		req := client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetBody(v)
		p, err := req.Post(endpoint)
		m.log.Info("", "req.Body", req.Body)
		m.log.Info("", "content-encoding", p.Header().Get("Content-Encoding"), "response", p.String())
		if err != nil {
			m.log.Warn("", "err", err, "i", i, "value", v)
		}
	}
	return
	m.log.Info("Post metrics", "increment", m.counter)
	return nil
}

func (m *MetricPost) Tick() error {
	return m.Post()
}
