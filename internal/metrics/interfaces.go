package metrics

import (
	"context"
)

type Setter interface {
	Set(ctx context.Context, m Metric) error
}

type Getter interface {
	Get(ctx context.Context, m *Metric) error
}

type ListSetter interface {
	SetList(ctx context.Context, in []Metric) error
}

type ListGetter interface {
	GetList(ctx context.Context) (out []Metric, err error)
}
