# Testing Terminologies

## Stub
A stub typically provides hardcoded or predefined responses for specific method calls. It’s very simple and often used to isolate parts of the system in tests. The purpose of a stub is to provide controlled and predictable responses for testing, not to maintain state or handle complex logic.

Style
Mocks vs Stubs = Behavioral testing vs State testing

Principle
According to the principle of Test only one thing per test, there may be several stubs in one test, but generally there is only one mock.


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