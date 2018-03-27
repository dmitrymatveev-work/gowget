package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

var progress = make(map[string]float64)
var mutex sync.RWMutex
var urls []string

func main() {
	urls = os.Args[1:]

	if len(urls) == 0 {
		fmt.Println("Please provide URLs list.")
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	var header string
	for _, url := range urls {
		header = fmt.Sprintf("%s\t%s", header, path.Base(url))
		progress[url] = 0
	}
	fmt.Println(header)

	go func() {
		for {
			printProgress()
			time.Sleep(time.Second)
		}
	}()

	for _, url := range urls {
		go download(url, &wg)
	}

	wg.Wait()
	printProgress()
}

func printProgress() {
	var row string
	for _, url := range urls {
		mutex.Lock()
		row = fmt.Sprintf("%s\t%3.f%%", row, progress[url])
		mutex.Unlock()
	}
	fmt.Println(row)
}

func download(url string, wg *sync.WaitGroup) {
	fileName := path.Base(url)
	out, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	counter := &Counter{
		UpdateProgress: func(s int) {
			mutex.RLock()
			progress[url] += float64(s) / float64(size) * 100
			mutex.RUnlock()
		},
	}

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		log.Fatal(err)
	}

	wg.Done()
}

type Counter struct {
	UpdateProgress func(int)
}

func (c *Counter) Write(chunk []byte) (int, error) {
	n := len(chunk)
	c.UpdateProgress(n)
	return n, nil
}
