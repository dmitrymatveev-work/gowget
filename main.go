package main

import (
	"log"
	"fmt"
	"os"
	"sync"
	"time"
	"path"
)

var statuses map[string]int = make(map[string]int)

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
}

func printStatus() {
	var result string
	for _, v := range statuses {
		result = fmt.Sprintf("%s\t%d%%", result, v)
	}
	fmt.Println(result)
}

func download(url string, wg *sync.WaitGroup) {
	fileName := path.Base(url)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	statuses[url] = 100
	time.Sleep(2*time.Second)

	f.WriteString("Test")

	wg.Done()
}