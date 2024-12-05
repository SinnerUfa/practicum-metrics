package jmetric

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func New() *Metrics {
	var (
		a int64
		b float64
	)
	return &Metrics{
		ID:    "",
		MType: "",
		Delta: &a,
		Value: &b,
	}
}
