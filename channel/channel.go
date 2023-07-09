package channel

import (
	"context"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
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

func New(name string, service *youtube.Service) *Channel {
	return &Channel{
		name:         name,
		service:      service,
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
			theme := "Spain"
			duration := uint64(601)
			tasks <- Task{Theme: theme, NeedDurationSec: duration, Result: c.videoResults}
			// нужен таумаут не вечное ожидание
			result := <-c.videoResults
			log.Printf("full video result %+v", result)

			upload := &youtube.Video{
				Snippet: &youtube.VideoSnippet{
					Title:       "beautiful " + theme,
					Description: "in this video beautiful " + theme, // can not use non-alpha-numeric characters
					CategoryId:  "22",
				},
				Status: &youtube.VideoStatus{PrivacyStatus: "unlisted"},
			}

			// The API returns a 400 Bad Request response if tags is an empty string.
			upload.Snippet.Tags = []string{theme, "upload", "api"}

			call := c.service.Videos.Insert([]string{"snippet", "status"}, upload)
			file, err := os.Open(result.Data.Path)
			if err != nil {
				log.Fatalf("Error opening %v: %v", result.Data.Path, err)
			}
			defer file.Close()

			response, err := call.Media(file).Do()
			if err != nil {
				log.Fatalf("Error making YouTube API call: %v", err)
			}
			fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
