package reflection

import (
	"reflect"
	"testing"
)

// Challenge: Write a function walk(x interface{}, fn func(string))
// which takes a struct x and calls fn for all strings fields found inside. difficulty level: recursively.

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func TestWalk(t *testing.T) {

	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"struct with one string field",
			struct {
				Name string
			}{"Chris"},
			[]string{"Chris"},
		},

		{
			"struct with two string field",
			struct {
				Name string
				City string
			}{"Chris", "London"},
			[]string{"Chris", "London"},
		},

		{
			"struct with non string field",
			struct {
				Name string
				Age  int
			}{"Chris", 33},
			[]string{"Chris"},
		},

		{
			"nested fields",
			Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("expected %v got %v", test.ExpectedCalls, got)
			}
		})
	}
}

// recursively traverses the fields of a struct,
// if it encounters a string field, it calls the provided function fn with the string value
func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.String {
			fn(field.String())
		}

		if field.Kind() == reflect.Struct {
			walk(field.Interface(), fn) //recursively calls walk to handle a struct field
		}
	}
}

/*
 The reflect package is used for inspecting the runtime type and value of an object,
 allowing you to interact with types and values dynamically.
*/
