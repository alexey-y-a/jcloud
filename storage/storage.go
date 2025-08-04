package storage

import (
	"encoding/json"
	"errors"
	"jcloud/bins"
	"jcloud/files"
	"os"
)

func SaveBins(binList bins.BinList, filename string) error {
	if !files.IsJason(filename) {
		return errors.New("file must have .json extension")
	}

	data, err := json.MarshalIndent(binList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadBins(filename string) (bins.BinList, error) {
	var binList bins.BinList

	if !files.IsJason(filename) {
		return binList, errors.New("file must have .json extension")
	}

	data, err := files.ReadFile(filename)
	if err != nil {
		return binList, err
	}

	err = json.Unmarshal(data, &binList)
	if err != nil {
		return binList, err
	}

	return binList, nil

}
