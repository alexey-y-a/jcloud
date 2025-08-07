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
	// verbose := flag.Bool("v", false, "Verbose output")

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
		fmt.Println("Update command not implemented yet")
		os.Exit(0)

	case *delete:
		if *binID == "" {
			fmt.Println("Error: --id is required for delete")
			os.Exit(1)
		}
		fmt.Println("Delete command not implemented yet")
		os.Exit(0)

	case *get:
		if *binID == "" {
			fmt.Println("Error: --id is required for get")
			os.Exit(1)
		}
		fmt.Println("Get command not implemented yet")
		os.Exit(0)

	case *list:
		fmt.Println("List command not implemented yet")
		os.Exit(0)

	default:
		fmt.Println("No valid command specified")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
