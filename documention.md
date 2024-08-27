# Testing Terminologies

## TEST DOUBLES (Stub, fake, spies, dummy, mock)

## Stub

A **stub** is a type of test double that provides predefined responses to method calls. It is commonly used when you want to test the behavior of a function that relies on an external dependency but do not care about the specifics of that dependency's behavior. Instead, you focus on the response needed to drive the test.

### Example Scenario

Let's say we have an application that checks if a user exists in a database before sending a welcome email. We want to test the function that checks the user's existence, but without actually connecting to a real database. Instead, we can use a stub that simulates the behavior of the database.

### Go Code Example Using a Stub

1. **Define the Interface**

   First, define an interface that represents the dependency, such as a database:

   ```go
   type UserRepository interface {
       UserExists(username string) bool
   }
   ```

   Here, `UserRepository` has a single method `UserExists` which checks if a user exists in the database.

2. **Define the Function to be Tested**

   Now, define a function that uses the `UserRepository` interface:

   ```go
   func SendWelcomeEmail(repo UserRepository, username string) string {
       if repo.UserExists(username) {
           return "User already exists. No email sent."
       }
       return "Welcome email sent to " + username
   }
   ```

   This function checks if the user exists. If the user does exist, it returns a message indicating no email was sent. Otherwise, it sends a welcome email.

3. **Create a Stub**

   Next, create a stub that implements the `UserRepository` interface. The stub will provide a predefined response to the `UserExists` method:

   ```go
   type UserRepositoryStub struct {
       ExistingUsers map[string]bool
   }

   func (stub *UserRepositoryStub) UserExists(username string) bool {
       // Simulate checking if a user exists by returning a predefined response
       exists, ok := stub.ExistingUsers[username]
       return ok && exists
   }
   ```

   The `UserRepositoryStub` struct contains a map of usernames to simulate the existing users. The `UserExists` method checks this map to determine if the user exists.

4. **Write the Test Case**

   Now, use the stub in a test case to verify the behavior of `SendWelcomeEmail`:

   ```go
   package main

   import "testing"

   func TestSendWelcomeEmail(t *testing.T) {
       // Create a stub with predefined existing users
       stub := &UserRepositoryStub{
           ExistingUsers: map[string]bool{
               "existingUser": true,
           },
       }

       // Test case where the user already exists
       result := SendWelcomeEmail(stub, "existingUser")
       expected := "User already exists. No email sent."
       if result != expected {
           t.Errorf("expected %s but got %s", expected, result)
       }

       // Test case where the user does not exist
       result = SendWelcomeEmail(stub, "newUser")
       expected = "Welcome email sent to newUser"
       if result != expected {
           t.Errorf("expected %s but got %s", expected, result)
       }
   }
   ```

### Explanation

1. **Interface and Function**: We define a `UserRepository` interface and a `SendWelcomeEmail` function that relies on this interface. The function checks if a user exists using the interface.

2. **Stub Implementation**: The `UserRepositoryStub` struct implements the `UserRepository` interface. It contains a map to simulate existing users. The `UserExists` method uses this map to provide a predefined response.

3. **Test Case**: The test creates an instance of `UserRepositoryStub` with some predefined users. It then calls `SendWelcomeEmail` to verify that the function behaves correctly depending on whether the user exists.

### Benefits of Using a Stub

- **Isolation**: The `SendWelcomeEmail` function can be tested independently of the real database.
- **Control**: You can control the responses from `UserExists` to test various scenarios.
- **Speed**: Tests run faster because they do not rely on actual database calls.
- **Predictability**: By using predefined responses, the test outcomes are predictable and reliable.

Style
Mocks vs Stubs = Behavioral testing vs State testing

Principle
According to the principle of Test only one thing per test, there may be several stubs in one test, but generally there is only one mock.

---

## Fakes 

Fakes are a powerful type of test double that simulate the behavior of real components by maintaining some form of internal state. They provide a more realistic environment for testing than stubs because they can handle multiple interactions and keep track of state changes, similar to what happens in a production environment.

### Detailed Example of Using Fakes

Let's consider an example of an online bookstore where we need to interact with a database to manage books. We'll create a fake that simulates a real database by keeping a list of books in memory. This fake will allow us to perform operations like adding, retrieving, and deleting books.

#### 1. Define the Interface

We'll start by defining an interface that represents our book repository. This interface will declare the methods that our fake (and real) repository needs to implement.

```go
package main

type Book struct {
    ID     string
    Title  string
    Author string
}

type BookRepository interface {
    AddBook(book Book) error
    GetBook(id string) (*Book, error)
    DeleteBook(id string) error
}
```

- The `BookRepository` interface declares three methods: `AddBook`, `GetBook`, and `DeleteBook`.
- The `Book` struct represents a book entity with an ID, title, and author.

