package channel

import (
	"context"
	"log"
	"time"
)

type Channel struct {
	name         string
	videoResults chan VideoResult
}

type Task struct {
	Theme           string
	NeedDurationSec uint64
	Result          chan VideoResult
}

type Data struct {
	Path     string
	Duration float64
}

type VideoResult struct {
	Data *Data
	Err  error
}

func New(name string) *Channel {
	return &Channel{
		name:         name,
		videoResults: make(chan VideoResult),
	}
}

func (c Channel) Start(ctx context.Context, interval time.Duration, tasks chan Task) {
	go c.run(ctx, interval, tasks)
}

func (c Channel) run(ctx context.Context, interval time.Duration, tasks chan Task) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			tasks <- Task{Theme: "Spain", NeedDurationSec: 601, Result: c.videoResults}
			// нужен таумаут не вечное ожидание
			result := <-c.videoResults
			log.Println("full video result %+v", result)
			// TDOD написать логику отправки тасок креатору
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
