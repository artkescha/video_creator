package sender

import (
	"bytes"
	"fmt"
	"os/exec"
)

type YoutubeSender struct {
	command string
}

func New(command string) *YoutubeSender {
	return &YoutubeSender{command: command}
}

func (s YoutubeSender) Send(filename, title, description string) (string, error) {
	cmd := exec.Command(s.command, "-filename",
		filename, "-title", title, "-description", description)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, stderr.String())
	}
	return fmt.Sprint("upload video success: %s" + out.String()), nil
}
