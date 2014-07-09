package main

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
)

func main() {
	Server(8080)
}

func Server(port int) {
	fmt.Println("Starting server " + strconv.Itoa(port) + " ....")
	http.ListenAndServe(":"+strconv.Itoa(port), &handle{port})
}

type handle struct {
	Port int
}

func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "ONE, ")
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "TWO, ")
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "THREE, ")
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "FOUR, ")
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "FIVE from %d\n", h.Port)
}
