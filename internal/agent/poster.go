package agent

import (
	"log/slog"
	"time"

	"context"
	// "encoding/json"
	// hash "github.com/SinnerUfa/practicum-metric/internal/hash"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	resty "github.com/go-resty/resty/v2"
)

type MetricPost struct {
	rep     repository.Storage
	adress  string
	counter uint
	noBatch bool
	key     string
}

func NewPoster(rep repository.Storage, adress string, noBatch bool, key string) *MetricPost {
	return &MetricPost{rep: rep, adress: adress, noBatch: noBatch, key: key}
}

// т.к. resty парсит и сжимает прямо перед запросом сложно достать обработанное тело запроса
func (m *MetricPost) Post() (err error) {
	l, err := m.rep.GetList(context.Background())
	if err != nil {
		return
	}
	client := resty.New()
	if m.noBatch {
		endpoint := "http://" + m.adress + "/update/"
		for i, v := range l {
			req := client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip")
			// if m.key != "" {
			// 	if b, err := json.Marshal(v); err != nil {
			// 		req.SetHeader("HashSHA256", hash.Hash(b, m.key))
			// 	}
			// }
			p, err := req.SetBody(v).Post(endpoint)
			slog.Debug("request", "body", req.Body, "encoding", p.Header().Get("Content-Encoding"), "response", p.String())
			if err != nil {
				slog.Warn("", "err", err, "i", i, "value", v)
			}
		}
	} else {
		endpoint := "http://" + m.adress + "/updates/"
		req := client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip")
		// if m.key != "" {
		// 	if b, err := json.Marshal(l); err != nil {
		// 		req.SetHeader("HashSHA256", hash.Hash(b, m.key))
		// 	}
		// }
		p, err := req.SetBody(l).Post(endpoint)
		slog.Debug("request", "body", req.Body, "encoding", p.Header().Get("Content-Encoding"), "response", p.String())
		if err != nil {
			slog.Warn("", "err", err, "l", l)
		}
	}
	m.counter++
	slog.Debug("post metrics", "increment", m.counter)
	return err
}

func (m *MetricPost) Tick() { // вообще у resty есть повторы, нужно разобраться
	for _, delay := range []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second} {
		if m.Post() == nil {
			break
		}
		time.Sleep(delay)
	}
}
