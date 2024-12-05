package gauge

type Gauge struct {
	value float64
}

func (c *Gauge) Set(i float64) {
	c.value = i
}

func (c *Gauge) Value() float64 {
	return c.value
}
