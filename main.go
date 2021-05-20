package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	count := flag.Int("a", 5, "Test")
	flag.Parse()
	fmt.Println(*count)
	fmt.Println(flag.Args())
}

func scanGitFolders(folder string, folders *[]string) []string {
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	currentPath := ""
	for _, file := range files {
		if file.IsDir() {
			currentPath = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path := strings.TrimSuffix(currentPath, "/.git")
				*folders = append(*folders, path)
			} else {
				if file.Name() == "node_modules" || file.Name() == "vendor" {
					continue
				}
				scanGitFolders(currentPath, folders)
			}
		}
	}

	return *folders
}
