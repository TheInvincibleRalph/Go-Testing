package main

import (
	"log"
	"net/http"
)

// type InMemoryPlayerStore struct{}

// func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
// 	return 123
// }

// func (i *InMemoryPlayerStore) RecordWin(name string) {}

func main() {
	server := &PlayerServer{NewInMemoryPlayerStore()}
	log.Fatal(http.ListenAndServe(":5000", server))
}

/*
Separation of Concerns:
The code separates the concerns of handling HTTP requests (PlayerServer),
retrieving data (PlayerStore), and running the server (main).
This makes the code more modular, easier to test, and maintainable.
*/
