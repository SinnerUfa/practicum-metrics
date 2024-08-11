package metrics

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
