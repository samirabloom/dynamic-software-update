package proxy_config

import (
	"net/http"
	"io/ioutil"
	"strings"
	"log"
)

//const url = "https://127.0.0.1:8080/server"

//const url = "https://foaas.herokuapp.com/shakespeare/Ali/Hilda"

//func main() {
//	//	go Server(8080)
//	uuid := PUTRequest("http://127.0.0.1:8080/server", "{\"name_one\":\"value_one\", \"name_two\":\"value_two\"}")
//	fmt.Printf("UUID Response %s\n", uuid)
//
//	jsonResponse := GETRequest("http://127.0.0.1:8080/server/" + uuid)
//	fmt.Printf("JSON Response %s", jsonResponse)
//
//	deleteResponse := DELETERequest("http://127.0.0.1:8080/server/" + uuid)
//	fmt.Printf("\nresponse after delete %s", deleteResponse)
//
//	reJsonResponse := GETRequest("http://127.0.0.1:8080/server/" + uuid)
//	fmt.Printf("JSON Response %s", reJsonResponse)
//
//}

func PUTRequest(url string, json string) string {
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, strings.NewReader(json))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
	return ""
}

func GETRequest(url string) ([]byte) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
		defer resp.Body.Close()
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return body
}


func DELETERequest(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
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
		return string(body)
	}
	return ""
}





