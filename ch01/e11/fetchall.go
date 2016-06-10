package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	urls := os.Args[1:]

	start := time.Now()

	progress := asyncDownload(urls)
	waitAndPrint(progress, len(urls))

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func asyncDownload(urls []string) chan string {
	progress := make(chan string)
	for i, url := range urls {
		savePath := fmt.Sprintf("download-%06d", i)
		go fetchAndSave(url, savePath, progress)
	}

	return progress
}

func fetchAndSave(url, path string, progress chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		progress <- fmt.Sprint(err)
		return
	}

	nbytes, err := save(path, resp.Body)
	errClose := resp.Body.Close()
	if err == nil {
		err = errClose
	}
	if err != nil {
		progress <- fmt.Sprintf("while saving %s: %v\n", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	progress <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}

func save(path string, contents io.Reader) (nr int64, err error) {
	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}

	defer func() {
		errClose := file.Close()
		if err == nil {
			err = errClose
		}
	}()

	return io.Copy(file, contents)
}

func waitAndPrint(ch <-chan string, n int) {
	for i := 0; i < n; i++ {
		fmt.Println(<-ch) // receive from channel ch
	}
}
