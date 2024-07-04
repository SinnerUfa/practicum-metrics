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
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.Tick()
			}
		}
	}(interval, s)
}
