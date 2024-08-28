package concurrency

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, err error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

// Whichever function writes to its channel first will have its code executed in the select,
// which results in its URL being returned
func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, err error) {
	select {
	case <-ping(a): // listening to the channel returned by ping
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout): // sends a signal if neither a and b returns
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch // returns an empty struct type channel used to signal the completion of the ping operation
}

func TestRacer(t *testing.T) {
	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
		// Setting Up the Test Servers
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		// Extracting the URLs
		slowURL := slowServer.URL
		fastURL := fastServer.URL

		// Insert the URL into the function under test
		want := fastURL
		got, err := Racer(slowURL, fastURL)

		if got != want {
			t.Errorf("expected %q got %q", want, got)
		}

		if err != nil {
			t.Fatalf("did not expect an error but got one %v", err)
		}
	})

	t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
		server := makeDelayedServer(25 * time.Second)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}

	})

}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
