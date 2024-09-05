package context

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type Store interface {
	Fetch(ctx context.Context) (string, error)
}

func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := store.Fetch(r.Context())

		if err != nil {
			return
		}
		fmt.Fprint(w, data)
	}
}

// Spy for store
type SpyStore struct {
	response string
	t        *testing.T
}

type SpyResponseWriter struct {
	written bool
}

func (s *SpyResponseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
	s.written = true
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	go func() {
		var result string
		for _, c := range s.response {
			select {
			case <-ctx.Done():
				log.Println("spy store got cancelled")
				return
			default:
				time.Sleep(10 * time.Millisecond)
				result += string(c)
			}
		}
		data <- result
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-data:
		return res, nil
	}
}

func TestServer(t *testing.T) {
	t.Run("returns data from store", func(t *testing.T) {
		data := "hello, world"
		store := &SpyStore{response: data, t: t}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(`expected %s, got %s`, data, response.Body.String())
		}
	})

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {

		data := "hello, world"
		store := &SpyStore{response: data, t: t}
		srv := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)

		// Create a new context with a timeout of 5ms
		// request.Context() gets and returns the context (the instance of context.Context interface) of the request. For incoming server requests, the context is canceled when the client's connection closes
		// context.WithCancel is used to change the default context to a new context that can be cancelled. It returns a new context and a cancel function.
		// cancellingCtx is therefore the child context of the request context
		cancellingCtx, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		request = request.WithContext(cancellingCtx) // request is now the child context of the previous request context (that was defined in request.Context())

		response := &SpyResponseWriter{}

		srv.ServeHTTP(response, request)

		if response.written {
			t.Error("a response should not have been written")
		}
	})

}
