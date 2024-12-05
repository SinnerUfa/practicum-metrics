package counter

type Counter struct {
	value int64
}

func (c *Counter) Set(i int64) {
	c.value += i
}

func (c *Counter) Value() int64 {
	return c.value
}
