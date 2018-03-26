package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide URLs list.")
	}

	for _, a := range os.Args[1:] {
		fmt.Println(a)
	}
}