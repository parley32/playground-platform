package main

import (
	"fmt"
)

func main() {
	circleA := Circle{radius: -4}
	circleB := Circle{radius: 2}
	rectA := Rect{height: 2, width: 4}
	rectB := Rect{height: 5, width: 5}
	triA := Tri{base: 5, height: 3}
	triB := Tri{base: 2, height: 2}

	shapeMassA := []Shape{circleA, rectA, triA}
	shapeMassB := []Shape{circleB, rectB, triB}
	areaA, err := totalArea(shapeMassA)
	if err != nil {
		fmt.Println("could not calculate shapeMassA:", err)
		return
	}
	fmt.Printf("shapeMass %v has an area of size %.04f\n", shapeMassA, areaA)

	areaB, err := totalArea(shapeMassB)
	if err != nil {
		fmt.Println("could not calculate shapeMassB:", err)
		return
	}
	fmt.Printf("shapeMass %v has an area of size %.04f\n", shapeMassB, areaB)

}
