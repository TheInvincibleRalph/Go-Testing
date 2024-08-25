package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func Greet(writer io.Writer, name string) { //injects the Writer interface to make our code testable (and predictable, meaning we can decide the behaviour of the function)
	fmt.Fprintf(writer, "Hello, %s", name)
}

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{} //instantiate a buffer
	Greet(&buffer, "Chris")

	got := buffer.String()
	want := "Hello, Chris"

	if got != want {
		t.Errorf("expected %s got %s", want, got)
	}
}

// fmt.Fprintf is like fmt.Printf but instead takes a Writer to send the string to, whereas fmt.Printf defaults to stdout.
// Dependency injection allows to write great general-purpose functions
// In "real life" you would inject in something that writes to stdout. But in test, you write to a anything that implements the Writer interface
// REMEMBER: A Writer is written TO while a Reader is read FROM (WT-RF)
// When you use a Writer, you are passing data to it, and it writes that data to its underlying destination.

/*

	type Writer interface {
    Write(p []byte) (n int, err error)
}
*/
