# Documentation starting from Tutorial 11: Concurrency

## On Goroutines and Anonymous functions

> Normally in Go when we call a function `doSomething()` we wait for it to return (even if it has no value to return, we still wait for it to finish). We say that this operation is *blocking* - it makes us wait for it to finish. An operation that does not block in Go will run in a separate process called a *goroutine.*

To tell Go to start a new goroutine we turn a function call into a `go` statement by putting the keyword `go` in front of it: `go doSomething()`.

**Anonymous functions in Go**

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func() {
			results[url] = wc(url)
		}()
	}

	return results
}
```

> Anonymous functions have a number of features which make them useful, two of which we're using above. Firstly, they can be executed at the same time that they're declared - this is what the `()` at the end of the anonymous function is doing. Secondly they maintain access to the lexical scope in which they are defined - all the variables that are available at the point when you declare the anonymous function are also available in the body of the function.

---

## Closures and Goroutines
> In programming, a closure is a function that "captures" the environment in which it is defined. This means that the function retains access to the variables that were in scope when the function was created, even after those variables would normally go out of scope. Closures allow functions to access non-local variables even after those variables have gone out of the context in which they were created.

```go
func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func() {
			results[url] = wc(url)
		}()

	}

	time.Sleep(2 * time.Second)
	return results
}
```

Here, the anonymous function is a closure that captures the url variable. However, because url is shared across iterations, all the closures capture the same url variable reference. When these closures are eventually executed (possibly after the loop has completed), they all use the latest value of url. This leads to unexpected results, like every goroutine using the last URL.

### Understanding the *Shared-URL* Problem

In the `CheckWebsites` function, you're spawning a new goroutine for each URL in the `urls` slice. The goroutine function is supposed to use the current value of `url` to check the website and store the result in the `results` map. However, all the goroutines end up using the same `url` variable, which results in them all referencing the same memory location.

Here's what's happening step-by-step:

1. **Variable Reuse**: The `url` variable in the `for _, url := range urls` loop is being reused for each iteration. It's a single variable, and in each iteration, it's assigned a new value from the `urls` slice.

2. **Closure Capturing**: When you create a goroutine with `go func() { ... }()`, the function literal captures the `url` variable by reference, not by value. This means that all goroutines share the same `url` variable. They don't create their own independent copies of `url` at the time they are created.

3. **Race Condition**: By the time the goroutine gets executed, the `url` variable might have already been updated to the next value in the loop, or even completed the loop, resulting in all goroutines referring to the last value of `url`. This is why you might end up with only one result, or multiple goroutines accessing the last `url` in the slice, depending on timing and concurrency.

### Why Does This Happen?

The key point is that the goroutine function is executed asynchronously, which means it might not run immediately when the `go` statement is called. By the time the goroutine actually runs, the loop may have progressed, and `url` may have a different value or may be at the last value.

### Solution

To fix this issue, you need to pass the `url` variable as a parameter to the goroutine. This way, each goroutine will receive its own copy of `url`, preserving the value at the time the goroutine was created. Here's how you can modify the function:

```go
func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func(u string) {
			results[u] = wc(u)
		}(url) // Pass the current value of url as an argument
	}

	time.Sleep(2 * time.Second)
	return results
}
```

### Explanation of the Fix

- **Closure with Parameter**: By defining the anonymous function with a parameter `(u string)` and passing `url` as an argument, each invocation of the goroutine captures its own copy of `url` as `u`. This ensures that the value of `url` at the time the goroutine was created is used within that goroutine.
  
- **Independent Copies**: Now, each goroutine has its own `u` variable, which is a copy of the `url` at that specific iteration. This prevents the race condition where all goroutines end up using the same `url` value.