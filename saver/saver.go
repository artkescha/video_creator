package saver

import (
	"os"
	"os/exec"
)

type Saver interface {
	DownloadVideo(url string, filepath string) error
}

type VideoDownloader struct{}

func (d VideoDownloader) DownloadVideo(url string, filepath string) error {
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
