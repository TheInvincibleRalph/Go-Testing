package concurrency

import (
	"reflect"
	"testing"
	"time"
)

type WebsiteChecker func(string) bool // a function type that takes in a string and returns a bool

type result struct {
	string
	bool
}

// func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
// 	results := make(map[string]bool)

// 	for _, url := range urls {
// 		results[url] = wc(url)
// 	}
// 	return results
// }

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)

	for _, url := range urls {
		go func(u string) {
			resultChannel <- result{u, wc(u)} // <- is called a send statement (variable to channel)
		}(url)

	}

	for i := 0; i < len(urls); i++ {
		r := <-resultChannel // := <- is called a receive expression (channel to variable)
		results[r.string] = r.bool
	}

	return results
}

func mockWebsiteChecker(url string) bool {
	return url != "waat://furhurterwe.geds"
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v got %v", want, got)
	}

}

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < len(urls); i++ {
		urls[i] = "a url"
	}
	b.ResetTimer() //this line resets the timer to ensure that the time spent on the setup abov is not included in the benchmark results.
	for i := 0; i < b.N; i++ {
		CheckWebsites(slowStubWebsiteChecker, urls)
	}
}

/* Task description

Instead of waiting for a website to respond before sending a request to the next website, we will tell our computer to make the next request while it is waiting.
*/
