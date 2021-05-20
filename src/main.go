package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

func main() {
	folders := make([]string, 0)
	scanPath := flag.String("path", "./", "Path to scan")
	username := flag.String("user", "", "user to scan")
	flag.Parse()
	scanGitFolders(*scanPath, &folders)
	writeSlicesToFile(getStatFilePath(), folders)
	stats(*username)
}

func writeSlicesToFile(filePath string, repos []string) {
	content := strings.Join(repos, "\n")
	ioutil.WriteFile(filePath, []byte(content), 0755)
}

func scanGitFolders(folder string, folders *[]string) {
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
}

func getStatFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	path := usr.HomeDir + "/.gitstats"
	return path
}

func readFileToSlices(filePath string) []string {
	lines := []string{}
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}

		} else {
			log.Fatalln(err)
		}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return lines
}
