package array

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("collection of 10 numbers", func(t *testing.T) {
		num := []int{3, 6, 9, 1, 4, 5, 9, 8, 0, 7}
		got := Sum(num)
		want := 52
		assert(t, num, got, want)
	})

	t.Run("collection of 5 numbers", func(t *testing.T) {
		num := []int{3, 6, 9, 1, 4}
		got := Sum(num)
		want := 23
		assert(t, num, got, want)
	})

}

func TestAll(t *testing.T) {

	got := SumAll([]int{1, 2, 3, 4, 5}, []int{6, 7, 8, 9, 10})
	want := []int{15, 40}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v got %v", want, got)
	}
}

func Sum(num []int) int {
	result := 0
	for _, numbers := range num {
		result += numbers
	}
	return result
}

func SumAll(numsToSum ...[]int) []int {
	var sums []int

	for _, numbers := range numsToSum {
		sums = append(sums, Sum(numbers))
	}
	return sums

}

func SumAllTails(numsToSum ...[]int) []int {
	var sums []int

	for _, numbers := range numsToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := Sum(numbers) - numbers[0] //or tail := numbers[1:]
			sums = append(sums, tail)         //then sums := append(sums, Sum(tail))
		}
	}
	return sums
}

func TestSumAllTail(t *testing.T) {

	assertEqual := func(t testing.TB, got, want []int) { //this is how to assign a function to a variable!
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v got %v", want, got)
		}
	}

	t.Run("sums the tails of a slice", func(t *testing.T) {
		got := SumAllTails([]int{1, 2, 3, 4, 5}, []int{6, 7, 8, 9, 10})
		want := []int{14, 34}

		assertEqual(t, got, want)
	})

	t.Run("safely sums empty slice", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{1, 2, 3, 4})
		want := []int{0, 9}

		assertEqual(t, got, want)
	})

}

func assert(t testing.TB, num []int, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("expected %d but got %d, given %v", want, got, num)

	}
}
