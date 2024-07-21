package ticker

import (
	"context"
	"time"
)

type SlaveTicker interface {
	Tick()
}

func NewAndRun(ctx context.Context, interval uint, s SlaveTicker) {
	go func(intv uint, slave SlaveTicker) {
		ticker := time.NewTicker(time.Duration(intv) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				slave.Tick()
			case <-ctx.Done():
				slave.Tick()
				return
			}
		}
	}(interval, s)
}
