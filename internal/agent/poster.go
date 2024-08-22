package agent

import (
	// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20
	slog "golang.org/x/exp/slog" // slog "log/slog"
	"time"

	"context"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	resty "github.com/go-resty/resty/v2"
)

type MetricPost struct {
	rep     repository.Storage
	adress  string
	counter uint
	noBatch bool
}

func NewPoster(rep repository.Storage, adress string, noBatch bool) *MetricPost {
	return &MetricPost{rep: rep, adress: adress, noBatch: noBatch}
}

func (m *MetricPost) Post() (err error) {
	l, err := m.rep.GetList(context.Background())
	if err != nil {
		return
	}

	client := resty.New()
	if m.noBatch {
		endpoint := "http://" + m.adress + "/update/"
		for i, v := range l {
			req := client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetBody(v)
			p, err := req.Post(endpoint)
			slog.Debug("request", "body", req.Body, "encoding", p.Header().Get("Content-Encoding"), "response", p.String())
			if err != nil {
				slog.Warn("", "err", err, "i", i, "value", v)
			}
		}
	} else {
		endpoint := "http://" + m.adress + "/updates/"
		req := client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetBody(l)
		p, err := req.Post(endpoint)
		slog.Debug("request", "body", req.Body, "encoding", p.Header().Get("Content-Encoding"), "response", p.String())
		if err != nil {
			slog.Warn("", "err", err, "l", l)
		}
	}
	m.counter++
	slog.Debug("post metrics", "increment", m.counter)
	return err
}

func (m *MetricPost) Tick() {
	delays := []int{1, 3, 5}
	for _, delay := range delays {
		err := m.Post()
		if err != nil {
			time.Sleep(time.Duration(delay) * time.Second)
		} else {
			break
		}
	}
}
