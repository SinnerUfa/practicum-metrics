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

type ListGetter interface {
	GetList() (out []Metric)
}

type ListSetter interface {
	SetList([]Metric) error
}

type ContextSetter interface {
	SetContext(ctx context.Context, m Metric) error
}

type ContextGetter interface {
	GetContext(ctx context.Context, m *Metric) error
}

type ContextListGetter interface {
	GetListContext(ctx context.Context) (out []Metric)
}

type ContextListSetter interface {
	SetListContext(ctx context.Context, in []Metric) error
}
