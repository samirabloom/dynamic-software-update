package network

import (
	"net/http"
	"io/ioutil"
	"strings"
	"log"
	"io"
	"fmt"
)

type Header struct {
	Name  string
	Value string
}

func PUTRequest(url string, json string) (body, status string) {
	return MakeRequest("PUT", url, strings.NewReader(json))
}

func GETRequest(url string) (body, status string) {
	return MakeRequest("GET", url, nil)
}

func GETRequestWithHeaders(url string, headers ...*Header) (body, status string) {
	return MakeRequest("GET", url, nil, headers...)
}

func DELETERequest(url string) (body, status string) {
	return MakeRequest("DELETE", url, nil)
}

func MakeRequest(method, url string, request io.Reader, headers ...*Header) (body, status string) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, request)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Connection", "close")
	for _, header := range headers {
		if header != nil {
			req.Header.Add(header.Name, header.Value)
		}
	}
	resp, err := client.Do(req)
	defer func() {
		fmt.Printf("5")
		if resp.Body != nil {
			fmt.Printf("6")
			resp.Body.Close()
		}
	}();

	fmt.Printf("0")
	if err != nil {
		fmt.Printf("1")
		log.Fatalln(err)
		fmt.Printf("2")
		return "", resp.Status
	} else {
		fmt.Printf("3")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("4")
		resp.Body.Close()
		fmt.Printf("5")
		return string(body), resp.Status
	}
}




