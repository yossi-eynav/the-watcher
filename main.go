package main

import (
	"flag"
	"fmt"
	"time"

	"strings"

	"github.com/yossi-eynav/the-watcher/lib"
)

const minutesInterval = 1

func setOptions() (string, []string) {
	repositoryPath := flag.String("repositoryPath", "", "a string")
	files := flag.String("files", "", "a string")
	flag.Parse()
	if len(*repositoryPath) == 0 {
		panic("Please enter a valid repositoryPath")
	}
	return *repositoryPath, strings.Split(*files, ",")
}

func main() {
	fmt.Println("Make sure you installed this tool before - https://github.com/chbrown/osx-notifier")
	ch := make(chan string)
	repositoryPath, fileNames := setOptions()
	fmt.Printf("Running with these options:\nrepositoryPath=%s, fileNames=%s", repositoryPath, fileNames)
	Git := lib.Git{RepositoryPath: repositoryPath}

	go func() {
		for {
			Git.Fetch()
			for _, filename := range fileNames {
				go Git.Log(ch, filename)
			}
			time.Sleep(time.Minute * minutesInterval)
		}
	}()

	for {
		Git.LogParser(<-ch)
	}
}
