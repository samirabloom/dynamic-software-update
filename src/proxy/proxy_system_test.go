package proxy

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

func writeConfigFile(proxyPort int, configPort int, uuid string, serverPorts []int, version string) string {
	fileName := "/tmp/system_test_config.json"
	data := "{\"proxy\":{\"ip\":\"localhost\",\"port\":" + strconv.Itoa(proxyPort) + "},\"configService\":{\"port\":" + strconv.Itoa(configPort) + "},\"cluster\":{\"servers\":["
	for index, serverPort := range serverPorts {
		if index > 0 {
			data += ","
		}
		data += "{\"ip\":\"127.0.0.1\",\"port\":"+strconv.Itoa(serverPort)+"}"
	}
	data += "]"
	if len(uuid) > 0 {
		data += ",\"uuid\":\""+uuid+"\""
	}
	if len(version) > 0 {
		data += ", \"version\": "+version
	}
	data += "}}"

	ioutil.WriteFile(fileName, []byte(data), 0644)
	return fileName
}

func makeProxyRequest(proxyPort int, uuidCookie string) string {
	body, _ := networkutil.GETCookiedRequest("http://127.0.0.1:"+strconv.Itoa(proxyPort), uuidCookie)
	return body
}

func Test_proxy_system_test_load_balancing_with_initial_config_file(testCtx *testing.T) {
	var (
		uuidCookieVersion0_0 string = "1027596f-1034-11e4-8334-600308a82410"
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
	)

	// given
	NewProxy(writeConfigFile(proxyPort, configPort, uuidCookieVersion0_0, serverPortsClusterOne, "")).Start(false)
	networkutil.Test_server(serverPortsClusterOne)

	// then - should load balance requests
	assertion.AssertDeepEqual("Initial Config - Correct 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Config - Correct 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Config - Correct 3rd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Config - Correct 4th response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, ""))
}

func Test_proxy_system_test_should_load_balance_with_session_cluster_transition(testCtx *testing.T) {
	var (
		uuidCookieVersion0_0 string = "1027596f-1034-11e4-8334-600308a82410"
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsCluster3re []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/server"
	)

	// given
	NewProxy(writeConfigFile(proxyPort, configPort, uuidCookieVersion0_0, serverPortsClusterOne, "")).Start(false)
	networkutil.Test_server(serverPortsClusterOne)

	// given - new cluster
	networkutil.Test_server(serverPortsClusterTwo)
	uuidCookieVersion1_1, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.1}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, ""))

	// then - should load balance requests against new cluster
	assertion.AssertDeepEqual("Latest Cluster When Matching UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))
	assertion.AssertDeepEqual("Latest Cluster When Matching UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))

	// then - should load balance previous clusters
	assertion.AssertDeepEqual("Initial Cluster When Matching UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion0_0))
	assertion.AssertDeepEqual("Initial Cluster When Matching UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion0_0))


	// given - another new cluster
	networkutil.Test_server(serverPortsCluster3re)
	uuidCookieVersion1_5, _ := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.5}}")

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, ""))

	// then - send request to previous cluster if they have previous uuid
	assertion.AssertDeepEqual("Previous Cluster When Matching UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))
	assertion.AssertDeepEqual("Previous Cluster When Matching UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When Matching UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_5))
	assertion.AssertDeepEqual("Latest Cluster When Matching UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_5))

	// then - should load balance previous clusters
	assertion.AssertDeepEqual("Initial Cluster When Matching UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion0_0))
	assertion.AssertDeepEqual("Initial Cluster When Matching UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion0_0))
}

func Test_proxy_system_test_should_load_balance_with_instant_cluster_transition(testCtx *testing.T) {
	var (
		uuidCookieVersion0_0 string = "1027596f-1034-11e4-8334-600308a82410"
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsCluster3re []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/server"
	)

	// given
	NewProxy(writeConfigFile(proxyPort, configPort, uuidCookieVersion0_0, serverPortsClusterOne, "")).Start(false)
	networkutil.Test_server(serverPortsClusterOne)

	// then - should load balance requests
	assertion.AssertDeepEqual("Initial Config - Correct 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Config - Correct 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Config - Correct 3rd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Config - Correct 4th response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, ""))


	// given - another new cluster - no upgrade defined
	networkutil.Test_server(serverPortsClusterTwo)
	uuidCookieVersion1_1, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.1}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// given - another new cluster - instant update
	networkutil.Test_server(serverPortsCluster3re)
	uuidCookieVersion1_5, _ := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"mode\":\"INSTANT\"}, \"version\": 1.5}}")

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, ""))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When Initial Cluster UUID - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion0_0))
	assertion.AssertDeepEqual("Latest Cluster When Initial Cluster UUID - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion0_0))

	// then - send request to previous cluster if they have previous uuid
	assertion.AssertDeepEqual("Latest Cluster When Updated Cluster UUID - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))
	assertion.AssertDeepEqual("Latest Cluster When Updated Cluster UUID - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_1))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When Matching Cluster UUID - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_5))
	assertion.AssertDeepEqual("Latest Cluster When Matching Cluster UUID - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, uuidCookieVersion1_5))
}

func Test_proxy_system_test_should_update_latest_cluster_with_cluster_removed(testCtx *testing.T) {
	var (
		uuidCookieVersion0_0 string = "1027596f-1034-11e4-8334-600308a82410"
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsCluster3re []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/server"
	)

	// given
	NewProxy(writeConfigFile(proxyPort, configPort, uuidCookieVersion0_0, serverPortsClusterOne, "")).Start(false)
	networkutil.Test_server(serverPortsClusterOne)

	// given - new cluster
	networkutil.Test_server(serverPortsClusterTwo)
	uuidCookieVersion1_1, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.1}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, ""))


	// given - another new cluster
	networkutil.Test_server(serverPortsCluster3re)
	uuidCookieVersion1_5, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.5}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, ""))


	// given - cluster removed
	networkutil.Test_server(serverPortsCluster3re)
	_, deleteStatus := networkutil.DELETERequest(configServiceUrl + "/" + uuidCookieVersion1_5)

	// then - should remove cluster
	assertion.AssertDeepEqual("Remove Cluster - Correct Delete Status", testCtx, "202 Accepted", deleteStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Previous Cluster After Cluster Removed - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Previous Cluster After Cluster Removed - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, ""))


	// given - cluster removed
	networkutil.Test_server(serverPortsCluster3re)
	_, deleteStatus = networkutil.DELETERequest(configServiceUrl+"/"+uuidCookieVersion1_1)

	// then - should remove cluster
	assertion.AssertDeepEqual("Remove Cluster - Correct Delete Status", testCtx, "202 Accepted", deleteStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Initial Cluster After Two Clusters Removed - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Initial Cluster After Two Clusters Removed - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, ""))
}

func Test_proxy_system_test_should_maintain_version_order_with_multiple_clusters(testCtx *testing.T) {
	var (
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsCluster3re []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/server"
	)

	// given
	NewProxy(writeConfigFile(proxyPort, configPort, "", serverPortsClusterOne, "1.0")).Start(false)
	networkutil.Test_server(serverPortsClusterOne)

	// given - new cluster - version less then initial cluster i.e. version 0.9
	networkutil.Test_server(serverPortsClusterTwo)
	_, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 0.9}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, ""))


	// given - new cluster - version same as initial cluster i.e. version 1.0
	networkutil.Test_server(serverPortsCluster3re)
	_, putStatus = networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"ip\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": 1.0}}")

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, ""))

}


