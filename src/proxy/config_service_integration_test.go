package proxy

import (
	"fmt"
	"testing"
	"code.google.com/p/go-uuid/uuid"
	networkutil "util/test/network"
	assertion "util/test/assertion"
	"strconv"
	"time"
	"proxy/contexts"
	"proxy/docker_client"
)

func Test_Config_PUT_GET_DELETE(testCtx *testing.T) {
	// given - a config server
	var (
		serverPort = networkutil.FindFreeLocalSocket(testCtx).Port
		serverUrl  = "http://127.0.0.1:" + strconv.Itoa(int(serverPort)) + "/configuration/cluster"
	)
	go ConfigServer(serverPort, &contexts.Clusters{}, &docker_client.DockerHost{})

	time.Sleep(150 * time.Millisecond)

	// when
	// - a PUT request
	uuidResponse, putStatus := networkutil.PUTRequest(serverUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":1024}, {\"hostname\":\"127.0.0.1\", \"port\":1025}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"1.1\"}}")

	// then
	assertion.AssertDeepEqual("Correct PUT Status", testCtx, "202 Accepted", putStatus)
	if uuid.Parse(uuidResponse) == nil {
		testCtx.Fatal(fmt.Errorf("\nInvalid UUID returned from request, response was [%s]", uuidResponse))
	}



	// when
	// - a GET request
	jsonResponse, getStatus := networkutil.GETRequest(serverUrl + "/" + uuidResponse)

	// then
	assertion.AssertDeepEqual("Correct PUT Status", testCtx, "200 OK", getStatus)
	assertion.AssertDeepEqual("Correct GET Response", testCtx, "{\n    \"cluster\": {\n        \"servers\": [\n            {\n                \"hostname\": \"127.0.0.1\",\n                \"port\": 1024\n            },\n            {\n                \"hostname\": \"127.0.0.1\",\n                \"port\": 1025\n            }\n        ],\n        \"upgradeTransition\": {\n            \"mode\": \"SESSION\",\n            \"sessionTimeout\": 1\n        },\n        \"uuid\": \"" + uuidResponse + "\",\n        \"version\": \"1.1\"\n    }\n}", jsonResponse)



	// when
	// - a DELETE request
	_, deleteStatus := networkutil.DELETERequest(serverUrl + "/" + uuidResponse)

	// then
	assertion.AssertDeepEqual("Correct DELETE Status", testCtx, "202 Accepted", deleteStatus)



	// when
	// - another GET response
	_, getAfterDeleteStatus := networkutil.GETRequest(serverUrl + "/" + uuidResponse)

	// then
	assertion.AssertDeepEqual("Correct PUT Status", testCtx, "404 Not Found", getAfterDeleteStatus)



	// when
	// - another DELETE request
	_, deleteAfterDeleteStatus := networkutil.DELETERequest(serverUrl + "/" + uuidResponse)

	// then
	assertion.AssertDeepEqual("Correct DELETE Status", testCtx, "404 Not Found", deleteAfterDeleteStatus)
}
