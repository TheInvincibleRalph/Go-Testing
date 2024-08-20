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

func assert(t testing.TB, num []int, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("expected %d but got %d, given %v", want, got, num)

	}
}
