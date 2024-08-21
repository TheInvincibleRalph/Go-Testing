package main

import (
	"fmt"
	"math"
	"testing"
)

func RootofEqn(a, b, c float64) (float64, float64) {
	d := b*b - 4*a*c
	root := math.Sqrt(d)
	num1 := (-b + root) / 2 * a
	num2 := (-b - root) / 2 * a

	return num1, num2
}

func TestRootofEqn(t *testing.T) {
	got1, got2 := RootofEqn(1, -2, 1)
	want1, want2 := 1.0, 1.0

	if got1 != want1 && got2 != want2 {
		t.Errorf("expected %.2f and %.2f got %.2f and %.2f", want1, want2, got1, got2)
	}
}
func main() {
	root1, root2 := RootofEqn(1, -2, 1)

	fmt.Printf("x is %f, y is %f", root1, root2)
}
