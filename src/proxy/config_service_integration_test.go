package proxy

import (
	"fmt"
	"testing"
	"code.google.com/p/go-uuid/uuid"
	networkutil "util/test/network"
	assertion "util/test/assertion"
	"strconv"
	"time"
	"proxy/stages"
)

func Test_Config_PUT_GET_DELETE(testCtx *testing.T) {
	// given - a config server
	var (
		serverPort = networkutil.FindFreeLocalSocket(testCtx).Port
		serverUrl  = "http://127.0.0.1:" + strconv.Itoa(int(serverPort)) + "/server"
	)
	go ConfigServer(serverPort, &stages.Clusters{})

	time.Sleep(150 * time.Millisecond)

	// when
	// - a PUT request
	uuidResponse, putStatus := networkutil.PUTRequest(serverUrl, "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":1024}, {\"ip\":\"127.0.0.1\", \"port\":1025}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.1}}")

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
	assertion.AssertDeepEqual("Correct GET Response", testCtx, "{\"cluster\":{\"servers\":[{\"ip\":\"127.0.0.1\",\"port\":1024},{\"ip\":\"127.0.0.1\",\"port\":1025}],\"upgradeTransition\":{\"mode\":\"SESSION\",\"sessionTimeout\":1},\"uuid\":\""+uuidResponse+"\",\"version\":1.1}}", jsonResponse)



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
