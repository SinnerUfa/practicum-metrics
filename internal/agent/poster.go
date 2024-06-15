package agent

import (
	"context"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	resty "github.com/go-resty/resty/v2"
)

type MetricPost struct {
	log     mlog.Logger
	rep     repository.Repository
	ctx     context.Context
	adress  string
	counter uint
}

func NewPoster(ctx context.Context, log mlog.Logger, rep repository.Repository, adress string) *MetricPost {
	return &MetricPost{ctx: ctx, log: log, rep: rep, adress: adress}
}

func (m *MetricPost) Post() error {
	l := m.rep.List()
	client := resty.New()

	endpoint := "http://" + m.adress + "/"

	for _, v := range l {

		_, err := client.R().SetHeader("Content-Type", "text/plain").Post(endpoint + v.ReguestString("update"))
		if err != nil {
			m.log.Warning(err)
			return err
		}

	}
	m.log.Info("Post increment:", m.counter)
	return nil
}

func (m *MetricPost) Tick() error {
	return m.Post()
}
