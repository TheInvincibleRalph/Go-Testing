package reflection

import (
	"reflect"
	"testing"
)

// Challenge: Write a function walk(x interface{}, fn func(string))
// which takes a struct x and calls fn for all strings fields found inside. difficulty level: recursively.

func TestWalk(t *testing.T) {

	expected := "Chris"

	var got []string

	x := struct {
		Name string
	}{expected}

	walk(x, func(input string) {
		got = append(got, input)
	})

	if len(got) != 1 {
		t.Errorf("wrong number of function calls, got %d want %d", len(got), 1)
	}

	if got[0] != expected {
		t.Errorf("expected %q got %q", expected, got[0])
	}
}

func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)
	field := val.Field(0)
	fn(field.String())
}
