package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const finalWord = "Go!"
const countdownStart = 3

type DefaultSleeper struct{}

type Sleeper interface {
	Sleep() //method Sleep() delays output to stdout by 1s
}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

func Countdown(out io.Writer, sleeper Sleeper) { //injects the Sleeper interface to make our code testable (and predictable, meaning we can decide the behaviour of the function)
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintf(out, "%d\n\t", i)
		sleeper.Sleep()
	}
	fmt.Fprint(out, finalWord)

}

func main() {
	sleeper := &DefaultSleeper{}
	Countdown(os.Stdout, sleeper)
}
