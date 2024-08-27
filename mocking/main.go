package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const finalWord = "Go!"
const countdownStart = 3

// Default sleeper
type DefaultSleeper struct{}

type Sleeper interface {
	Sleep() //method Sleep() delays output to stdout by 1s
}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

// Configurable sleeper
type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}

func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}

func Countdown(out io.Writer, sleeper Sleeper) { //injects the Sleeper interface to make our code testable (and predictable, meaning we can decide the behaviour of the function)
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
		sleeper.Sleep()
	}

	fmt.Fprint(out, finalWord)

}

func main() {
	sleeper := &DefaultSleeper{}
	Countdown(os.Stdout, sleeper)
}

// Main for configurable sleeper:

// func main() {
// 	sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
// 	Countdown(os.Stdout, sleeper)
// }

// The hallmark of this code is Separation of Concern (SoC),
// the external Sleep method was separated from the main function so that it can be injected as a dependency.
// this is done for the sake of testability. But how exactly does injecting an interface enhances testability?

/*
Here are a few thoughts for the case of dependency injection

1. Decoupling from ^Concrete Implementations
By using the Sleeper interface instead of directly calling time.Sleep(),
the Countdown function is decoupled from the specific implementation of the sleep behavior.
This allows the function to remain flexible and adaptable to various needs.
In testing, instead of relying on time.Sleep()—which would introduce real time delays—the
interface allows us to inject a different implementation that behaves in a controlled or predictable manner.

2. Controlled Test Environment
With the ability to inject different implementations of Sleeper,
you can create a custom version of Sleeper that simulates the sleep
behavior without actually causing a delay. This makes the tests run much
faster and more reliably. For example, you might use a SpySleeper or a FakeSleeper
in your tests to count the number of times Sleep() was called or to skip actual waiting time.

3. Using dependency injection with interfaces, is a common practice
in writing robust and testable Go code, especially when dealing with
external interactions, ?side effects, or dependencies that may not be
suitable for direct use in a testing environment.

^A concrete implementation refers to a specific instance of a type
 or a struct that provides the actual behavior or functionality defined by methods.
 It is a direct, tangible form of a type that can be instantiated and interacted with, a
 s opposed to an interface, which only defines a set of methods without providing the underlying implementation.

 A concrete implementation is a type that:

- Implements the behavior specified by methods.
- Can be directly instantiated and used.
- Provides the actual code that performs the tasks.

In the code above, DefaultSleeper is a concrete implementation because:

- It is a specific type (DefaultSleeper) that you can create an instance of using &DefaultSleeper{}.
- It provides a concrete method (Sleep()) that performs a specific action (sleeping for one second).

?check documentation for definition


RULE OF THUMB:

Instead of hardcoding dependencies, use interfaces to inject the dependencies.
This allows changing the behavior without modifying the function that relies on the dependency.
*/
