package main

import (
	"fmt"
	"jcloud/bins"
	"jcloud/storage"
	"log"
)

func main() {

	bin1 := bins.NewBin("123", true, "Test Bin 1")
	bin2 := bins.NewBin("456", false, "Test Bin 2")

	binList := bins.NewBinList()
	binList.Bins = append(binList.Bins, bin1, bin2)

	err := storage.SaveBins(binList, "bins.json")
	if err != nil {
		log.Fatal("Save error:", err)
	}

	loadedBins, err := storage.LoadBins("bins.json")
	if err != nil {
		log.Fatal("Load error:", err)
	}

	fmt.Printf("Loaded bins: %+v\n", loadedBins)
}
