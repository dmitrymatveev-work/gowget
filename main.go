package main

import (
	"io"
	"strconv"
	"net/http"
	"log"
	"fmt"
	"os"
	"sync"
	"time"
	"path"
)

var statuses map[string]float64 = make(map[string]float64)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide URLs list.")
	}

	var wg sync.WaitGroup
	wg.Add(len(args))

	var header string
	for _, a := range args {
		header = fmt.Sprintf("%s\t%s", header, path.Base(a))
		statuses[a] = 0
	}
	fmt.Println(header)

	go func() {
		for {
			printStatus()
			time.Sleep(time.Second)
		}
	}()

	for _, a := range args {
		go download(a, &wg)
	}

	wg.Wait()
	printStatus();
}

func printStatus() {
	var result string
	for _, s := range statuses {
		result = fmt.Sprintf("%s\t%d%%", result, int(s))
	}
	fmt.Println(result)
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

	counter := &Counter {
		UpdateStatus: func(s int) {
			statuses[url] += float64(s) / float64(size) * 100
		},
	}

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		log.Fatal(err)
	}

	wg.Done()
}

type Counter struct {
	UpdateStatus func(int)
}

func (c *Counter) Write(p []byte) (int, error) {
	n := len(p)
	c.UpdateStatus(n)
	return n, nil
}