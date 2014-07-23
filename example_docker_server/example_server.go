package main

import (
	"strconv"
	"net/http"
	"time"
	"fmt"
	"flag"
)

func main() {
	portFlag := flag.String("port", "1025", "Set the server's port")
	flag.Parse()

	port, _ := strconv.Atoi(*portFlag)
	fmt.Printf("Starting server %d ...\n", port)
	http.ListenAndServe(":"+strconv.Itoa(port), &handle{port})
}

// ==== SERVER - START

type handle struct {
	Port int
}

func (h *handle) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate");
	response.Header().Set("Pragma", "no-cache");
	response.Header().Set("Expires", "0");
	fmt.Fprintf(response, "Port: %d\n", h.Port)
	time.Sleep(50 * time.Millisecond)
	fmt.Fprintf(response, "50 ms, ")
	response.(http.Flusher).Flush()
	time.Sleep(50 * time.Millisecond)
	fmt.Fprintf(response, "100 ms, ")
	response.(http.Flusher).Flush()
	time.Sleep(50 * time.Millisecond)
	fmt.Fprintf(response, "150 ms, ")
	response.(http.Flusher).Flush()
	time.Sleep(50 * time.Millisecond)
	fmt.Fprintf(response, "200 ms, ")
	response.(http.Flusher).Flush()
	time.Sleep(50 * time.Millisecond)
	fmt.Fprintf(response, "250 ms\n")
	response.(http.Flusher).Flush()
}

// ==== SERVER - END
