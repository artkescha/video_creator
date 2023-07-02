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
	videoCreator := creator.New(cli, saver.VideoDownloader{}, rootVideoPath)
	videoCreator.Start(ctx, 10*time.Second)

	// handle ctr+c.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	cancel()
}
