package main

import (
	"math"
)

type Shape interface {
	area() float64
}

type Circle struct {
	radius float64
}

type Rect struct {
	width, height float64
}

type Tri struct {
	base, height float64
}

func (c Circle) area() float64 {
	return c.radius * c.radius * math.Pi
}
func (r Rect) area() float64 {
	return r.height * r.width
}

func (t Tri) area() float64 {
	return t.height * t.base / 2
}

func totalArea(g []Shape) float64 {
	var total float64
	for _, shape := range g {
		total += shape.area()
	}
	return total
}
