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

---

## What is `interface{}`?

In Go, `interface{}` (pronounced as "empty interface") is a special type that can hold values of any type. It is defined as an interface that has no methods:

### Key Characteristics of `interface{}`

1. **Empty Interface:**
   - An empty interface does not specify any methods. This means that any type satisfies this interface since it doesn't require the type to implement any specific methods.
   - It can be thought of as a "catch-all" type, capable of storing values of any data type.

2. **Dynamic Typing:**
   - While Go is statically typed, `interface{}` allows for a form of dynamic typing. When a value of any type is assigned to an `interface{}` variable, the actual type information is stored alongside the value.
   - This makes `interface{}` a powerful tool for situations where you need to handle values of unknown or varying types.

### Usage Examples

1. **Storing Any Type of Value:**

   You can use `interface{}` to store any type of value:

   ```go
   var x interface{}
   x = 42         // int
   x = "hello"    // string
   x = 3.14       // float64
   x = []int{1, 2, 3} // slice of int
   ```

### Use Cases

- **Generic Data Structures:** `interface{}` is commonly used in generic data structures like slices and maps where the type of elements may vary.
  
- **APIs and Libraries:** Functions that interact with external systems (like databases, network communication, etc.) often use `interface{}` to handle a wide range of input and output types.

- **Decoupling Components:** Using `interface{}` can help in decoupling components, making the code more flexible and easier to test by accepting any type of value.

### Drawbacks

- **Type Safety:** Using `interface{}` means sacrificing some type safety. It's easy to accidentally store or retrieve a value of the wrong type, leading to runtime errors.
  
- **Performance:** Operations involving `interface{}` can be slower due to the need for type assertions and dynamic type handling.

> *In Go `any` is an alias for `interface{}`*

---

## Context

A **context** in Go is a powerful concept used to manage and control the lifecycle of operations, typically within concurrent or distributed systems like HTTP requests, database queries, or goroutines. It helps manage things like **timeouts**, **deadlines**, and **cancellations**. The context allows passing request-scoped values, cancellation signals, and timeouts between function calls.

### Key Uses of Context:
1. **Cancellation**: A context can be cancelled, which allows you to stop or clean up operations early, especially useful when a user aborts an HTTP request.
2. **Timeouts**: You can set timeouts or deadlines for operations to ensure that long-running tasks are stopped after a certain amount of time.
3. **Request-scoped Data**: You can pass additional data, like user authentication details, that can be accessed by downstream functions.

### How it works:
- **Context is passed down a call chain**: Context is usually passed as the first argument in functions that perform operations like handling HTTP requests, making database queries, or starting goroutines.
- **Propagation of cancellation or timeout**: When the parent context is cancelled or reaches its timeout, it propagates this signal to all functions that share this context, allowing them to stop their work.

---

### Types of Contexts in Go:

1. **`context.Background()`**: This is the root context. It's typically used when no higher context is available. It's commonly used to start a top-level request or process.
   
   ```go
   ctx := context.Background()
   ```

2. **`context.TODO()`**: This is a placeholder when you are unsure about which context to use. It's often used during development and testing.

   ```go
   ctx := context.TODO()
   ```

3. **`context.WithCancel(parent)`**: This creates a derived context from the parent, and the `cancel()` function is used to signal cancellation to the child context.

   ```go
   ctx, cancel := context.WithCancel(parentCtx)
   ```

4. **`context.WithTimeout(parent, timeout)`**: This creates a context that automatically cancels after the specified timeout duration. It helps to avoid long-running operations.

   ```go
   ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
   ```

5. **`context.WithDeadline(parent, time)`**: Similar to `WithTimeout`, but you specify an exact point in time for the cancellation to occur.

   ```go
   ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(5*time.Second))
   ```

---

### Example: Cancellation with Context

Imagine an HTTP server handling a request. If the client cancels the request (e.g., closes their browser), the server should stop processing that request. Using a context, you can detect that the request was cancelled and stop the work early.

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // Get context from the request

    select {
    case <-time.After(10 * time.Second):
        fmt.Fprintln(w, "Processed Request")
    case <-ctx.Done(): // Check if context is cancelled
        fmt.Fprintln(w, "Request was cancelled")
    }
}
```

In this example:
- **`ctx.Done()`** listens for a cancellation event, and if the context is cancelled (such as if the client cancels the request), the server responds accordingly.

---

### Why is Context Important?

1. **Graceful Handling of Operations**: It allows for graceful cancellation of tasks, making sure resources are not wasted on abandoned or long-running tasks.
2. **Concurrency**: Context works well in concurrent environments like Go, where multiple goroutines may be running operations, and you need a way to control their lifecycle.
3. **Timeouts and Deadlines**: Setting timeouts or deadlines ensures that operations don’t block indefinitely.

By using context, you can make your code more efficient, responsive, and resilient to errors or unexpected events like user cancellations or timeouts.

The `context.Context` interface looks like this: 
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```