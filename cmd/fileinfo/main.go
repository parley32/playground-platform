package main

import (
	"fmt"
	"os"

	"github.com/parley-32/playground-platform/internal/fileutil"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: fileinfo <path>")
		os.Exit(1)
	}

	path := os.Args[1]

	fi, err := fileutil.GetFileInfo(path)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	ok, err := fileutil.IsRegularFile(path)

	if !ok {
		fmt.Println("not a regular file")
		os.Exit(1)
	}

	fmt.Printf("file: %s (%d bytes)\n", fi.Path, fi.Size)

}
