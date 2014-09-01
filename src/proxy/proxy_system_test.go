package proxy

import (
	"testing"
	networkutil "util/test/network"
	assertion "util/test/assertion"
	"strconv"
	"io/ioutil"
	"util/test/vagrant"
	"github.com/franela/goreq"
	"time"
	"os"
)

func WriteConfigSystemTestConfigFile(proxyPort int, configPort int, uuid string, serverPorts []int, version string) string {
	fileName := "/tmp/system_test_config.json"
	data := "{\"proxy\":{\"hostname\":\"localhost\",\"port\":" + strconv.Itoa(proxyPort) + "},\"configService\":{\"port\":" + strconv.Itoa(configPort) + "},\"cluster\":{\"servers\":["
	for index, serverPort := range serverPorts {
		if index > 0 {
			data += ","
		}
		data += "{\"hostname\":\"127.0.0.1\",\"port\":"+strconv.Itoa(serverPort)+"}"
	}
	data += "]"
	if len(uuid) > 0 {
		data += ",\"uuid\":\""+uuid+"\""
	}
	if len(version) > 0 {
		data += ",\"version\":\""+version+"\""
	}
	data += "}}"

	ioutil.WriteFile(fileName, []byte(data), 0644)
	return fileName
}

func WriteConfigSystemTestWordPressDockerConfigFile(proxyPort int, configPort int, dockerHost string, dockerPort int, uuid string, serverPorts []int, version string) string {
	fileName := "/tmp/system_test_docker_config.json"
	data := "{\n" +
			"    \"proxy\": {\n" +
			"        \"port\": " + strconv.Itoa(proxyPort) + "\n" +
			"    },\n" +
			"    \"configService\": {\n" +
			"        \"port\": " + strconv.Itoa(configPort) + "\n" +
			"    },\n" +
			"    \"dockerHost\": {\n" +
			"        \"ip\": \"" + dockerHost + "\",\n" +
			"        \"port\": " + strconv.Itoa(dockerPort) + ",\n" +
			"        \"log\": true\n" +
			"    },\n" +
			"    \"cluster\": {\n" +
			"        \"containers\": [\n";
	for index, serverPort := range serverPorts {
		if index > 0 {
			data += ","
		}
		data += "         {\n" +
				"             \"image\": \"mysql\",\n" +
				"             \"tag\": \"latest\",\n" +
				"             \"alwaysPull\": true,\n" +
				"             \"name\": \"some-mysql-test\",\n" +
				"             \"environment\": [\n" +
				"                 \"MYSQL_ROOT_PASSWORD=mysecretpassword\"\n" +
				"             ],\n" +
				"             \"volumes\": [\n" +
				"                 \"/var/lib/mysql:/var/lib/mysql\"\n" +
				"             ]\n" +
				"         },\n" +
				"         {\n" +
				"             \"image\": \"wordpress\",\n" +
				"             \"tag\": \"3.9.2\",\n" +
				"             \"alwaysPull\": true,\n" +
				"             \"portToProxy\": "+strconv.Itoa(serverPort)+",\n" +
				"             \"name\": \"some-wordpress-test\",\n" +
				"             \"links\": [\n" +
				"                 \"some-mysql-test:mysql\"\n" +
				"             ],\n" +
				"             \"portBindings\": {\n" +
				"                 \"80/tcp\": [\n" +
				"                     {\n" +
				"                         \"HostIp\": \"0.0.0.0\",\n" +
				"                         \"HostPort\": \""+strconv.Itoa(serverPort)+"\"\n" +
				"                     }\n" +
				"                 ]\n" +
				"             }\n" +
				"         }\n";
	}
	data += "        ]\n";
	if len(uuid) > 0 {
		data += ",   \"uuid\":\""+uuid+"\""
	}
	if len(version) > 0 {
		data += ",   \"version\":\""+version+"\""
	}
	data += "    }\n" +
			"}\n"
	ioutil.WriteFile(fileName, []byte(data), 0644)
	return fileName
}

