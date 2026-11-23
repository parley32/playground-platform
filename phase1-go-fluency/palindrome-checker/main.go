package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run . <word>")
		fmt.Println("example: go run . nun")
		os.Exit(1)
	}

	word := os.Args[1]
	n := len(word) - 1
	for i := 0; i < n; i++ {
		if word[i] != word[n] {
			fmt.Printf("%v is not a palindrome\n", word)
			os.Exit(1)
		}
		n--
	}
	fmt.Printf("Congrats! %v is a palindrome\n", word)
}
