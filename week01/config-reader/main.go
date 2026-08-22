package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run . <file>")
		os.Exit(1)
	}

	file, err := openFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error whilst opening the file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	content, err := readContent(file)
	if err != nil {
		fmt.Printf("Error whilst reading the file: %v\n", err)
		os.Exit(1)
	}

	result, err := parseNumber(content)
	if err != nil {
		fmt.Printf("Error whilst parsing the file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("The in-file number is: %v\n", result)

}

func openFile(fileName string) (*os.File, error) {
	return os.Open(fileName)
}

func readContent(file *os.File) ([]byte, error) {
	return os.ReadFile(file.Name())
}

func parseNumber(content []byte) (int, error) {
	result := int(0)
	isNegative := false
	for _, char := range content {
		switch {
		case unicode.IsSpace(rune(char)):
			{
				continue
			}
		case result == 0 && rune(char) == '+':
			{

				isNegative = false

			}
		case result == 0 && rune(char) == '-':
			{

				isNegative = true
			}
		case unicode.IsDigit(rune(char)):
			{
				digit, err := strconv.Atoi(string(char))
				if err != nil {
					return 0, fmt.Errorf("invalid digit %q: %w", char, err)
				}
				result = result*10 + digit
			}
		default:
			{
				return 0, fmt.Errorf("char %q is not a number", char)
			}
		}
	}
	if isNegative {
		result *= -1
	}
	return result, nil
}
