package network

import (
	"net/http"
	"io/ioutil"
	"strings"
	"log"
	"io"
)

func PUTRequest(url string, json string) (body, status string) {
	return MakeRequest("PUT", url, strings.NewReader(json), "")
}

func GETRequest(url string) (body, status string) {
	return MakeRequest("GET", url, nil, "")
}

func GETCookiedRequest(url string, uuidCookie string) (body, status string) {
	return MakeRequest("GET", url, nil, uuidCookie)
}

func DELETERequest(url string) (body, status string) {
	return MakeRequest("DELETE", url, nil, "")
}

func MakeRequest(method, url string, request io.Reader, uuidCookie string) (body, status string) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, request)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Connection", "close")
	if len(uuidCookie) > 0 {
		req.Header.Add("Cookie", "dynsoftup="+uuidCookie+";")
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




