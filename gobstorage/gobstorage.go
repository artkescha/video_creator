package gobstorage

import (
	"fmt"
	"os"

	"encoding/gob"
	"path/filepath"
)

const (
	path      = "/gob"
	extension = ".gob"
)

type GobStorage struct {
	rootPath string
}

func NewGobStorage() *GobStorage {
	return &GobStorage{rootPath: filepath.Dir(os.Args[0])}
}

func (s GobStorage) Load(name string, data interface{}) error {
	filename := filepath.Join(s.rootPath, path, name+extension)
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open file failed: %s", err)
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("decode data from file failed: %s", err)
	}
	return nil
}

func (s GobStorage) Save(name string, data interface{}) error {
	dirname := filepath.Join(s.rootPath, path)
	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		return fmt.Errorf("make dir filed: %s", err)
	}
	filename := filepath.Join(dirname, name+extension)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file filed: %s", err)
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode data to file failed: %s", err)
	}
	return nil
}
