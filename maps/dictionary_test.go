package main

import (
	"errors"
	"testing"
)

type Dictionary map[string]string

var ErrNotFound = errors.New("could not find the word you are looking for")

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"} //instantiating a map

	t.Run("known word", func(t *testing.T) {
		got, _ := dictionary.Search("test") //calling Search on the Dictionary instance(or type)
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, got := dictionary.Search("unknown")
		if got == nil {
			t.Fatal("expected to get an error.")
		}

		assertError(t, got, ErrNotFound) //the .Error() method gets the string in the err message
	})
}

func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word] // The second value is a boolean which indicates if the key was found successfully.
	if !ok {
		return "", ErrNotFound
	}
	return definition, nil
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("expected %q got %q", want, got)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}
