package main

import (
	"testing"
	"time"
)

type SpyTime struct {
	durationSlept time.Duration
}

// The Sleep method sets the durationSlept field
// to the duration passed as an argument.
// Instead of actually sleeping (like time.Sleep would),
// it just records the duration.
func (s *SpyTime) Sleep(duration time.Duration) {
	s.durationSlept = duration
}

// This is a test function that verifies the behavior of ConfigurableSleeper
func TestConfigurableSleeper(t *testing.T) {
	sleepTime := 5 * time.Second

	spyTime := &SpyTime{}
	sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
	sleeper.Sleep()

	if spyTime.durationSlept != sleepTime {
		t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
	}
}
