package channel

import (
	"context"
	"fmt"
	"log"
	"time"
	"video_creator/sender"
)

type Channel struct {
	name         string
	videoResults chan VideoResult
	sender       sender.Sender
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

func New(name string, sender sender.Sender) *Channel {
	return &Channel{
		name:         name,
		videoResults: make(chan VideoResult),
		sender:       sender,
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
			theme := "Kosmos"
			duration := uint64(601)
			tasks <- Task{Theme: theme, NeedDurationSec: duration, Result: c.videoResults}
			// нужен таумаут, не вечное ожидание
			result := <-c.videoResults
			log.Printf("full video result %+v", result)
			if result.Err != nil {
				fmt.Printf("create video failed %s", result.Err)
				continue
			}
			output, err := c.sender.Send(result.Data.Path, "Beautiful "+theme, "Beautiful "+theme)
			if err != nil {
				fmt.Printf("send video with path %s failed %s", result.Data.Path, err)
				continue
			}
			fmt.Printf("send video with path %s result %s", result.Data.Path, output)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