func WriteConfigSystemTestMockServerDockerConfigFile(proxyPort int, configPort int, dockerHost string, dockerPort int, uuid string, serverPorts []int, version string) string {
	fileName := "/tmp/system_test_docker_config.json"
	data := "{\n" +
			"    \"proxy\": {\n" +
			"        \"port\": " + strconv.Itoa(proxyPort) + "\n" +
			"    },\n" +
			"    \"configService\": {\n" +
			"        \"port\": " + strconv.Itoa(configPort) + "\n" +
			"    },\n" +
			"    \"dockerHost\": {\n" +
			"        \"ip\": \"" + dockerHost + "\",\n" +
			"        \"port\": " + strconv.Itoa(dockerPort) + ",\n" +
			"        \"log\": true\n" +
			"    },\n" +
			"    \"cluster\": {\n" +
			"        \"containers\": [\n";
	for index, serverPort := range serverPorts {
		if index > 0 {
			data += ","
		}
		data += "         {\n" +
				"             \"image\": \"jamesdbloom/mockserver\",\n" +
				"             \"portToProxy\": "+strconv.Itoa(serverPort)+",\n" +
				"             \"portBindings\": {\n" +
				"                 \"8080/tcp\": [\n" +
				"                     {\n" +
				"                         \"HostIp\": \"0.0.0.0\",\n" +
				"                         \"HostPort\": \""+strconv.Itoa(serverPort)+"\"\n" +
				"                     }\n" +
				"                 ]\n" +
				"             }\n" +
				"         }\n";
	}
	data += "        ]\n";
	if len(uuid) > 0 {
		data += ",   \"uuid\":\""+uuid+"\""
	}
	if len(version) > 0 {
		data += ",   \"version\":\""+version+"\""
	}
	data += "    }\n" +
			"}\n"
	ioutil.WriteFile(fileName, []byte(data), 0644)
	return fileName
}

func makeProxyRequest(proxyPort int, path string, sessionUuidCookie string, gradualTransitionUuidCookie string) string {
	var (
		sessionUuidCookieHeader *networkutil.Header
		gradualTransitionUuidCookieHeader *networkutil.Header
	)
	if len(sessionUuidCookie) > 0 {
		sessionUuidCookieHeader = &networkutil.Header{"Cookie", "dynsoftup="+sessionUuidCookie+";"}
	}
	if len(gradualTransitionUuidCookie) > 0 {
		gradualTransitionUuidCookieHeader = &networkutil.Header{"Cookie", "dynsoftup="+gradualTransitionUuidCookie+";"}
	}
	body, _ := networkutil.GETRequestWithHeaders("http://127.0.0.1:"+strconv.Itoa(proxyPort)+path, sessionUuidCookieHeader, gradualTransitionUuidCookieHeader)
	return body
}

