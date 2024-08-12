package ticker

import (
	"context"
	"sync"
	"time"
)

type SlaveTicker interface {
	Tick()
}
type Ticker struct {
	wg     sync.WaitGroup
	ticker *time.Ticker
}

func NewAndRun(ctx context.Context, interval uint, slave SlaveTicker) *Ticker {
	t := &Ticker{}
	t.ticker = time.NewTicker(time.Duration(interval) * time.Second)
	t.wg.Add(1)
	go func(ctx context.Context) {
		defer t.ticker.Stop()
		defer t.wg.Done()
		for {
			select {
			case <-t.ticker.C:
				slave.Tick()
			case <-ctx.Done():
				slave.Tick()
				return
			}
		}
	}(ctx)
	return t
}

func (t *Ticker) Close() {
	t.wg.Wait()
}
