package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run . <temperature><C|F>")
		fmt.Println("example: go run . 32C")
		os.Exit(1)
	}

	inputs := strings.ToUpper(os.Args[1])
	temperature := inputs[:len(inputs)-1]
	scale := inputs[len(inputs)-1:]

	tempConverter(temperature, scale)
}

func tempConverter(temperature string, scale string) {
	temp, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		fmt.Printf("invalid temperature value: %v\n", scale)
		os.Exit(1)
	}
	switch scale {
	case "C":
		fmt.Printf("%vC corresponds to %.2fF\n", temperature, cToFConversion(temp))
	case "F":
		fmt.Printf("%vF corresponds to %.2fC\n", temperature, fToCConversion(temp))
	default:
		fmt.Printf("Invalid scale '%v'. Please use 'C' for Celsius or 'F' for Fahrenheit.\n", scale)
	}
}

func cToFConversion(temperature float64) float64 {
	result := temperature*9/5 + 32
	return result
}

func fToCConversion(temperature float64) float64 {
	result := (temperature - 32) * 5 / 9
	return result
}
