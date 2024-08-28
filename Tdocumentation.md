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

---

## Testing HTTP Handlers

> Testing code that uses HTTP is so common that Go has tools in the standard library to help you test it. In the standard library, there is a package called net/http/httptest which enables users to easily create a mock HTTP server  ~Chris James


- `http.HandlerFunc` is a type that looks like this: `type HandlerFunc func(ResponseWriter, *Request)`. *All it's really saying is it needs a function that takes a ResponseWriter and a Request, which is not too surprising for an HTTP server.*

- `httptest.NewServer` takes an `http.HandlerFunc`

- `w.WriteHeader(http.StatusOK)` writes an `OK` response back to the caller

---

## Empty Struct (`struct{}`) and Memory Allocation

1. **What is `struct{}`?**
   - In Go, `struct{}` represents an **empty struct** type. It has no fields, which means it occupies zero bytes of memory. 
   - It’s the smallest possible data type in Go. Unlike other types, such as `bool`, `int`, or `string`, which require some memory allocation, `struct{}` is truly zero-sized.

2. **Memory Efficiency:**
   - **Zero Allocation:** The primary reason for using `struct{}` over other types is its zero allocation characteristic. It doesn't take up any space in memory. When you use a `chan struct{}`, sending and receiving operations don’t actually transfer any data, which means no memory needs to be allocated or freed.
   - **Minimal Overhead:** When working with concurrency, minimizing overhead is crucial. Since `struct{}` has no fields, there’s no unnecessary allocation or copying involved. This makes it very efficient for signaling or synchronization purposes.

3. **Signaling with Channels:**
   - In Go, channels are often used for signaling purposes, such as to notify when a goroutine should stop or when an event has occurred. In these cases, the data being sent is often irrelevant; only the occurrence of the event matters.
   - Using `chan struct{}` in such cases indicates that the channel is used purely for signaling. The only concern is whether a signal has been received, not the content of the signal. This is a clear, idiomatic way to express that no data is needed.

4. **Comparison with Other Types:**
   - **`chan bool`:** A boolean channel could be used to send a true or false signal. However, a boolean value still requires memory allocation (typically 1 byte), and using `bool` implies that the true or false value is significant, which may not be the case for signaling.
   - **`chan int`:** Similarly, using an integer channel requires memory allocation and might imply that the actual integer value has some meaning, which is not the intention in signaling cases.
   - **`chan struct{}`:** By contrast, `chan struct{}` conveys that the signal itself is what matters, not any particular value. It makes it explicit to anyone reading the code that the presence of a value in the channel is the signal, not the value itself.

5. **Use Cases of `chan struct{}`:**
   - **Signaling Completion or Stop:** It’s common to use `chan struct{}` to signal the completion of a task or to stop a goroutine. The receiver simply waits for a signal, which indicates it should proceed or terminate.
   - **Mutex Implementation:** `struct{}` channels are sometimes used in custom mutex implementations or to control access to a shared resource.
   - **Broadcasting Events:** In some designs, `struct{}` channels can be used to broadcast events to multiple listeners. Since the event data itself isn’t important, only the occurrence of the event, `struct{}` is a natural fit.

### Conclusion

Using `struct{}` in channels is a common idiom in Go for signaling without any data transfer. It leverages the zero-memory footprint of the empty struct type to efficiently manage synchronization and communication in concurrent programs. This approach helps to express the intent of signaling without transferring actual data, thereby adhering to best practices in Go programming.

---

## `make`-ing channels

**Always `make` channels**

 Always use the `make` function when creating a channel; rather than say `var ch chan struct{}`. 
 
 When you use `var` the variable will be initialised with the "zero" value of the type. So for `string` it is `""`, `int` it is 0, etc.

For channels the zero value is `nil` and if you try and send to it with `<-` it will block forever because you cannot send to `nil` channels

---

## The `select` statement

`select` allows you to wait on multiple channels. The first one to send a value "wins" and the code underneath the `case` is executed.