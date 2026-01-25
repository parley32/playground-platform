package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/parley-32/playground-platform/internal/fileutil"
)

func main() {
	recursive := flag.Bool("r", false, "recursively list files in subdirectories")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: fileinfo [-r] <path>")
		os.Exit(1)
	}

	path := args[0]

	fi, err := fileutil.GetFileInfo(path)
	if err == nil {
		fmt.Printf("file: %s (%d bytes)\n", fi.Path, fi.Size)
		return
	}

	files, err := fileutil.ListFiles(path, *recursive)

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("no files found in directory")
		return
	} else {
		for _, f := range files {
			fmt.Printf("file: %s (%d bytes)\n", f.Path, f.Size)
		}
	}
}
