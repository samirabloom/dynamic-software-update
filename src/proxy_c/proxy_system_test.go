package proxy_c

import (
	"testing"
	networkutil "util/test/network"
	assertion "util/test/assertion"
	"io/ioutil"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeConfigFile(proxyPort int, configPort int, serverPorts []int) string {
	fileName := "system_test_config.json"
	data := "{\"proxy\":{\"ip\":\"localhost\",\"port\":" + strconv.Itoa(proxyPort) + "},\"configService\":{\"port\":" + strconv.Itoa(configPort) + "},\"cluster\":{\"servers\":["
	for index, serverPort := range serverPorts {
		if index > 0 {
			data += ","
		}
		data += "{\"ip\":\"127.0.0.1\",\"port\":"+strconv.Itoa(serverPort)+"}"
	}
	data += "]}}"

	ioutil.WriteFile("system_test_config.json", []byte(data), 0644)
	return fileName
}

func makeProxyRequest(proxyPort int, uuidCookie string) string {
	body, _ := networkutil.GETCookiedRequest("http://127.0.0.1:"+strconv.Itoa(proxyPort), uuidCookie)
	return body
}

func Test_proxy_system(testCtx *testing.T) {
	logLevel = new(string)
	*logLevel = "INFO"

	var (
		proxyPort int               = 1236
		configPort int              = 1237
		serverPortsClusterOne []int = []int{1055, 1056}
		serverPortsClusterTwo []int = []int{1055, 1056}
		serverPortsCluster3re []int = []int{1057, 1058}
		serverPortsCluster4or []int = []int{1059, 1060}
		serverPortsCluster5iv []int = []int{1061, 1062}
	)

	// given
	go Proxy(writeConfigFile(proxyPort, configPort, serverPortsClusterOne))
	networkutil.Test_server(serverPortsClusterOne)

	// then - should load balance requests
	assertion.AssertDeepEqual("Initial Config - Correct 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Config - Correct 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Config - Correct 3rd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Config - Correct 4th response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, ""))



	// given - new cluster - with session timeout and default upgrade
	networkutil.Test_server(serverPortsClusterTwo)
	uuidCookieVersion1_1, putStatus := networkutil.PUTRequest("http://127.0.0.1:"+strconv.Itoa(configPort)+"/server", "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.1}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Default Update - Correct 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Default Update - Correct 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, ""))

	// then - should load balance requests against new cluster
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Default Update - Correct 3rd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Default Update - Correct 4th response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))



	// given - another new cluster - with session timeout and session upgrade
	networkutil.Test_server(serverPortsCluster3re)
	uuidCookieVersion1_5, _ := networkutil.PUTRequest("http://127.0.0.1:"+strconv.Itoa(configPort)+"/server", "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"mode\":\"SESSION\",\"sessionTimeout\":1}, \"version\": 1.5}}")

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Session Update - Correct 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Session Update - Correct 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, ""))

	// then - send request to previous cluster if they have previous uuid
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Session Update - Correct 3rd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Session Update - Correct 4th response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Session Update - Correct 5th response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_5))
	assertion.AssertDeepEqual("Update Cluster Session Timeout & Session Update - Correct 6th response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_5))



	// given - another new cluster - no upgrade defined
	networkutil.Test_server(serverPortsCluster4or)
	uuidCookieVersion2_0, _ := networkutil.PUTRequest("http://127.0.0.1:"+strconv.Itoa(configPort)+"/server", "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster4or[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster4or[1])+"}], \"version\": 2.0}}")

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Update Cluster No Update - Correct 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster4or[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Update Cluster No Update - Correct 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster4or[1])+"\n", makeProxyRequest(proxyPort, ""))

	// then - send request to previous cluster if they have previous uuid
	assertion.AssertDeepEqual("Update Cluster No Update - Correct 3rd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster4or[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))
	assertion.AssertDeepEqual("Update Cluster No Update - Correct 4th response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster4or[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Update Cluster No Update - Correct 5th response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster4or[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion2_0))
	assertion.AssertDeepEqual("Update Cluster No Update - Correct 6th response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster4or[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion2_0))



	// given - another new cluster - instant update
	networkutil.Test_server(serverPortsCluster5iv)
	uuidCookieVersion2_5, _ := networkutil.PUTRequest("http://127.0.0.1:"+strconv.Itoa(configPort)+"/server", "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster5iv[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster5iv[1])+"}], \"upgradeTransition\":{\"mode\":\"INSTANT\"}, \"version\": 2.5}}")

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Update Instance Update - Correct 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster5iv[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Update Instance Update - Correct 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster5iv[1])+"\n", makeProxyRequest(proxyPort, ""))

	// then - send request to previous cluster if they have previous uuid
	assertion.AssertDeepEqual("Update Instance Update - Correct 3rd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster5iv[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))
	assertion.AssertDeepEqual("Update Instance Update - Correct 4th response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster5iv[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Update Instance Update - Correct 5th response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster5iv[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion2_5))
	assertion.AssertDeepEqual("Update Instance Update - Correct 6th response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster5iv[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion2_5))
}



