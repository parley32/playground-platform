package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: fileinfo <path>")
		os.Exit(1)
	}

	fmt.Println("path:", os.Args[1])
}