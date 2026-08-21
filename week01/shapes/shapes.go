package main

import (
	"errors"
	"fmt"
	"math"
)

type Shape interface {
	area() (float64, error)
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

func (c Circle) area() (float64, error) {
	if c.radius <= 0 {
		return 0, errors.New("radius must be greater than 0")
	}
	return c.radius * c.radius * math.Pi, nil
}
func (r Rect) area() (float64, error) {
	if r.height <= 0 || r.width < 0 {
		return 0, errors.New("height and width must be greater than 0")
	}
	return r.height * r.width, nil
}

func (t Tri) area() (float64, error) {
	if t.base <= 0 || t.height < 0 {
		return 0, errors.New("base and height must be greater than 0")
	}
	return t.height * t.base / 2, nil
}

func totalArea(shapes []Shape) (float64, error) {
	var total float64
	for _, shape := range shapes {
		area, err := shape.area()
		if err != nil {
			return 0, fmt.Errorf("could not calculate area for %T: %w", shape, err)
		}
		total += area
	}
	return total, nil
}
