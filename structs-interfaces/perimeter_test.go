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
		shape Shape
		want  float64
	}{
		{Rectangle{Width: 10.0, Height: 15.0}, 150.0},
		{Circle{Radius: 10}, 314.1592653589793},
		{Triangle{Base: 12, Height: 6}, 36},
	}

	for _, tt := range areaTest {
		got := tt.shape.Area()
		if got != tt.want {
			t.Errorf("expected %g got %g", got, tt.want)
		}
	}

	// checkArea := func(t testing.TB, shape Shape, want float64) {
	// 	t.Helper()

	// 	got := shape.Area()
	// 	if got != want {
	// 		t.Errorf("expected %g got %g", want, got)
	// 	}
	// }
	// t.Run("returns the area of a rectangle", func(t *testing.T) {
	// 	rectangle := Rectangle{10.0, 20.0}
	// 	checkArea(t, rectangle, 200.0)
	// })

	// t.Run("returns the area of a circle", func(t *testing.T) {
	// 	circle := Circle{10}
	// 	checkArea(t, circle, 314.1592653589793)

	// })
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