#### 2. Implement the Fake

Now, we will implement a fake that adheres to the `BookRepository` interface. This fake will store books in memory using a map, simulating a database.

```go
package main

import (
    "errors"
    "fmt"
)

// BookRepositoryFake is a fake implementation of the BookRepository interface.
type BookRepositoryFake struct {
    books map[string]Book // In-memory store to hold books
}

// NewBookRepositoryFake initializes a new fake book repository.
func NewBookRepositoryFake() *BookRepositoryFake {
    return &BookRepositoryFake{
        books: make(map[string]Book),
    }
}

// AddBook adds a book to the fake repository.
func (repo *BookRepositoryFake) AddBook(book Book) error {
    if _, exists := repo.books[book.ID]; exists {
        return errors.New("book already exists")
    }
    repo.books[book.ID] = book
    return nil
}

// GetBook retrieves a book from the fake repository by ID.
func (repo *BookRepositoryFake) GetBook(id string) (*Book, error) {
    book, exists := repo.books[id]
    if !exists {
        return nil, errors.New("book not found")
    }
    return &book, nil
}

// DeleteBook removes a book from the fake repository by ID.
func (repo *BookRepositoryFake) DeleteBook(id string) error {
    if _, exists := repo.books[id]; !exists {
        return errors.New("book not found")
    }
    delete(repo.books, id)
    return nil
}
```

- `BookRepositoryFake` is a struct that implements the `BookRepository` interface.
- It uses a map (`books map[string]Book`) to store books in memory, simulating a database.
- Methods like `AddBook`, `GetBook`, and `DeleteBook` manipulate this in-memory map, providing a realistic way to interact with a "database."

#### 3. Example Usage of the Fake

Let's create some functions to interact with the `BookRepositoryFake` and write a test to demonstrate how this fake maintains state.

```go
package main

import (
    "fmt"
)

func main() {
    // Create a new fake repository
    repo := NewBookRepositoryFake()

    // Add a book
    book1 := Book{ID: "1", Title: "1984", Author: "George Orwell"}
    err := repo.AddBook(book1)
    if err != nil {
        fmt.Println("Error:", err)
    }

    // Try to add the same book again
    err = repo.AddBook(book1)
    if err != nil {
        fmt.Println("Error:", err) // Expected: "book already exists"
    }

    // Retrieve the book
    retrievedBook, err := repo.GetBook("1")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Retrieved Book: %+v\n", *retrievedBook)
    }

    // Delete the book
    err = repo.DeleteBook("1")
    if err != nil {
        fmt.Println("Error:", err)
    }

    // Try to retrieve the deleted book
    retrievedBook, err = repo.GetBook("1")
    if err != nil {
        fmt.Println("Error:", err) // Expected: "book not found"
    }
}
```

#### Explanation

1. **Initialization**: `NewBookRepositoryFake()` initializes a new instance of the fake book repository. It creates an empty map to store books.

2. **Adding a Book**: The `AddBook` method adds a book to the map. If the book already exists (i.e., the ID is already a key in the map), it returns an error.

3. **Retrieving a Book**: The `GetBook` method looks up a book by its ID in the map. If the book is found, it returns a pointer to the `Book` struct; otherwise, it returns an error.

4. **Deleting a Book**: The `DeleteBook` method deletes a book from the map using its ID. If the book is not found, it returns an error.

#### Why Fakes are Useful

- **State Management**: Fakes maintain state. In this example, the `BookRepositoryFake` keeps track of books that are added, retrieved, and deleted. This allows you to test interactions that depend on previous operations (e.g., checking that a book can be retrieved only if it was added first).

- **Realistic Simulation**: Fakes provide a more realistic simulation of how your application interacts with external systems. This is especially useful in integration tests where the interaction with components (like databases) needs to be close to reality.

- **Isolated Testing**: You can test business logic without depending on a real database connection. This makes tests faster and more reliable because they are not subject to the availability or state of external systems.


### Fakes when used in testing

```go
// BookRepositoryFake is a fake implementation of the BookRepository interface.
type BookRepositoryFake struct {
    books map[string]Book // In-memory store to hold books
}

// NewBookRepositoryFake initializes a new fake book repository.
func NewBookRepositoryFake() *BookRepositoryFake {
    return &BookRepositoryFake{
        books: make(map[string]Book),
    }
}
```

### Scenario

Suppose you have a function in your application that checks if a book exists in the inventory and returns its details if it does. You want to write tests for this function using the `BookRepositoryFake` to simulate the database.

### Test function

Here’s a function that checks if a book exists in the repository and retrieves it:

```go
package main

import (
    "fmt"
)

// GetBookDetails checks if a book exists by its ID and returns the details.
func GetBookDetails(repo BookRepository, id string) (string, error) {
    book, err := repo.GetBook(id)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("Title: %s, Author: %s", book.Title, book.Author), nil
}
```

