package proxy_config

import (
	"testing"
	"bytes"
	"net/http"
	"util/test/mock"
	"fmt"
	"encoding/json"
	"regexp"
	"net/url"
	"util/test/assertion"
)

func TestShouldPutJsonObject(testCtx *testing.T) {
	// given
	var (
		jsonObjectMaps map[string]interface{} = make(map[string]interface{})
		responseWriter                        = &mock.MockResponseWriter{WritenBodyBytes: make(map[int][]byte)}
		request                               = &http.Request{}
		uuid string                           = "uuid"
		bodyByte                              = []byte("{\"name_one\":\"value_one\", \"name_two\":\"value_two\"}")
		jsonObject interface{}
	)
	json.Unmarshal(bodyByte, &jsonObject)
	request.Body = &mock.MockBody{BodyBytes: bodyByte}

	// when
	PUTHandler(func() string {
		return uuid
	})(jsonObjectMaps, responseWriter, request)

	// then
	if !bytes.Equal([]byte(uuid), responseWriter.WritenBodyBytes[0]) {
		testCtx.Fatal(fmt.Errorf("\nUUID response incorrect\n\nExpected:\n[%s]\nActual:\n[%s]", []byte(uuid), responseWriter.WritenBodyBytes[0]))
	}

	jsonObject.(map[string]interface{})["id"] = uuid
	assertion.AssertDeepEqual("Json objects match", testCtx, jsonObject, jsonObjectMaps[uuid])
}

func TestShouldGetJsonObject(testCtx *testing.T) {
	// given
	var (
		jsonObjectMaps map[string]interface{} = make(map[string]interface{})
		responseWriter                        = &mock.MockResponseWriter{WritenBodyBytes: make(map[int][]byte)}
		request                               = &http.Request{}
		bodyByte                              = []byte("{\"id\":\"uuid\",\"name_one\":\"value_one\",\"name_two\":\"value_two\"}")
		jsonObject interface{}
	)
	urlRegex := regexp.MustCompile("/server/([a-z0-9-]*){1}")
	request.URL = &url.URL{}
	request.URL.Path = "/server/uuid"
	json.Unmarshal(bodyByte, &jsonObject)
	jsonObjectMaps["uuid"] = jsonObject

	// when
	GETHandler(urlRegex)(jsonObjectMaps, responseWriter, request)
	println(string(responseWriter.WritenBodyBytes[0]))
	value := make([]byte, len(responseWriter.WritenBodyBytes[0]))
	copy(value, responseWriter.WritenBodyBytes[0])
	fmt.Printf("\nBodyByte: [%s]\nResponse: <%s>\n\n", bodyByte, value)

	// then
	if !bytes.Equal(bodyByte, value) {
		testCtx.Fatalf("expected: [%s] actual: [%s]", bodyByte, value)
		testCtx.Fatalf("[%s] - [%s]", bodyByte, value)
	}
	assertion.AssertDeepEqual("Json objects match", testCtx, jsonObject, jsonObjectMaps["uuid"])
}

func TestShouldDeleteJsonObject(testCtx *testing.T) {
	// given
	var (
		jsonObjectMaps map[string]interface{} = make(map[string]interface{})
		responseWriter                        = &mock.MockResponseWriter{WritenBodyBytes: make(map[int][]byte)}
		request                               = &http.Request{}
		bodyByte                              = []byte("{\"id\":\"uuid\",\"name_one\":\"value_one\", \"name_two\":\"value_two\"}")
		jsonObject interface{}
	)
	urlRegex := regexp.MustCompile("/server/([a-z0-9-]*){1}")
	request.URL = &url.URL{}
	request.URL.Path = "/server/uuid"
	json.Unmarshal(bodyByte, &jsonObject)
	jsonObjectMaps["uuid"] = jsonObject

	// when
	DeleteHandler(urlRegex)(jsonObjectMaps, responseWriter, request)

	// then
	fmt.Printf("\nresponse object map %v\n", jsonObjectMaps["uuid"])
	if jsonObjectMaps["uuid"] != nil {
		testCtx.Fatal(fmt.Errorf("\nexpected: nil \nactual: %s\n ", jsonObjectMaps["uuid"]))
	}
}
