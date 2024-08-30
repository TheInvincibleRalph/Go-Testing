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

		{
			"pointersnto things",
			&Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},

		{
			"slices",
			[]Profile{
				{33, "London"},
				{34, "Ralph"},
			},
			[]string{"London", "Ralph"},
		},

		{
			"arrays",
			[2]Profile{
				{33, "London"},
				{34, "Ralph"},
			},
			[]string{"London", "Ralph"},
		},

		{
			"maps",
			map[string]string{
				"Love language": "Go",
				"Love life":     "BackOps",
			},
			[]string{"Go", "BackOps"},
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
	val := getValue(x)

	walkValue := func(value reflect.Value) {
		walk(value.Interface(), fn) // recursively calls walk to handle a struct or slice or array field, also the value.Interface() converts the reflect.Value back to an interface{}
	}

	switch val.Kind() {
	case reflect.String:
		fn(val.String())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			walkValue(val.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			walkValue(val.Field(i))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walkValue(val.MapIndex(key))
		}
	}
}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	return val
}

/*
 The reflect package is used for inspecting the runtime type and value of an object,
 allowing you to interact with types and values dynamically.

 reflect.ValueOf(x) allows you to wrap a value inside a reflect.Value object,
 which provides methods to inspect the underlying value, modify it,
 or even interact with it based on its type.

 reflect.Value is a struct provided by the reflect package that holds a reference
 to the actual value of the object you are reflecting upon.
 It includes information about the value's type,
 kind (whether it is an int, string, struct, etc.), and the value itself.
*/
