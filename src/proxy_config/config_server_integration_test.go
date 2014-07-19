package proxy_config

import (
	"fmt"
	"testing"
	"code.google.com/p/go-uuid/uuid"
	networkutil "util/test/network"
	assertion "util/test/assertion"
	"strconv"
)

func Test_Config_PUT_GET_DELETE(testCtx *testing.T) {
	// given - a config server
	var (
		serverPort = networkutil.FindFreeLocalSocket(testCtx).Port
		serverUrl = "http://127.0.0.1:" + strconv.Itoa(serverPort) + "/server"
	)
	go ConfigServer(serverPort, nil)



	// when
	// - a PUT request

	uuidResponse, putStatus := PUTRequest(serverUrl, "{\"name_one\":\"value_one\", \"name_two\":\"value_two\"}")

	// then
	assertion.AssertDeepEqual("Correct PUT Status", testCtx, "202 Accepted", putStatus)
	if uuid.Parse(uuidResponse) == nil {
		testCtx.Fatal(fmt.Errorf("\nInvalid UUID returned from request, response was [%s]", uuidResponse))
	}



	// when
	// - a GET request
	jsonResponse, getStatus := GETRequest(serverUrl + "/" + uuidResponse)

	// then
	assertion.AssertDeepEqual("Correct PUT Status", testCtx, "200 OK", getStatus)
	assertion.AssertDeepEqual("Correct GET Response", testCtx, "{\"id\":\"" + uuidResponse + "\",\"name_one\":\"value_one\",\"name_two\":\"value_two\"}", jsonResponse)



	// when
	// - a DELETE request
	_, deleteStatus := DELETERequest(serverUrl + "/" + uuidResponse)

	// then
	assertion.AssertDeepEqual("Correct DELETE Status", testCtx, "202 Accepted", deleteStatus)



	// when
	// - another GET response
	_, getAfterDeleteStatus := GETRequest(serverUrl + "/" + uuidResponse)

	// then
	assertion.AssertDeepEqual("Correct PUT Status", testCtx, "404 Not Found", getAfterDeleteStatus)



	// when
	// - another DELETE request
	_, deleteAfterDeleteStatus := DELETERequest(serverUrl + "/" + uuidResponse)

	// then
	assertion.AssertDeepEqual("Correct DELETE Status", testCtx, "404 Not Found", deleteAfterDeleteStatus)
}

