package network

import (
	"time"
	"fmt"
	"net/http"
)



// ==== SERVER - START

type Handle1 struct {
Port int
}

func (h *Handle1) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	bodyBytes := make([]byte, 1024)
	len, _ := request.Body.Read(bodyBytes)
	fmt.Printf("\nRequest to the server with status code: 500 and Port: %d\n%s", h.Port, bodyBytes[0:len])
	response.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(response, "Response from the server with status code: 500 and Port: %d\n", h.Port)
	time.Sleep(50 * time.Millisecond)
}

type Handle2 struct {
	Port int
}

func (h *Handle2) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	bodyBytes := make([]byte, 1024)
	len, _ := request.Body.Read(bodyBytes)
	fmt.Printf("\nRequest to the server with status code: 200 and Port: %d\n%s", h.Port, bodyBytes[0:len])
	response.WriteHeader(http.StatusOK)
	fmt.Fprintf(response, "Server responded with status code: 200 and Port: %d\n", h.Port)
	time.Sleep(50 * time.Millisecond)

}
