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

type SetterWithContext interface {
	SetWithContext(ctx context.Context, m Metric) error
}

type GetterWithContext interface {
	GetWithContext(ctx context.Context, m *Metric) error
}

type ListSetterWithContext interface {
	SetListWithContext(ctx context.Context, in []Metric) error
}

type ListGetterWithContext interface {
	GetListWithContext(ctx context.Context) (out []Metric, err error)
}