### Writing Tests Using Fakes

We’ll use Go’s testing package to write unit tests for the `GetBookDetails` function. The `BookRepositoryFake` will act as a stand-in for the real database.

1. **Test Setup**: Initialize a new `BookRepositoryFake` and populate it with some test data.
2. **Test Execution**: Call `GetBookDetails` with different scenarios (book exists, book does not exist).
3. **Assertions**: Check the results to see if they match the expected outcomes.

### Test Code

```go
package main

import (
    "testing"
)

func TestGetBookDetails(t *testing.T) {
    // Initialize the fake repository
    repo := NewBookRepositoryFake()

    // Add a book to the fake repository
    book := Book{ID: "1", Title: "1984", Author: "George Orwell"}
    err := repo.AddBook(book)
    if err != nil {
        t.Fatalf("failed to add book: %v", err)
    }

    // Test: Retrieve existing book details
    result, err := GetBookDetails(repo, "1")
    if err != nil {
        t.Errorf("expected no error, but got %v", err)
    }
    expected := "Title: 1984, Author: George Orwell"
    if result != expected {
        t.Errorf("expected '%s', but got '%s'", expected, result)
    }

    // Test: Retrieve non-existing book details
    _, err = GetBookDetails(repo, "2")
    if err == nil {
        t.Errorf("expected error for non-existing book, but got none")
    }
}
```

### Explanation of the Test Code

1. **Test Setup**:
   - `repo := NewBookRepositoryFake()`: Creates a new instance of `BookRepositoryFake`. This is an in-memory store that will hold our test data.
   - `repo.AddBook(book)`: Adds a book to the fake repository. This simulates having a book already stored in a real database.

2. **Test Execution and Assertions**:
   - **Test Case 1 (Book Exists)**:
     - We call `GetBookDetails(repo, "1")` to retrieve the details of the book with ID `"1"`.
     - We check if the returned details match the expected string.
     - We also check that no error is returned.
   
   - **Test Case 2 (Book Does Not Exist)**:
     - We call `GetBookDetails(repo, "2")` to try to retrieve a book with ID `"2"`, which was not added to the fake.
     - We check that an error is returned, indicating the book was not found.

---

## Spies

**Spies** are a type of test double used in unit testing to verify interactions between components. They are used to check whether certain methods or functions have been called, with what arguments, how many times, and in what order. In Go, spies are typically implemented using custom types or functions that conform to specific interfaces and record relevant data for assertions.

### How Spies Work in Go

1. **Recording Behavior**: Spies keep track of method calls, including the number of times a method was called, the arguments it was called with, and possibly other information like return values.

2. **Verifying Interactions**: After running the test, spies allow you to assert whether the expected interactions occurred. For example, you can verify if a function was called the correct number of times or with the correct arguments.

3. **Non-Intrusive**: Spies in Go often do not change the original behavior of the code under test. Instead, they capture interactions, making them less intrusive than mocks, which may require predefined behavior.

### Example of Using Spies in Go

Let's consider a scenario where we have a `Notifier` interface, and we want to test a function that depends on this interface. We'll create a spy that implements the `Notifier` interface to verify that certain methods are called as expected.

#### Step 1: Define the Interface

```go
type Notifier interface {
    Notify(message string)
}
```

#### Step 2: Implement the Spy

```go
type SpyNotifier struct {
    Calls []string // To store the messages passed to Notify
}

func (s *SpyNotifier) Notify(message string) {
    s.Calls = append(s.Calls, message) // Record the message
}
```

#### Step 3: Use the Spy in a Test

```go
func TestNotifyUser(t *testing.T) {
    spy := &SpyNotifier{}
    NotifyUser(spy, "Hello, World!")

    if len(spy.Calls) != 1 {
        t.Errorf("expected 1 call to Notify, but got %d", len(spy.Calls))
    }

    if spy.Calls[0] != "Hello, World!" {
        t.Errorf("expected 'Hello, World!' but got '%s'", spy.Calls[0])
    }
}
```

In this example:

- `SpyNotifier` implements the `Notifier` interface and records any calls to `Notify`.
- The test function `TestNotifyUser` verifies that the `Notify` method is called exactly once and with the correct message.

### Significance of Spies in Go Testing

1. **Verifying Side Effects**: Spies are useful for testing the side effects of a function. For instance, you can verify that a logging function was called correctly without needing to inspect the log files.

2. **Testing Behavior, Not State**: While traditional tests often check the state of an object after a function runs, spies allow you to verify that the correct behavior occurred. This is useful when testing functions that interact with other systems or components.

3. **Ensuring Correct Function Interactions**: Spies can ensure that functions interact with their dependencies as expected. For example, if a function is supposed to send a notification under certain conditions, a spy can confirm that the notification method is called with the right arguments.

