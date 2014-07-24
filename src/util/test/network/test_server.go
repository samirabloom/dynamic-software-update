package network

import (
	"strconv"
	"net/http"
	"fmt"
	"time"
)

// ==== TEST_SERVER - START

func Test_server(ports []int) {
	for _, port := range ports {
		fmt.Printf("Starting server %d ...\n", port)
		go http.ListenAndServe(":"+strconv.Itoa(port), &handle{port})
	}
	time.Sleep(1000 * time.Millisecond)
}

type handle struct {
	Port int
}

func (h *handle) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Port: %d\n", h.Port)

}

// ==== TEST_SERVER - END
