package main

import (
	"math"
	"testing"
)

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

type Triangle struct {
	Base   float64
	Height float64
}

type Shape interface {
	Area() float64
}

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 15.0}
	got := Perimeter(rectangle)
	want := 50.0

	if got != want {
		t.Errorf("expected %.2f got %.2f", want, got)
	}
}

func TestArea(t *testing.T) {

	areaTest := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{Width: 10.0, Height: 15.0}, hasArea: 150.0},
		{name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36},
	}

	for _, tt := range areaTest {
		// using tt.name from the test cases as the t.Run test name
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v expected %g got %g", tt.shape, got, tt.hasArea)
			}
		})

	}
}

func Perimeter(rectangle Rectangle) float64 {
	return 2 * (rectangle.Width + rectangle.Height)
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (c Circle) Area() float64 {
	return math.Pi * (c.Radius * c.Radius)
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

// Methods are very similar to functions but they are called by invoking them on an instance of a particular type.
// Where you can just call functions wherever you like, such as Area(rectangle) you can only call methods on "things".

// Interfaces are a very powerful concept in statically typed languages like Go because they allow you to make
// functions that can be used with different types and create highly-decoupled code whilst still maintaining type-safety.

// If a type does not implement an interface, it means that type is not associated with the methods defined within the interface.
// Interfaces can be used or implemented by different types, as expressed in this package.
