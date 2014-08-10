package network

import (
	"net/http"
	"io/ioutil"
	"strings"
	"log"
	"io"
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
	defer resp.Body.Close()

	if err != nil {
		log.Fatalln(err)
		return "", resp.Status
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), resp.Status
	}
}




