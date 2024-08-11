package agent

import (
	"context"
	"log/slog"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	resty "github.com/go-resty/resty/v2"
)

type MetricPost struct {
	rep     repository.Repository
	ctx     context.Context
	adress  string
	counter uint
}

func NewPoster(ctx context.Context, rep repository.Repository, adress string) *MetricPost {
	return &MetricPost{ctx: ctx, rep: rep, adress: adress}
}

func (m *MetricPost) Post() {
	l := m.rep.List()
	client := resty.New()

	endpoint := "http://" + m.adress + "/update/"
	slog.Debug("", "l", l)
	for i, v := range l {
		req := client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetBody(v)
		p, err := req.Post(endpoint)
		slog.Debug("", "req.Body", req.Body)
		slog.Debug("", "content-encoding", p.Header().Get("Content-Encoding"), "response", p.String())
		if err != nil {
			slog.Warn("", "err", err, "i", i, "value", v)
		}
	}
	slog.Debug("Post metrics", "increment", m.counter)
}

func (m *MetricPost) Tick() {
	m.Post()
}
