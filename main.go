package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var statuses map[string]int = make(map[string]int)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide URLs list.")
	}

	var wg sync.WaitGroup
	wg.Add(len(args))

	for _, a := range args {
		statuses[a] = 0
	}

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
	statuses[url] = 100
	time.Sleep(2*time.Second)
	wg.Done()
}