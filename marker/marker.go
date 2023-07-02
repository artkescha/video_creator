package marker

import (
	"fmt"
	"log"
	"video_creator/gobstorage"
)

type Markers struct {
	fileName   string
	markers    map[string]int
	gobStorage *gobstorage.GobStorage
}

func NewMarkers(fileName string) *Markers {
	data := make(map[string]int)
	gob := gobstorage.NewGobStorage()
	err := gob.Load(fileName, &data)
	if err != nil {
		log.Printf("load markers with filename %s failed %s", fileName, err)
	}
	return &Markers{
		fileName:   fileName,
		markers:    data,
		gobStorage: gob,
	}
}

func (m *Markers) Get(theme string) (int, error) {
	return m.markers[theme], nil
}

func (m *Markers) Set(theme string, currentPage int) error {
	m.markers[theme] = currentPage
	if err := m.gobStorage.Save(m.fileName, m.markers); err != nil {
		return fmt.Errorf("save Markers: %w", err)
	}
	return nil
}
