package ticker

import (
	"context"
	"time"
)

type SlaveTicker interface {
	Tick() error
}

func NewAndRun(ctx context.Context, interval uint, s SlaveTicker) {
	go func(uint, SlaveTicker) {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.Tick()
			case <-ctx.Done():
				s.Tick()
				return
			}
		}
	}(interval, s)
}