func Test_Proxy_System_Test_Load_Balancing_With_Initial_Config_File(testCtx *testing.T) {
	var (
		uuidCookieVersion0_0 string = "1027596f-1034-11e4-8334-600308a82410"
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
	)

	// given
	proxy, err := LoadConfig(WriteConfigSystemTestConfigFile(proxyPort, configPort, uuidCookieVersion0_0, serverPortsClusterOne, ""), os.Stdout)
	if err != nil {
		testCtx.Fatalf("Failed to proxy - err: %s\n", err)
	}
	proxy.Start(false)
	networkutil.Test_server(serverPortsClusterOne, false)

	// then - should load balance requests
	assertion.AssertDeepEqual("Initial Config - Correct 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Initial Config - Correct 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Initial Config - Correct 3rd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Initial Config - Correct 4th response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))
}

func Test_Proxy_System_Test_Should_Load_Balance_With_Session_Cluster_Transition(testCtx *testing.T) {
	var (
		uuidCookieVersion0_0 string = "1027596f-1034-11e4-8334-600308a82410"
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsCluster3re []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/configuration/cluster"
	)

	// given
	proxy, err := LoadConfig(WriteConfigSystemTestConfigFile(proxyPort, configPort, uuidCookieVersion0_0, serverPortsClusterOne, ""), os.Stdout)
	if err != nil {
		testCtx.Fatalf("Failed to proxy - err: %s\n", err)
	}
	proxy.Start(false)
	networkutil.Test_server(serverPortsClusterOne, false)

	// given - new cluster
	networkutil.Test_server(serverPortsClusterTwo, false)
	uuidCookieVersion1_1, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"1.1\"}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))

	// then - should load balance requests against new cluster
	assertion.AssertDeepEqual("Latest Cluster When Matching UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_1, ""))
	assertion.AssertDeepEqual("Latest Cluster When Matching UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_1, ""))

	// then - should load balance previous clusters
	assertion.AssertDeepEqual("Initial Cluster When Matching UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion0_0, ""))
	assertion.AssertDeepEqual("Initial Cluster When Matching UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion0_0, ""))


	// given - another new cluster
	networkutil.Test_server(serverPortsCluster3re, false)
	uuidCookieVersion1_5, _ := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"1.5\"}}")

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))

	// then - send request to previous cluster if they have previous uuid
	assertion.AssertDeepEqual("Previous Cluster When Matching UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_1, ""))
	assertion.AssertDeepEqual("Previous Cluster When Matching UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_1, ""))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When Matching UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_5, ""))
	assertion.AssertDeepEqual("Latest Cluster When Matching UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_5, ""))

	// then - should load balance previous clusters
	assertion.AssertDeepEqual("Initial Cluster When Matching UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion0_0, ""))
	assertion.AssertDeepEqual("Initial Cluster When Matching UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion0_0, ""))
}

func Test_Proxy_System_Test_Should_Load_Balance_With_Instant_Cluster_Transition(testCtx *testing.T) {
	var (
		uuidCookieVersion0_0 string = "1027596f-1034-11e4-8334-600308a82410"
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsCluster3re []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/configuration/cluster"
	)

	// given
	proxy, err := LoadConfig(WriteConfigSystemTestConfigFile(proxyPort, configPort, uuidCookieVersion0_0, serverPortsClusterOne, ""), os.Stdout)
	if err != nil {
		testCtx.Fatalf("Failed to proxy - err: %s\n", err)
	}
	proxy.Start(false)
	networkutil.Test_server(serverPortsClusterOne, false)

	// then - should load balance requests
	assertion.AssertDeepEqual("Initial Config - Correct 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Initial Config - Correct 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Initial Config - Correct 3rd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Initial Config - Correct 4th response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))


	// given - another new cluster - no upgrade defined
	networkutil.Test_server(serverPortsClusterTwo, false)
	uuidCookieVersion1_1, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"1.1\"}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// given - another new cluster - instant update
	networkutil.Test_server(serverPortsCluster3re, false)
	uuidCookieVersion1_5, _ := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"mode\":\"INSTANT\"}, \"version\": \"1.5\"}}")

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When Initial Cluster UUID - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion0_0, ""))
	assertion.AssertDeepEqual("Latest Cluster When Initial Cluster UUID - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion0_0, ""))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When Updated Cluster UUID - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_1, ""))
	assertion.AssertDeepEqual("Latest Cluster When Updated Cluster UUID - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_1, ""))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When Matching Cluster UUID - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_5, ""))
	assertion.AssertDeepEqual("Latest Cluster When Matching Cluster UUID - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, "", uuidCookieVersion1_5, ""))
}

func Test_Proxy_System_Test_Should_Update_Latest_Cluster_With_Cluster_Removed(testCtx *testing.T) {
	var (
		uuidCookieVersion0_0 string = "1027596f-1034-11e4-8334-600308a82410"
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsCluster3re []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/configuration/cluster"
	)

	// given
	proxy, err := LoadConfig(WriteConfigSystemTestConfigFile(proxyPort, configPort, uuidCookieVersion0_0, serverPortsClusterOne, ""), os.Stdout)
	if err != nil {
		testCtx.Fatalf("Failed to proxy - err: %s\n", err)
	}
	proxy.Start(false)
	networkutil.Test_server(serverPortsClusterOne, false)

	// given - new cluster
	networkutil.Test_server(serverPortsClusterTwo, false)
	uuidCookieVersion1_1, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"1.1\"}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))


	// given - another new cluster
	networkutil.Test_server(serverPortsCluster3re, false)
	uuidCookieVersion1_5, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"1.5\"}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster Again - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))


	// given - cluster removed
	networkutil.Test_server(serverPortsCluster3re, false)
	_, deleteStatus := networkutil.DELETERequest(configServiceUrl + "/" + uuidCookieVersion1_5)

	// then - should remove cluster
	assertion.AssertDeepEqual("Remove Cluster - Correct Delete Status", testCtx, "202 Accepted", deleteStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Previous Cluster After Cluster Removed - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Previous Cluster After Cluster Removed - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))


	// given - cluster removed
	networkutil.Test_server(serverPortsCluster3re, false)
	_, deleteStatus = networkutil.DELETERequest(configServiceUrl+"/"+uuidCookieVersion1_1)

	// then - should remove cluster
	assertion.AssertDeepEqual("Remove Cluster Again - Correct Delete Status", testCtx, "202 Accepted", deleteStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Initial Cluster After Two Clusters Removed - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Initial Cluster After Two Clusters Removed - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))
}

func Test_Proxy_System_Test_Should_Maintain_Version_Order_With_Multiple_Clusters(testCtx *testing.T) {
	var (
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsCluster3re []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/configuration/cluster"
	)

	// given
	proxy, err := LoadConfig(WriteConfigSystemTestConfigFile(proxyPort, configPort, "", serverPortsClusterOne, "1.0"), os.Stdout)
	if err != nil {
		testCtx.Fatalf("Failed to proxy - err: %s\n", err)
	}
	proxy.Start(false)
	networkutil.Test_server(serverPortsClusterOne, false)

	// given - new cluster - version less then initial cluster i.e. version 0.9
	networkutil.Test_server(serverPortsClusterTwo, false)
	_, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"0.9\"}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After One New Cluster - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))


	// given - new cluster - version same as initial cluster i.e. version 1.0
	networkutil.Test_server(serverPortsCluster3re, false)
	_, putStatus = networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"sessionTimeout\":1}, \"version\": \"1.0\"}}")

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Latest Cluster When No UUID After Two New Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsCluster3re[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))

}

func Test_Proxy_System_Test_Should_Route_Concurrently(testCtx *testing.T) {
	var (
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/configuration/cluster"
	)

	// given
	proxy, err := LoadConfig(WriteConfigSystemTestConfigFile(proxyPort, configPort, "", serverPortsClusterOne, "1.0"), os.Stdout)
	if err != nil {
		testCtx.Fatalf("Failed to proxy - err: %s\n", err)
	}
	proxy.Start(false)
	networkutil.Test_server(serverPortsClusterOne, false)

	// given - new concurrent cluster
	networkutil.Test_server(serverPortsClusterTwo, true)
	_, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"mode\":\"CONCURRENT\"}, \"version\": \"1.1\"}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Latest Cluster When Concurrent Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Latest Cluster When Concurrent Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterTwo[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))

	// then - should load balance requests against latest cluster
	assertion.AssertDeepEqual("Previous Cluster When Concurrent Clusters And Latest Crashes - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "/crash", "", ""))
	assertion.AssertDeepEqual("Previous Cluster When Concurrent Clusters And Latest Crashes - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "/crash", "", ""))

}

func Test_Proxy_System_Test_Should_Fallback_On_Connection_Error(testCtx *testing.T) {
	var (
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsClusterTwo []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		serverPortsCluster3re []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port, networkutil.FindFreeLocalSocket(testCtx).Port}
		configServiceUrl            = "http://127.0.0.1:" + strconv.Itoa(configPort) + "/configuration/cluster"
	)

	// given
	proxy, err := LoadConfig(WriteConfigSystemTestConfigFile(proxyPort, configPort, "", serverPortsClusterOne, "1.0"), os.Stdout)
	if err != nil {
		testCtx.Fatalf("Failed to proxy - err: %s\n", err)
	}
	proxy.Start(false)
	networkutil.Test_server(serverPortsClusterOne, false)

	// given - new instant cluster (that will fail immediately)
	_, putStatus := networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsClusterTwo[1])+"}], \"upgradeTransition\":{\"mode\":\"INSTANT\"}, \"version\": \"1.1\"}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update INSTANT Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// given - new session cluster (that will fail immediately)
	_, putStatus = networkutil.PUTRequest(configServiceUrl, "{\"cluster\": {\"servers\": [{\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[0])+"}, {\"hostname\":\"127.0.0.1\", \"port\":"+strconv.Itoa(serverPortsCluster3re[1])+"}], \"upgradeTransition\":{\"mode\":\"SESSION\", \"sessionTimeout\":1}, \"version\": \"1.2\"}}")

	// then - should update cluster configuration
	assertion.AssertDeepEqual("Update SESSION Cluster - Correct PUT Status", testCtx, "202 Accepted", putStatus)

	// then - should fallback to previous cluster
	assertion.AssertDeepEqual("Latest Cluster When Concurrent Clusters - 1st response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[0])+"\n", makeProxyRequest(proxyPort, "", "", ""))
	assertion.AssertDeepEqual("Latest Cluster When Concurrent Clusters - 2nd response", testCtx, "Port: "+strconv.Itoa(serverPortsClusterOne[1])+"\n", makeProxyRequest(proxyPort, "", "", ""))
}

func Test_Proxy_System_Test_Should_Start_Up_With_WordPress_Docker_Containers(testCtx *testing.T) {
	var (
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port}
	)

	// given
	vagrant.CreateVagrantDockerBox()
	proxy, err := LoadConfig(WriteConfigSystemTestWordPressDockerConfigFile(proxyPort, configPort, "192.168.50.5", 2375, "", serverPortsClusterOne, "1.0"), os.Stdout)
	if err != nil {
		testCtx.Fatalf("Failed to proxy - err: %s\n", err)
	}
	proxy.Start(false)
	time.Sleep(1 * time.Second)

	// then - should proxy requests to wordpress
	res, err := goreq.Request{
		Method: "GET",
		Uri: "http://127.0.0.1:" + strconv.Itoa(proxyPort) + "/",
	}.Do()

	if res == nil {
		testCtx.Fatalf("Response from WordPress cluster is nil - err: %s\n", err)
	}
	assertion.AssertDeepEqual("Wordpress Header Returned", testCtx, []string{"Apache/2.2.22 (Debian)"}, res.Header["Server"])
	assertion.AssertDeepEqual("Wordpress Header Returned", testCtx, []string{"PHP/5.4.4-14+deb7u12"}, res.Header["X-Powered-By"])
}

func Test_Proxy_System_Test_Should_Start_Up_With_MockServer_Docker_Containers(testCtx *testing.T) {
	var (
		proxyPort int               = networkutil.FindFreeLocalSocket(testCtx).Port
		configPort int              = networkutil.FindFreeLocalSocket(testCtx).Port
		serverPortsClusterOne []int = []int{networkutil.FindFreeLocalSocket(testCtx).Port}
	)

	// given
	vagrant.CreateVagrantDockerBox()
	proxy, err := LoadConfig(WriteConfigSystemTestMockServerDockerConfigFile(proxyPort, configPort, "192.168.50.5", 2375, "", serverPortsClusterOne, "1.0"), os.Stdout)
	if err != nil {
		testCtx.Fatalf("Failed to proxy - err: %s\n", err)
	}
	proxy.Start(false)
	time.Sleep(1 * time.Second)

	// then - should setup expectation on MockServer
	res, err := goreq.Request{
		Method: "PUT",
		Uri: "http://127.0.0.1:"+strconv.Itoa(proxyPort)+"/expectation",
		Body: " {\n" +
				"    \"httpRequest\": {\n" +
				"        \"path\": \"/test\"\n" +
				"    }, \n" +
				"    \"httpResponse\": {\n" +
				"        \"body\": \"this is a test response\"\n" +
				"    }, \n" +
				"    \"times\": {\n" +
				"        \"remainingTimes\": 1, \n" +
				"        \"unlimited\": false\n" +
				"    }\n" +
				"}",
	}.Do()
	if res == nil {
		testCtx.Fatalf("Failed to setup expectation on MockServer - err: %s\n", err)
	}

	// then - should proxy requests to MockServer
	res, err = goreq.Request{
		Method: "GET",
		Uri: "http://127.0.0.1:"+strconv.Itoa(proxyPort)+"/test",
	}.Do()

	if res == nil {
		testCtx.Fatalf("Failed to receive response from MockServer - err: %s\n", err)
	}
	var body = make([]byte, 1024)
	size, err := res.Body.Read(body)
	if err != nil {
		testCtx.Fatalf("Failed read response body from MockServer - err: %s\n", err)
	}
	assertion.AssertDeepEqual("MockServer returned correct body", testCtx, "this is a test response", string(body[0:size]))
}


