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
	limit   uint
}

func NewPoster(rep repository.Storage, adress string, noBatch bool, key string, limit uint) *MetricPost {
	return &MetricPost{rep: rep, adress: adress, noBatch: noBatch, key: key, limit: limit}
}

func (m *MetricPost) Post() (err error) {
	l, err := m.rep.GetList(context.Background())
	if err != nil {
		return
	}
	if len(l) == 0 {
		return
	}
	client := resty.New()
	if m.noBatch {
		endpoint := "http://" + m.adress + "/update/"
		var limit uint = uint(len(l))
		if m.limit != 0 {
			limit = m.limit
		}
		jobs := make(chan any, limit)
		results := make(chan any, limit)
		for w := 0; w < int(limit); w++ {
			go func(jobs <-chan any, results chan<- any) {
				for j := range jobs {
					req := client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip")
					// if m.key != "" {
					// 	if b, err := json.Marshal(v); err != nil {
					// 		req.SetHeader("HashSHA256", hash.Hash(b, m.key))
					// 	}
					// }
					p, err := req.SetBody(j).Post(endpoint)
					slog.Debug("request", "body", req.Body, "encoding", p.Header().Get("Content-Encoding"), "response", p.String())
					if err != nil {
						slog.Warn("", "err", err, "value", j)
					}
					results <- p.String()
				}
			}(jobs, results)
		}
		for _, v := range l {
			jobs <- v
		}
		close(jobs)
		for range l {
			p := <-results
			slog.Debug("request", "response done", p.(string))
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

func (m *MetricPost) Tick() {
	for _, delay := range []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second} {
		if m.Post() == nil {
			break
		}
		time.Sleep(delay)
	}
}

