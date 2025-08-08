package main

import (
	"flag"
	"fmt"
	"jcloud/api"
	"jcloud/config"
	"jcloud/files"
	"jcloud/storage"
	"os"
)

func main() {
	create := flag.Bool("create", false, "Create new bin")
	update := flag.Bool("update", false, "Update existing bin")
	delete := flag.Bool("delete", false, "Delete bin")
	get := flag.Bool("get", false, "Get bin content")
	list := flag.Bool("list", false, "List all bins")
	file := flag.String("file", "", "Path to JSON file")
	name := flag.String("name", "", "Bin name")
	binID := flag.String("id", "", "Bin ID")

	flag.Parse()

	cfg := config.Load()
	fileSystem := files.FileSystem{}
	storage := storage.NewJsonStorage(fileSystem)
	apiClient := api.New(cfg, storage)

	switch {
	case *create:
		if *file == "" || *name == "" {
			fmt.Println("Error: both --file and --name are required for create")
			os.Exit(1)
		}
		bin, err := apiClient.CreateBin(*file, *name)
		if err != nil {
			fmt.Printf("Error creating bin: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created bin: ID=%s, Name=%s\n", bin.ID, bin.Name)

	case *update:
		if *file == "" || *binID == "" {
			fmt.Println("Error: both --file and --id are required for update")
			os.Exit(1)
		}
		bin, err := apiClient.UpdateBin(*file, *binID)
		if err != nil {
			fmt.Printf("Error updating bin: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Updated bin: ID=%s, Name=%s\n", bin.ID, bin.Name)

	case *delete:
		if *binID == "" {
			fmt.Println("Error: --id is required for delete")
			os.Exit(1)
		}
		if err := apiClient.DeleteBin(*binID); err != nil {
			fmt.Printf("Error deleting bin: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted bin with ID: %s\n", *binID)

	case *get:
		if *binID == "" {
			fmt.Println("Error: --id is required for get")
			os.Exit(1)
		}
		bin, err := apiClient.GetBin(*binID)
		if err != nil {
			fmt.Printf("Error getting bin: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Bin info: ID=%s, Name=%s, CreatedAt=%s\n", bin.ID, bin.Name, bin.CreatedAt)

	case *list:
		binList, err := apiClient.ListBins()
		if err != nil {
			fmt.Printf("Error listing bins: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Your bins:")
		for _, bin := range binList.Bins {
			fmt.Printf("- ID: %s, Name: %s\n", bin.ID, bin.Name)
		}

	default:
		fmt.Println("No valid command specified")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
