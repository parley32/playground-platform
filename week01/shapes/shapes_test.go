package main

import (
	"math"
	"testing"
)

func TestCircleArea(t *testing.T) {
	tests := []struct {
		name     string
		radius   float64
		wantArea float64
		wantErr  bool
	}{
		{
			name:     "valid circle",
			radius:   2,
			wantArea: 4 * math.Pi,
		},
		{
			name:    "invalid radius",
			radius:  -1,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			circle := Circle{radius: test.radius}
			got, err := circle.area()
			if (err != nil) != test.wantErr {
				t.Fatalf("area() error = %v, wantErr %v", err, test.wantErr)
			}
			if got != test.wantArea {
				t.Errorf("area() = %v, want %v", got, test.wantArea)
			}
		})
	}
}

func TestRectangleArea(t *testing.T) {
	tests := []struct {
		name     string
		width    float64
		height   float64
		wantArea float64
		wantErr  bool
	}{
		{
			name:     "valid rectangle",
			width:    2,
			height:   2,
			wantArea: 4,
		},
		{
			name:    "invalid width",
			width:   -1,
			height:  2,
			wantErr: true,
		},
		{
			name:    "invalid height",
			width:   3,
			height:  -2,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rectangle := Rect{width: test.width, height: test.height}
			got, err := rectangle.area()
			if (err != nil) != test.wantErr {
				t.Fatalf("area() error = %v, wantErr %v", err, test.wantErr)
			}
			if got != test.wantArea {
				t.Errorf("area() = %v, want %v", got, test.wantArea)
			}
		})
	}
}

func TestTriangleArea(t *testing.T) {
	tests := []struct {
		name     string
		base     float64
		height   float64
		wantArea float64
		wantErr  bool
	}{
		{
			name:     "valid triangle",
			base:     4,
			height:   3,
			wantArea: 6,
		},
		{
			name:    "invalid height",
			base:    2,
			height:  -2,
			wantErr: true,
		},

		{
			name:    "invalid base",
			base:    -2,
			height:  2,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			triangle := Tri{base: test.base, height: test.height}
			got, err := triangle.area()
			if (err != nil) != test.wantErr {
				t.Fatalf("area() error = %v, wantErr %v", err, test.wantErr)
			}
			if got != test.wantArea {
				t.Errorf("area() = %v, want %v", got, test.wantArea)
			}
		})
	}
}

func TestTotalArea(t *testing.T) {
	tests := []struct {
		name     string
		shapes   []Shape
		wantArea float64
		wantErr  bool
	}{
		{
			name: "valid total area",
			shapes: []Shape{
				Rect{width: 2, height: 2},
				Tri{base: 2, height: 3},
			},
			wantArea: 7,
		},
		{
			name: "invalid shape",
			shapes: []Shape{
				Rect{width: 2, height: 2},
				Circle{radius: -1},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := totalArea(test.shapes)
			if (err != nil) != test.wantErr {
				t.Fatalf("totalArea() error = %v, wantErr %v", err, test.wantErr)
			}
			if got != test.wantArea {
				t.Errorf("totalArea() = %v, want %v", got, test.wantArea)
			}
		})
	}
}
