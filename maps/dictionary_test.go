package main

import "testing"

func TestSearch(t *testing.T) {
	dictionary := map[string]string{"test": "this is just a test"}

	got := Search(dictionary, test)
	want := "this is just a test"

	if got =! want {
		t.Errorf("expected %q got %q, given %q", want, got, "test")
	}
}
