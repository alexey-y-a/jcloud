package storage

import (
	"encoding/json"
	"errors"
	"jcloud/bins"
	"jcloud/files"
)

type Storage interface {
	SaveBins(binList bins.BinList, filename string) error
	LoadBins(filename string) (bins.BinList, error)
}

type JsonStorage struct {
	Files files.FileOperations
}

func NewJsonStorage(f files.FileOperations) *JsonStorage {
	return &JsonStorage{Files: f}
}

func (s *JsonStorage) SaveBins(binList bins.BinList, filename string) error {
	if !s.Files.IsJson(filename) {
		return errors.New("file must have .json extension")
	}

	data, err := json.MarshalIndent(binList, "", "  ")
	if err != nil {
		return err
	}

	return s.Files.WriteFile(filename, data, 0644)
}

func (s *JsonStorage) LoadBins(filename string) (bins.BinList, error) {
	var binList bins.BinList

	if !s.Files.IsJson(filename) {
		return binList, errors.New("file must have .json extension")
	}

	data, err := s.Files.ReadFile(filename)
	if err != nil {
		return binList, err
	}

	err = json.Unmarshal(data, &binList)
	if err != nil {
		return binList, err
	}

	return binList, nil
}
