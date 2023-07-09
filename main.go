package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"video_creator/channel"

	"github.com/kosa3/pexels-go"
	"video_creator/creator"
	"video_creator/saver"
)

func main() {
	cli := pexels.NewClient(os.Args[1])
	proxy, err := url.Parse("http://201.91.82.155:3128")
	if err != nil {
		log.Fatal(err)
		return
	}
	cli.HTTPClient = &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
		Timeout:   20 * time.Second,
	}
	ctx, cancel := context.WithCancel(context.Background())
	rootVideoPath := filepath.Join(".", "videos")
	tasksChan := make(chan channel.Task)
	videoCreator := creator.New(cli, saver.VideoDownloader{}, rootVideoPath)
	videoCreator.Start(ctx, tasksChan)
	interval := 288 * time.Minute // 4.8 часа (5 видео в сутки)
	channel.New("channel_1").Start(ctx, interval, tasksChan)

	// handle ctr+c.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	cancel()
}
