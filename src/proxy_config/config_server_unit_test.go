package proxy_config

import (
	"testing"
	"bytes"
	"net/http"
	"fmt"
	"encoding/json"
	"regexp"
	"net/url"
	mock "util/test/mock"
	assertion "util/test/assertion"
)

// TODO - this file needs to be cleaned up!!

func Test_Config_PUT_With_Valid_Json_Object(testCtx *testing.T) {
	// given
	var (
		jsonObjectMaps map[string]interface{} = make(map[string]interface{})
		responseWriter                        = &mock.MockResponseWriter{WritenBodyBytes: make(map[int][]byte)}
		bodyByte                              = []byte("{\"name_one\":\"value_one\", \"name_two\":\"value_two\"}")
		request                               = &http.Request{Body: &mock.MockBody{BodyBytes: bodyByte}}
		uuid string                           = "uuid"
		jsonObject interface{}
	)
	json.Unmarshal(bodyByte, &jsonObject)

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

func Test_Config_GET_With_Existing_Object(testCtx *testing.T) {
	// given
	var (
		jsonObjectMaps map[string]interface{} = make(map[string]interface{})
		responseWriter                        = &mock.MockResponseWriter{WritenBodyBytes: make(map[int][]byte)}
		bodyByte                              = []byte("{\"id\":\"uuid\",\"name_one\":\"value_one\",\"name_two\":\"value_two\"}")
		request                               = &http.Request{URL: &url.URL{Path: "/server/uuid"}}
		urlRegex                              = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		jsonObject interface{}
	)
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

func Test_Config_DELETE_With_Existing_Object(testCtx *testing.T) {
	// given
	var (
		jsonObjectMaps map[string]interface{} = make(map[string]interface{})
		responseWriter                        = &mock.MockResponseWriter{WritenBodyBytes: make(map[int][]byte)}
		bodyByte                              = []byte("{\"id\":\"uuid\",\"name_one\":\"value_one\", \"name_two\":\"value_two\"}")
		request                               = &http.Request{URL: &url.URL{Path: "/server/uuid"}}
		urlRegex                              = regexp.MustCompile("/server/([a-z0-9-]*){1}")
		jsonObject interface{}
	)
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
