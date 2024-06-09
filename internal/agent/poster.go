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

	for _, v := range l {
		endpoint := "http://" + m.adress + "/" + v.ReguestString("update")
		// request, err := http.NewRequestWithContext(m.ctx, http.MethodPost, endpoint, http.NoBody)
		// if err != nil {
		// 	return err
		// }
		// request.Header.Add("Content-Type", "text/plain")

		// response, err := client.Do(request)
		// if err != nil {
		// 	return nil // нужно подумать
		// }
		// defer response.Body.Close()
		// _, err = io.ReadAll(response.Body)
		// if err != nil {
		// 	return nil
		// }
		_, err := client.R().SetHeader("Content-Type", "text/plain").Post(endpoint)
		if err != nil {
			return nil
		}

	}
	m.log.Info("Post increment:", m.counter)
	return nil
}

func (m *MetricPost) Tick() error {
	return m.Post()
}
