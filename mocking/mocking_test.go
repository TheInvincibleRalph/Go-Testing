package main

import (
	"bytes"
	"testing"
)

type SpySleeper struct {
	Calls int //stores the number of times Sleep() is called
}

func (s *SpySleeper) Sleep() {
	s.Calls++ //increments the number of times Sleep() is called
}

func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}
	SpySleeper := &SpySleeper{}

	Countdown(buffer, SpySleeper)

	got := buffer.String()
	want := `3
	2
	1
	Go!`

	if got != want {
		t.Errorf("expected %q got %q", want, got)
	}

	if SpySleeper.Calls != 3 {
		t.Errorf("not enough calls to sleeper, want 3 got %d", SpySleeper.Calls)
	}
}
