package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide URLs list.")
	}

	var wg sync.WaitGroup
	wg.Add(len(args))

	go func() {
		for _, a := range os.Args[1:] {
			fmt.Println(a)
			wg.Done()
		}
	}()

	wg.Wait()
}