4. **Isolation**: By using spies, you can isolate the function under test from its dependencies. This means that the test only fails if the function itself is incorrect, not because of the behavior of its dependencies.

5. **Improved Test Coverage**: Spies allow for more detailed testing of interactions, which can lead to better coverage of edge cases and a deeper understanding of how functions behave under various conditions.


### 1. **State**
- **State** refers to the data or the internal properties that an object (or struct) holds. In the context of a struct, the state is represented by the fields or attributes of the struct.

- **Empty Struct and State**: An empty struct (`struct{}`) has no fields, and thus it doesn't hold any state. This means it cannot store any data internally, and its instances are purely identifiers with no data associated with them.

  ```go
  type EmptyStruct struct{}
  ```

  In this case, `EmptyStruct` has no state because it has no fields.

### 2. **Behavior**
- **Behavior** refers to the actions or methods that an object can perform. In Go, behavior is implemented via methods that are associated with a struct type.

- **Empty Struct and Behavior**: Even though an empty struct has no state, *you can still define methods on it, which gives it behavior*. These methods can perform operations, often using external inputs or interacting with other parts of the program, without relying on internal state.

  ```go
  type Logger struct{}

  func (l Logger) LogMessage(message string) {
      fmt.Println("Log:", message)
  }
  ```

  Here, `Logger` is an empty struct, but it has a behavior defined by the `LogMessage` method. Even though `Logger` holds no internal state, it can still perform actions (like logging a message).

### Combining State and Behavior
- **Structs with State**: In many cases, structs are used to bundle both state and behavior. Fields represent the state, and methods define the behavior that operates on this state.

  ```go
  type User struct {
      Name  string
      Email string
  }

  func (u *User) UpdateEmail(newEmail string) {
      u.Email = newEmail
  }
  ```

  In this example, `User` has state (`Name` and `Email`) and behavior (`UpdateEmail` method) that modifies this state.

- **Empty Structs**: When using empty structs, you separate state and behavior. The behavior is purely in the methods, with no internal state to modify or rely on. This can be useful for stateless operations, utility functions, or signaling mechanisms.

### Summary of Behavior and State in Empty Structs
- **State**: An empty struct has no state because it has no fields to hold data.
- **Behavior**: Despite having no state, an empty struct can still have behavior through methods. These methods do not rely on internal data but can still perform meaningful actions.

This separation can be beneficial in certain design patterns where you need *an identifiable type with associated behavior* but no internal data storage, such as in the case of utility objects, service-like functions, or simple event signaling.

### Service-like Functions

**Service-like functions** refer to functions or methods in programming that perform specific operations or tasks, often resembling services in the way they encapsulate and provide functionality without necessarily relying on internal state. These functions are typically associated with utility, processing, or performing actions that don't depend on the data within an instance but instead provide a service or operation that can be reused.

### Characteristics of Service-like Functions:
1. **Stateless**: They often don't depend on or modify internal state. Instead, they operate on parameters passed to them.
2. **Reusable**: They provide reusable functionality that can be called from various parts of a program.
3. **Encapsulated Logic**: They encapsulate specific logic or processes, making the code more modular and maintainable.
4. **Utility-based**: They often perform common tasks, such as formatting, logging, processing data, or interacting with external systems like databases or APIs.

### Examples of Service-like Functions

#### 1. **Logging Service**
A common example is a logging function that logs messages to a console or file. This function doesn't need to know about the internal state of any object; it just needs the message to log.

```go
type Logger struct{}

// Service-like function for logging
func (l Logger) Log(message string) {
    fmt.Println("Log:", message)
}

func main() {
    logger := Logger{}
    logger.Log("This is a log message.")
}
```
Here, `Log` is a service-like function. It provides logging services without needing internal state.

#### 2. **Validation Service**
You can have functions that validate inputs, ensuring data conforms to certain rules before being processed.

```go
type Validator struct{}

// Service-like function for email validation
func (v Validator) ValidateEmail(email string) bool {
    // Simple email validation logic
    if email == "" {
        return false
    }
    return strings.Contains(email, "@")
}

func main() {
    validator := Validator{}
    isValid := validator.ValidateEmail("user@example.com")
    fmt.Println("Is valid email:", isValid)  // Output: Is valid email: true
}
```
Here, `ValidateEmail` is a service-like function that checks if an email string is valid.

#### 3. **Formatting Service**
A function that formats data into a specific structure or string is also a good example.

```go
type Formatter struct{}

// Service-like function for formatting a user's name
func (f Formatter) FormatName(firstName, lastName string) string {
    return fmt.Sprintf("%s %s", firstName, lastName)
}

func main() {
    formatter := Formatter{}
    fullName := formatter.FormatName("John", "Doe")
    fmt.Println("Full Name:", fullName)  // Output: Full Name: John Doe
}
```
Here, `FormatName` is a service-like function that formats and returns a full name string.

