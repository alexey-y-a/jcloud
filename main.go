package main

import (
	"fmt"
	"jcloud/api"
	"jcloud/bins"
	"jcloud/config"
	"jcloud/files"
	"jcloud/storage"
)

func main() {
	cfg := config.Load()

	apiClient := api.New(cfg)
	apiClient.MakeRequest()

	fileSystem := files.FileSystem{}
	storage := storage.NewJsonStorage(fileSystem)

	binList := bins.NewBinList()
	binList.Bins = append(binList.Bins, bins.NewBin("1", false, "Test Bin"))

	err := storage.SaveBins(binList, "bins.json")
	if err != nil {
		fmt.Printf("Error saving bins: %v\n", err)
		return
	}

	loadedBins, err := storage.LoadBins("bins.json")
	if err != nil {
		fmt.Printf("Error loading bins: %v\n", err)
		return
	}

	fmt.Printf("Loaded bins: %+v\n", loadedBins)
}
