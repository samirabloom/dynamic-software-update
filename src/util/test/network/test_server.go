package network

import (
	"strconv"
	"net/http"
	"fmt"
	"time"
	"regexp"
)

// ==== TEST_SERVER - START

func Test_server(ports []int, crash bool) {
	for _, port := range ports {
		go http.ListenAndServe(":"+strconv.Itoa(port), &handle{port: port, crash: crash})
	}
	time.Sleep(150 * time.Millisecond)
}

type handle struct {
	port int
	crash bool
}

func (h *handle) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if (h.crash && regexp.MustCompile("/crash").MatchString(request.URL.Path)) {
		panic("simulating server crash")
	}
	fmt.Fprintf(response, "Port: %d\n", h.port)
}

// ==== TEST_SERVER - END
