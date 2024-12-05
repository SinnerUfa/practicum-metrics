package metrics

import (
	"context"
)

type Setter interface {
	Set(m Metric) error
}

type Getter interface {
	Get(m *Metric) error
}

type ListSetter interface {
	SetList([]Metric) error
}

type ListGetter interface {
	GetList() (out []Metric, err error)
}

type ContextSetter interface {
	SetContext(ctx context.Context, m Metric) error
}

type ContextGetter interface {
	GetContext(ctx context.Context, m *Metric) error
}

type ContextListSetter interface {
	SetListContext(ctx context.Context, in []Metric) error
}

type ContextListGetter interface {
	GetListContext(ctx context.Context) (out []Metric, err error)
}
