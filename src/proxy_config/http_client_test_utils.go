package proxy_config

import (
	"net/http"
	"io/ioutil"
	"strings"
	"log"
	"io"
)

func PUTRequest(url string, json string) (body, status string) {
	return makeRequest("PUT", url, strings.NewReader(json))
}

func GETRequest(url string) (body, status string) {
	return makeRequest("GET", url, nil)
}


func DELETERequest(url string) (body, status string) {
	return makeRequest("DELETE", url, nil)
}

func makeRequest(method, url string, request io.Reader) (body, status string) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, request)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	} else {
		defer resp.Body.Close()
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), resp.Status
	}
	return "", resp.Status
}