#### 4. **Math Service**
Functions that perform mathematical calculations or operations can also be service-like.

```go
type Calculator struct{}

// Service-like function for addition
func (c Calculator) Add(a, b int) int {
    return a + b
}

func main() {
    calculator := Calculator{}
    result := calculator.Add(5, 3)
    fmt.Println("Result of addition:", result)  // Output: Result of addition: 8
}
```
Here, `Add` is a service-like function that performs addition.

#### 5. **HTTP Request Service**
A function that makes HTTP requests can be seen as a service-like function. It performs the task of making a network call without holding any state.

 ```go
type HttpRequester struct{}

// Service-like function for making a GET request
func (hr HttpRequester) Get(url string) (*http.Response, error) {
    return http.Get(url)
}

func main() {
    requester := HttpRequester{}
    response, err := requester.Get("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Response Status:", response.Status)
}
```
Here, `Get` is a service-like function that sends an HTTP GET request.


### Use case of a Stub

You want to test a function that calculates the total price of items in a shopping cart. The function depends on an external service that provides the price of an item. You use a stub to provide a controlled response for this external service.

```go
package main

import (
	"testing"
)

// ItemPriceService defines the interface for fetching item prices.
type ItemPriceService interface {
	GetPrice(itemID string) float64
}

// ShoppingCart calculates the total price of items.
func ShoppingCart(service ItemPriceService, items []string) float64 {
	total := 0.0
	for _, item := range items {
		total += service.GetPrice(item)
	}
	return total
}

// Stub for ItemPriceService.
type StubItemPriceService struct{}

func (s *StubItemPriceService) GetPrice(itemID string) float64 {
	prices := map[string]float64{
		"item1": 10.0,
		"item2": 20.0,
	}
	return prices[itemID]
}

func TestShoppingCart(t *testing.T) {
	service := &StubItemPriceService{}
	items := []string{"item1", "item2"}
	total := ShoppingCart(service, items)
	expected := 30.0
	if total != expected {
		t.Errorf("expected %f, got %f", expected, total)
	}
}
```

> **Extras**
> *The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers. If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.*

---

## Interface Implementation

> An interface in Go is like a contract. It defines a set of methods that a type must implement so that any type that implements all the methods specified in the interface is said to fulfill the contract, and is considered to have implemented the interface. 

**The Contract**

```go
// UserRepository defines the contract for interacting with user data
type UserRepository interface {
    SaveUser(user *User) error
    GetUserByID(id int) (*User, error)
}
```

**Implementation:**

```go
// SQLUserRepository is a concrete implementation of UserRepository for a SQL database
type SQLUserRepository struct {
    db *sql.DB
}

func (repo *SQLUserRepository) SaveUser(user *User) error {
    // SQL logic to save user
    _, err := repo.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
    return err
}

func (repo *SQLUserRepository) GetUserByID(id int) (*User, error) {
    // SQL logic to retrieve a user by ID
    row := repo.db.QueryRow("SELECT id, username, password FROM users WHERE id = ?", id)
    var user User
    err := row.Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```
---

## Domain-Driven Architecture (Domain-from-Side-effect theory)

### **Domain Code**:
- **Domain** refers to the core logic and rules that govern the business or problem your application is solving. It’s the heart of your application, where you define how things should work according to the business rules, independent of any external systems or frameworks.
- **Domain code** should be pure, meaning it doesn’t interact directly with external systems (like databases, APIs, file systems, etc.). Instead, it focuses on implementing the core business logic.

### **The Outside World (Side-Effects)**:
- The **outside world** includes anything external to your core business logic, like databases, web services, user interfaces, and file systems.
- **Side-effects** are interactions with the outside world. For example, reading or writing to a database, sending HTTP requests, or logging to a file are all side-effects because they involve interacting with systems outside of your core logic.

### **Separation of Concerns**:
- **Separating your domain code from the outside world** means keeping your core business logic isolated from code that deals with side-effects. This separation allows your domain code to remain clean, testable, and independent of external systems.
- You can achieve this by using design patterns like **dependency injection**, **interfaces**, and **repositories**. These patterns allow your domain code to depend on abstractions (like interfaces) rather than concrete implementations (like a specific database).

---

## Type-safe

When a function is not type-safe, it means that it can accept and operate on values of different data types without strict enforcement by the compiler or runtime. This can lead to situations where you might inadvertently pass an incorrect type to the function, causing runtime errors or unexpected behavior when attempting to perform operations that are not valid for the given type(s).

`reflect.DeepEqual`, for instance, is not type  safe:

