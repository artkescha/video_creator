package creator

import (
	"context"
	"fmt"
	"github.com/artkescha/moviego"
	"github.com/kosa3/pexels-go"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"video_creator/marker"
	"video_creator/saver"
)

type VideoCreator struct {
	pexlsClient   *pexels.Client
	markers       *marker.Markers
	saver         saver.Saver
	videoRootPath string
}

func New(client *pexels.Client, saver saver.Saver, videoRootPath string) *VideoCreator {
	return &VideoCreator{
		pexlsClient:   client,
		markers:       marker.NewMarkers("./videos"),
		saver:         saver,
		videoRootPath: videoRootPath,
	}
}

func (creator *VideoCreator) Start(ctx context.Context, createVideoInterval time.Duration) {
	go creator.run(ctx, createVideoInterval)
}

func (creator *VideoCreator) run(ctx context.Context, createVideoInterval time.Duration) {
	ticker := time.NewTicker(createVideoInterval)
	// TODO раскомментировать при дальнейшей разработке!
	for {
		select {
		case <-ticker.C:
			theme := "Praga"
			themePath := filepath.Join(creator.videoRootPath, theme)
			if err := createFolder(themePath); err != nil {
				log.Printf("create theme folder failed %s", err)
				break
			}
			videoPath := filepath.Join(themePath, time.Now().String())
			if err := createFolder(videoPath); err != nil {
				log.Printf("create video folder failed %s", err)
				break
			}
			log.Printf("SA videoPath %s", videoPath)
			err := creator.createVideo(ctx, theme, videoPath, 1200)
			if err != nil {
				log.Println("create video fail fail", err)
			}

		case <-ctx.Done():
			ticker.Stop()
			log.Println("video creator stop")
			return
		}
	}
}

func (creator *VideoCreator) createVideo(ctx context.Context, theme string, basePath string, needDuration int) error {
	currentDuration := 0
	partsPaths := make([]string, 0)
	currentPage, err := creator.markers.Get(theme)
	log.Printf("current page: %d", currentPage)

	if err != nil {
		return fmt.Errorf("load current page number with teme %s failed %w", theme, err)
	}

	for currentDuration < needDuration {
		// fix for many request
		log.Println("fix for many requests, wait 121 seconds before next request")
		// time.Sleep(121 * time.Second)

		vs, err := creator.pexlsClient.VideoService.Search(ctx, &pexels.VideoParams{
			Query: theme,
			Page:  currentPage,
		})
		if err != nil {
			log.Println("search failed ", err)
			if strings.HasSuffix(err.Error(), "with status code 429") {
				log.Println("many request - wait 1 Hour")
				time.Sleep(1 * time.Hour)
			}
			continue
		}
		for _, video := range vs.Videos {
			path := ""
			for idx, file := range video.VideoFiles {
				if file.Quality != "sd" || file.FPS != 25.00 || file.Width != 960 || file.Height != 540 || file.FileType != "video/mp4" {
					continue
				}
				log.Println(file, idx)
				path = filepath.Join(basePath, fmt.Sprintf("%s-%d-%d.mp4", theme, time.Now().UnixNano(), idx))
				err := creator.saver.DownloadVideo(file.Link, path)
				if err != nil {
					log.Printf("download video part with path %s failed %s", path, err)
					break
				}
				partsPaths = append(partsPaths, path)
				currentDuration += video.Duration
				break
			}
		}
		if currentDuration >= needDuration {
			currentPage++
			log.Printf("duration is completed %d", currentDuration)
			if err := creator.markers.Set(theme, currentPage); err != nil {
				return fmt.Errorf("save current page number with teme %s failed %w", theme, err)
			}
			break
		}
		currentPage++
		if err := creator.markers.Set(theme, currentPage); err != nil {
			return fmt.Errorf("save current page number with teme %s failed %w", theme, err)
		}
	}
	fullVideo, err := createVideoFromParts(partsPaths, basePath)
	if err != nil {
		return err
	}
	log.Printf("SA full video %+v", fullVideo)
	return nil
}

func createVideoFromParts(partsPaths []string, basePath string) (moviego.Video, error) {
	videos := make([]moviego.Video, 0)
	log.Println("partsPaths", partsPaths)
	for _, filePath := range partsPaths {
		videoPart, err := moviego.Load(filePath)
		if err != nil {
			return moviego.Video{}, err
		}
		videos = append(videos, videoPart)
	}
	// Объединение нескольких видео в одно.
	finalVideo, err := moviego.Concat(videos, basePath)
	if err != nil {
		log.Printf("moviego.Concat %s", err)
		return moviego.Video{}, err
	}
	//renderErr := finalVideo.Output("final.mp4").Run()
	//if err != nil {
	//	return moviego.Video{}, renderErr
	//}
	// removed parts video
	for _, filePath := range partsPaths {
		if err := os.Remove(filePath); err != nil {
			log.Printf("removed file with path %s failed %s", filePath, err)
			continue
		}
	}
	return finalVideo, nil
}

func createFolder(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
