package channel

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"log"
	"os/exec"
	"time"
)

type Channel struct {
	name         string
	videoResults chan VideoResult
	service      *youtube.Service
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
			theme := "Ocean"
			duration := uint64(601)
			tasks <- Task{Theme: theme, NeedDurationSec: duration, Result: c.videoResults}
			// нужен таумаут не вечное ожидание
			result := <-c.videoResults
			log.Printf("full video result %+v", result)
			if result.Err != nil {
				fmt.Printf("create video failed %s", result.Err)
				continue
			}
			cmd := exec.Command("./youtubeuploader", "-filename",
				result.Data.Path, "-title", "Beautiful "+theme, "-description", "Beautiful "+theme)
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				continue
			}
			fmt.Println("upload video success: " + out.String())
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