```go
func TestAll(t *testing.T) {

	got := SumAll([]int{1, 2, 3, 4, 5}, []int{6, 7, 8, 9, 10})
	want := "gopher"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v got %v", want, got)
	}
}
```
The code compiles (doesn't flag an error) but throws an error at runtime.

## Value Receiver vs Pointer Receiver

In Go, methods can have either value receivers or pointer receivers. The choice between them affects how the method interacts with the receiver and can have implications for performance and behavior.

### Value Receiver

- **Definition:** A method with a value receiver operates on a copy of the value it is called on.
- **Usage:** Use value receivers when the method does not need to modify the receiver or when the receiver is a small, simple type.

Example:
```go
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}
```

### Pointer Receiver
- **Definition:** A method with a pointer receiver operates on the actual value it is called on, allowing the method to modify the receiver.
- **Usage:** Use pointer receivers when the method needs to modify the receiver, or when the receiver is a large struct to avoid copying.

Example:
```go
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}
```
---

## Dependency Injection

Dependency Injection (DI) is a design pattern where an object or function's dependencies are provided (injected) by an external entity rather than the object or function creating the dependencies itself. The key idea behind DI is `to decouple the creation of dependencies from their usage`, making the code more modular, testable, and easier to maintain.

### Function Parameters vs. Dependency Injection

1. **Function Parameters (like in `Perimeter`)**:
   - When you pass a `Rectangle` to the `Perimeter` function, you’re just passing data (not a functionality like in the case of DI). The `Rectangle` is already created or defined somewhere in your code, and the `Perimeter` function merely operates on this data. This is not dependency injection because the `Rectangle` is not a "dependency" in the sense of a service, resource, or component that the function relies on to operate. It’s simply an argument.

```go
   type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}
```

2. **Dependency Injection**:
   - Dependency injection involves passing in **dependencies** that the function or struct needs to perform its tasks. These dependencies are typically objects, services, or interfaces that provide `specific functionality`, like logging, database access, or external services. The key point in DI is that the function or struct does not create these dependencies itself; they are injected from outside.

### Creating custom (domain-specific) types from existing types:

```go
type Bitcoin int
```

This creates a type called `Bitcoin` from the existing type `int`. To make `Bitcoin` you just use the syntax `Bitcoin(999)`; this converts `999` to Bitcoin 

### On the `Stringer` Interface

```go
type Stringer interface {
	String() string
}
```
This interface is defined in the fmt package and lets you define how your (user-defined) type is printed when used with the `%s` format string in prints.

```go
func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```
This ensures that `Bitcoin(10)`, the user-defined type, is returned as `10 BTC`

### On `error.New` method

`errors.New` creates a new `error` with a message of your choosing.

---
### Formatting verbs used:

- `%q`--adds quote to a returned string
- `%s`--for strings
- `%d`--for integers
- `%v`--for slices
- `%p`--prints memory addresses in base 16 notation with leading `0x`s
- `%#v`--print out struct with the values in its field
- `%f`--for float
- `%g`--prints more precise float decimals
---

> the `.Error()` method when called on a variable gets and returns the string in that variable

---

In Go, understanding the distinction between **interface types** and **concrete types** is fundamental for leveraging Go's type system effectively, especially when it comes to polymorphism and designing flexible APIs.

### 1. **Concrete Types**

**Definition:**  
*Concrete types are specific implementations that define the structure and behavior of data. These types provide both the data and the methods that operate on that data.* They are defined using the `struct`, `int`, `float`, `string`, arrays, slices, maps, and other built-in or user-defined types.

**Characteristics:**

- **Definite Structure:** Concrete types have a well-defined structure. For example, a `struct` type has explicitly defined fields, and an `int` type represents a specific kind of integer.
  
  ```go
  type Rectangle struct {
      width, height float64
  }
  ```

- **Direct Instantiation:** You can directly create instances of concrete types. For example:
  
  ```go
  rect := Rectangle{width: 10, height: 5}
  ```

- **Method Implementations:** Concrete types can have methods that are defined on them directly. These methods operate on the data defined within the type.

  ```go
  func (r Rectangle) Area() float64 {
      return r.width * r.height
  }
  ```

- **Specificity:** Concrete types are specific; they are exact implementations that define how data is stored and manipulated.

- **Memory Layout:** Concrete types have a specific memory layout, meaning the size and memory structure are known at compile time.

### 2. **Interface Types**

**Definition:**  
*Interfaces define a set of method signatures but do not implement them. They are abstract types that specify a contract or a set of behaviors that other types must implement.* An interface type is satisfied by any type that implements its methods, making it a cornerstone of polymorphism in Go.

**Characteristics:**

- **Abstract Contract:** Interfaces are abstract and do not hold any data themselves. They define behaviors that concrete types must adhere to.
  
  ```go
  type Shape interface {
      Area() float64
  }
  ```

- **Decoupling:** Interfaces allow you to decouple the definition of methods from their implementation. This decoupling enables you to write more flexible and reusable code.

- **Implicit Implementation:** In Go, a type implements an interface simply by implementing its methods. There's no explicit declaration needed. If a `Rectangle` type implements all the methods of a `Shape` interface, then `Rectangle` is considered to satisfy the `Shape` interface.

  ```go
  // No explicit declaration needed
  var s Shape = Rectangle{width: 10, height: 5} // Rectangle implements Shape
  ```

- **Dynamic Behavior:** Interfaces can hold any value that implements the defined methods, allowing you to treat different types uniformly based on shared behavior.

- **Polymorphism:** Interfaces enable polymorphism, where a single function can operate on different types of objects. For example, a function can accept a `Shape` interface and work with any type that implements the `Shape` interface, regardless of its concrete type.

  ```go
  func PrintArea(s Shape) {
      fmt.Println(s.Area())
  }
  ```

### 4. **Example to Illustrate Differences**

Let's use a practical example to illustrate how concrete and interface types work together.

#### Concrete Type Example:

```go
type Circle struct {
    radius float64
}

func (c Circle) Area() float64 {
    return 3.14 * c.radius * c.radius
}
```

- Here, `Circle` is a concrete type with a specific structure (`radius`) and a method (`Area`) that operates on that structure.

#### Interface Type Example:

```go
type Shape interface {
    Area() float64
}
```

- `Shape` is an interface type. It doesn't know how `Area` is calculated; it only specifies that any type that claims to be a `Shape` must have an `Area` method.

#### Using Both Together:

```go
func PrintArea(s Shape) {
    fmt.Println(s.Area())
}

func main() {
    c := Circle{radius: 5}
    PrintArea(c) // Circle is treated as a Shape
}
```

- In this example, `PrintArea` accepts a `Shape`. It can operate on any type (like `Circle`) that implements the `Shape` interface. This shows how interfaces allow for flexibility and polymorphism, as you can add other types (e.g., `Rectangle`) later without changing `PrintArea`.

### Summary

- **Concrete Types** are actual implementations with specific data and methods. They define how something is structured and how it behaves.
- **Interface Types** define a set of behaviors (methods) without specifying how those behaviors are implemented. They allow different concrete types to be treated uniformly based on shared behavior.

---

## The Single Responsibility Principle (SRP) and Seperation of Concern (SoC)

> Ref: `mocking/main.go` and `mocking/mocking_test.go`

Separating from specific implementations, like the sleep behavior in your `Countdown` example, is important for several reasons. These revolve around making the codebase more flexible, maintainable, and testable. Here’s why this separation is beneficial:

### 1. **Testability and Faster Tests**

- **Problem without Separation**: If your code directly uses a specific implementation like `time.Sleep(1 * time.Second)`, it introduces real time delays into your tests. This makes tests slow, which is problematic when you need to run them frequently during development.

- **Solution with Separation**: By separating the sleep behavior into an interface (`Sleeper`), you can inject a mock or a fake implementation during testing. For example, a `SpySleeper` that does nothing but record calls can be used. This way, your tests can run instantly without actual waiting:

    ```go
    type SpySleeper struct {
        Calls int
    }

    func (s *SpySleeper) Sleep() {
        s.Calls++
    }
    ```

This allows you to verify that the `Sleep()` method was called the expected number of times without introducing unnecessary delays.

### 2. **Decoupling for Flexibility**

- **Problem without Separation**: When your code is tightly coupled to a specific implementation, such as directly calling `time.Sleep()`, it becomes rigid. Changing how the sleep behavior works (e.g., to use a different timing mechanism or to add logging) would require modifying the code that uses it.

- **Solution with Separation**: By defining an interface and using dependency injection, you can easily swap out the implementation. This makes the code more flexible and adaptable to change. For instance, if you need a different `Sleeper` that logs every sleep action for debugging purposes, you can simply create a new type:

    ```go
    type LoggingSleeper struct{}

    func (l *LoggingSleeper) Sleep() {
        fmt.Println("Sleeping...")
        time.Sleep(1 * time.Second)
    }
    ```

Now, you can inject `LoggingSleeper` wherever you need it, without changing the core logic of `Countdown`.

### 3. **Enhanced Maintainability**

- **Problem without Separation**: Tightly coupled code is harder to maintain and understand. If the sleep logic is embedded within the main function logic, any changes to the sleep behavior would require touching the main function code. This increases the risk of bugs and makes the code less modular.

- **Solution with Separation**: When behavior is encapsulated in separate implementations that adhere to an interface, you can update or replace these implementations without affecting other parts of the codebase. It keeps the code modular and the responsibilities clearly defined, leading to easier maintenance.

### 4. **Promoting Reusability**

- **Problem without Separation**: If sleep behavior is not separated, any other part of the application that needs to sleep will have to duplicate the logic or rely on copying and pasting code. This results in code duplication and potential inconsistencies.

- **Solution with Separation**: By using an interface, you can define a `Sleeper` type once and reuse it across different parts of your application. This promotes code reuse and consistency. For instance, both a `Countdown` function and another function that simulates delays can use the same `Sleeper` interface without duplicating code.

### 5. **Better Abstraction and Clean Code Principles**

- **Problem without Separation**: Code that handles multiple responsibilities can become complex and harder to understand. If your `Countdown` function handles both the countdown logic and the specific sleep behavior, it violates the Single Responsibility Principle (SRP).

- **Solution with Separation**: Separating the sleep behavior into its own component adheres to clean code principles. Each part of the code has a clear, distinct role. The `Countdown` function focuses solely on counting down, while the `Sleeper` interface and its implementations manage sleeping.

### 6. **Easier to Mock and Spy for Behavior Verification**

- **Problem without Separation**: Without the ability to inject different implementations, it's hard to verify certain behaviors. For instance, how do you confirm that your countdown logic properly waits for each second to pass if you're directly calling `time.Sleep()`?

- **Solution with Separation**: Using a `SpySleeper` allows you to verify that `Sleep()` is called the expected number of times. You can inspect the internal state of the spy to assert that behavior is as expected, which is crucial for reliable testing.

---

## Side-effects

In the context of programming, **side effects** refer to any observable changes or interactions a function or expression has with the outside world, beyond returning a value. These changes can include modifications to variables or data structures outside the function's scope, I/O operations (like printing to the console, writing to a file, or making network requests), or altering the state of a system in any other way.

### Examples of Side Effects

1. **Modifying Global or External State**: A function that changes a global variable or modifies the state of an object that exists outside its local scope introduces a side effect.

    ```go
    var counter int

    func incrementCounter() {
        counter++ // Modifying a global variable
    }
    ```

    Here, `incrementCounter()` has a side effect because it changes the value of the global variable `counter`.

2. **I/O Operations**: Functions that perform input/output operations such as printing to the console, writing to a file, or reading user input have side effects.

    ```go
    func greetUser(name string) {
        fmt.Printf("Hello, %s!\n", name) // Output to console
    }
    ```

    The call to `fmt.Printf` is a side effect because it changes the program's output to the console.

3. **Network Requests**: Functions that make network requests, such as sending HTTP requests or connecting to a database, produce side effects by interacting with external systems.

    ```go
    func sendMessage(url string, message string) error {
        _, err := http.Post(url, "application/json", strings.NewReader(message))
        return err
    }
    ```

    This function has a side effect because it sends data over the network.

4. **Mutating Input Arguments**: If a function changes the content of its input arguments (assuming they are pointers or references), it has a side effect.

    ```go
    func updateValue(val *int) {
        *val = 10 // Modifying the value at the pointer address
    }
    ```

    Here, `updateValue` modifies the value that the pointer `val` points to, creating a side effect.

5. **Modifying Data Structures**: A function that alters a data structure passed to it also introduces side effects.

    ```go
    func addElement(slice *[]int, element int) {
        *slice = append(*slice, element) // Modifying the slice
    }
    ```

    This function modifies the slice that is passed to it, leading to a side effect.

### Why Are Side Effects Important?

1. **Predictability and Debugging**: Functions without side effects (pure functions) are easier to predict, understand, and test. If a function depends only on its input arguments and has no side effects, it will always produce the same output for the same input, which simplifies debugging and reasoning about code.

2. **Testing and Test Isolation**: Side effects can make functions harder to test. If a function modifies global state, interacts with the file system, or depends on external systems, you need to set up and clean up these states during testing, which can complicate test cases.

3. **Concurrency Issues**: Side effects can lead to concurrency problems like race conditions. When multiple threads or routines modify shared state simultaneously, it can cause unpredictable behavior and bugs.

4. **Functional Programming Principles**: Many functional programming paradigms emphasize writing pure functions that do not have side effects. This approach promotes immutability, makes code more modular, and reduces unintended interactions between different parts of a program.

### Managing Side Effects

- **Encapsulation**: Encapsulate side effects within specific parts of your codebase. For instance, use service layers to handle external interactions (e.g., database access or network calls), keeping core logic free of side effects.

- **Dependency Injection**: Use dependency injection to control side effects. Inject dependencies like loggers, database connections, or network clients, so you can substitute them with mock implementations during testing.

- **Pure Functions**: Whenever possible, write pure functions that do not have side effects. This makes your code more predictable and easier to test.

- **Explicit Interfaces**: Clearly define interfaces that separate pure logic from functions that handle side effects. For example, use an interface for a logger, so you can inject a no-op logger in tests to avoid actual I/O operations.
