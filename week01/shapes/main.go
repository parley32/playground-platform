package main

import (
	"fmt"
)

func main() {
	circleA := Circle{radius: 4}
	circleB := Circle{radius: 2}
	rectA := Rect{height: 2, width: 4}
	rectB := Rect{height: 5, width: 5}
	triA := Tri{base: 5, height: 3}
	triB := Tri{base: 2, height: 2}

	shapeMassA := []Shape{circleA, rectA, triA}
	shapeMassB := []Shape{circleB, rectB, triB}

	fmt.Printf("shape %v has an area of size %.04f\n", circleA, circleA.area())
	fmt.Printf("shape %v has an area of size %.04f\n", rectA, rectA.area())
	fmt.Printf("shape %v has an area of size %.04f\n", triA, triA.area())
	fmt.Printf("shapeMass %v has an area of size %.04f\n", shapeMassA, totalArea(shapeMassA))

	fmt.Printf("shape %v has an area of size %.04f\n", rectB, rectB.area())
	fmt.Printf("shape %v has an area of size %.04f\n", circleB, circleB.area())
	fmt.Printf("shape %v has an area of size %.04f\n", triB, triB.area())
	fmt.Printf("shapeMass %v has an area of size %.04f\n", shapeMassB, totalArea(shapeMassB))

}
