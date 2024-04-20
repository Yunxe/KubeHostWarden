package common

import (
	"context"
	"time"
)

type Collector struct {
	MetricType string
	ErrCh      chan error
	PointCh    chan Point
}

type MetricConfig struct {
	Type        string
	CollectFunc func(ctx context.Context, collector *Collector)
}

func NewCollector(mt string) Collector {
	return Collector{
		MetricType: mt,
		ErrCh:      make(chan error, 1),
		PointCh:    make(chan Point, 10),
	}
}

func (c *Collector) ReturnPonit() <-chan Point {
	return c.PointCh
}

func (c *Collector) ReturnError() <-chan error {
	return c.ErrCh
}

func (c *Collector) Start(ctx context.Context, collectFunc func(ctx context.Context, c *Collector)) {
	go func() {
		ticker := time.NewTicker(3 * time.Second) // Adjust the interval as necessary
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				close(c.PointCh)
				close(c.ErrCh)
				return
			case <-ticker.C:
				collectFunc(ctx, c)
			}
		}
	}()
}
