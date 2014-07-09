package proxy_config

import (
	"fmt"
	"bytes"
	"strings"
	"testing"
)

func TestConfigServerEndToEnd(testCtx *testing.T) {
	// put request
	uuid := PUTRequest("http://127.0.0.1:8080/server", "{\"name_one\":\"value_one\", \"name_two\":\"value_two\"}")
	fmt.Printf("UUID Response %s\n", uuid)

	// get request
	// given
	var ExpectedJsonResponse = []byte("{\"id\":\"" + uuid + "\",\"name_one\":\"value_one\",\"name_two\":\"value_two\"}")

	//	// when
	jsonResponse := GETRequest("http://127.0.0.1:8080/server/" + uuid)
	fmt.Printf("JSON Response %s\n", jsonResponse)

	// then
	if !bytes.Equal(ExpectedJsonResponse, jsonResponse) {
		testCtx.Fatal(fmt.Errorf("\nUUID response incorrect\n\nExpected:\n[%s]\nActual:\n[%s]", ExpectedJsonResponse, jsonResponse))
	}

	// delete request
	// given
	var response string = "Accepted"
	// when
	deleteResponse := DELETERequest("http://127.0.0.1:8080/server/" + uuid)
	fmt.Printf("\nresponse after delete %s\n", deleteResponse)
	// then
	if !strings.EqualFold(response, deleteResponse) {
		testCtx.Fatal(fmt.Errorf("\nUUID response incorrect\n\nExpected:\n[%s]\nActual:\n[%s]", response, deleteResponse))
	}

	// get response
	ResponseAfterDelete := GETRequest("http://127.0.0.1:8080/server/" + uuid)

	// then
	var expectedResponse string = "404 page not found\n"
	fmt.Printf("\nresponse after delete: %s\n", ResponseAfterDelete)

	if !strings.EqualFold(expectedResponse, string(ResponseAfterDelete)) {
		testCtx.Fatal(fmt.Errorf("\nexpected:[%s]\n Actual:<%s>", expectedResponse, ResponseAfterDelete))
	}
}

func TestToDeleteTheNonExistingJsonObject(testCtx *testing.T) {
	// delete request
	// given
	var response string = "404 page not found\n"
	// when
	deleteResponse := DELETERequest("http://127.0.0.1:8080/server/" + "Non_existing_uuid")
	fmt.Printf("\nresponse after delete: %s\n", deleteResponse)
	// then
	if !strings.EqualFold(response, deleteResponse) {
		testCtx.Fatal(fmt.Errorf("\nUUID response incorrect\n\nExpected:\n[%s]\nActual:\n[%s]", response, deleteResponse))
	}
}

