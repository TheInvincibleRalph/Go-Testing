package main

import (
	"bytes"
	"reflect"
	"testing"
)

type SpyCountDownOperations struct { //stores how many times each spies are called
	Calls []string
}

func (s *SpyCountDownOperations) Sleep() { //spies on the sleep operation of Countdown
	s.Calls = append(s.Calls, sleep)
}

func (s *SpyCountDownOperations) Write(p []byte) (n int, err error) { //spies on the main countdown logic
	s.Calls = append(s.Calls, write)
	return
}

const write = "write"
const sleep = "sleep"

func TestCountdown(t *testing.T) {

	t.Run("prints 3 to Go!", func(t *testing.T) {

		buffer := &bytes.Buffer{}

		Countdown(buffer, &SpyCountDownOperations{})

		got := buffer.String()
		want := `3
2
1
Go!`

		if got != want {
			t.Errorf("expected %q got %q", want, got)
		}

	})

	t.Run("sleep before every print", func(t *testing.T) {
		spySleepPrinter := &SpyCountDownOperations{}
		Countdown(spySleepPrinter, spySleepPrinter)

		want := []string{
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
			t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
		}
	})
}
