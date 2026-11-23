package main

import (
	"fmt"
	"os"
	"strconv"
)

type BMICategory string

const (
	Underweight BMICategory = "underweight"
	Healthy     BMICategory = "healthy"
	Overweight  BMICategory = "overweight"
	Obese       BMICategory = "obese"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: go run . <height in cm> <weight in kg>")
		fmt.Println("example: go run . 160 100")
		os.Exit(1)
	}

	height, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil || height <= 0 {
		fmt.Printf("%v is not a valid height\n", os.Args[1])
		os.Exit(1)
	}
	heightInM := height / 100
	weight, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil || weight <= 0 {
		fmt.Printf("%v is not a valid weight\n", os.Args[2])
		os.Exit(1)
	}

	bmi := weight / (heightInM * heightInM)
	fmt.Printf("If you measure %vm and weigh %vkg, your BMI is %.1f\n", heightInM, weight, bmi)
	printFindings(classifyBMI(bmi))
}

func classifyBMI(bmi float64) BMICategory {
	switch {
	case bmi < 18.5:
		return Underweight
	case bmi < 25:
		return Healthy
	case bmi < 30:
		return Overweight
	default:
		return Obese
	}
}

func printFindings(bmi BMICategory) {
	switch bmi {
	case Underweight:
		fmt.Println("You are UNDERWEIGHT. Get to bulkin', before some child decides to toss you around like the stick figure that you are")
	case Healthy:
		fmt.Println("Hmmm. 'Healthy', as in 'Average'. GET TO BULKIN'!!")
	case Overweight:
		fmt.Println("Quite litteraly the definition of peak raw athleticism. Remember: if you can see your own abs, it just means you are underfed.")
	case Obese:
		fmt.Println("Haters call it Obesity, but really you're just a professional-grade bulker. Also, should you be taking criticism from anyone you could easily defeat by sitting on?")
	}
}
