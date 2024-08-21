package agent

import (
	"log/slog"
	"time"

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
	l, err := m.rep.GetList()
	if err != nil {
		return
	}
	client := resty.New()
	if m.noBatch {
		endpoint := "http://" + m.adress + "/update/"
		for i, v := range l {
			req := client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetBody(v)
			p, err := req.Post(endpoint)
			slog.Debug("", "req.Body", req.Body)
			slog.Debug("", "content-encoding", p.Header().Get("Content-Encoding"), "response", p.String())
			if err != nil {
				slog.Warn("", "err", err, "i", i, "value", v)
			}
		}
	} else {
		endpoint := "http://" + m.adress + "/updates/"
		req := client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetBody(l)
		p, err := req.Post(endpoint)
		slog.Debug("", "req.Body", req.Body)
		slog.Debug("", "content-encoding", p.Header().Get("Content-Encoding"), "response", p.String())
		if err != nil {
			slog.Warn("", "err", err, "l", l)
		}
	}
	m.counter++
	slog.Debug("Post metrics", "increment", m.counter)
	return err
}

func (m *MetricPost) Tick() {
	var delay int

	for counter := 0; counter < 3; counter++ {
		err := m.Post()
		switch counter {
		case 0:
			delay = 1
		case 1:
			delay = 3
		case 2:
			delay = 5
		}
		if err != nil {
			time.Sleep(time.Duration(delay) * time.Second)
		} else {
			break
		}
	}
}
