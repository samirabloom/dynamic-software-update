package proxy_c

import (
	"testing"
	assertion "util/test/assertion"
	"net"
	"errors"
	"fmt"
)

func Test_Parse_Cluster_Config_When_Config_Valid(testCtx *testing.T) {
	// given
	var (
		expectedError error = nil
		backendBaseAddr, _  = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		expectedMockRouter  = &RangeRoutingContext{backendBaseAddr: backendBaseAddr, clusterSize: 8}
		proxyConfig     = map[string]interface{}{"ip": "127.0.0.1", "port": 1024, "clusterSize": "8"}
		jsonConfig          = map[string]interface{}{"server_range": proxyConfig}
	)

	// when
	actualMockRouter, actualError := parseClusterConfig(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, actualMockRouter, expectedMockRouter)
}



func Test_Parse_Cluster_When_No_server_range_IP(testCtx *testing.T) {
	// given
	var (
		proxyConfig     = map[string]interface{}{"port": 1024, "clusterSize": "8"}
		jsonConfig          = map[string]interface{}{"server_range": proxyConfig}
		expectedError      = errors.New("Invalid server range configuration please provide \"server_range\" with an \"ip\" and \"port\" in configuration - address provided was [" + fmt.Sprintf("%s:%v", proxyConfig["ip"], proxyConfig["port"]) + "]")
	)

	// when
	actualMockRouter, actualError := parseClusterConfig(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, actualMockRouter, nil)
}

func TestParse_Cluster_When_IP_Invalid(testCtx *testing.T) {
	// given
	var (
		proxyConfig     = map[string]interface{}{"ip": "inv@lid", "port": 1024, "clusterSize": "8"}
		jsonConfig          = map[string]interface{}{"server_range": proxyConfig}
		expectedError      = errors.New("Invalid server range configuration please provide \"server_range\" with an \"ip\" and \"port\" in configuration - address provided was [" + fmt.Sprintf("%s:%v", proxyConfig["ip"], proxyConfig["port"]) + "]")
	)

	// when
	actualMockRouter, actualError := parseClusterConfig(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, actualMockRouter, nil)
}

func Test_Parse_Cluster_When_No_Port(testCtx *testing.T) {
	// given
	var (
		proxyConfig     = map[string]interface{}{"ip": "127.0.0.1", "clusterSize": "8"}
		jsonConfig          = map[string]interface{}{"server_range": proxyConfig}
		expectedError      = errors.New("Invalid server range configuration please provide \"server_range\" with an \"ip\" and \"port\" in configuration - address provided was [" + fmt.Sprintf("%s:%v", proxyConfig["ip"], proxyConfig["port"]) + "]")
	)

	// when
	actualMockRouter, actualError := parseClusterConfig(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, actualMockRouter, nil)
}

func Test_Parse_Cluster_When_Port_Invalid(testCtx *testing.T) {
	// given
	var (
		proxyConfig     = map[string]interface{}{"ip": "127.0.0.1", "port": "invalid", "clusterSize": "8"}
		jsonConfig          = map[string]interface{}{"server_range": proxyConfig}
		expectedError      = errors.New("Invalid server range configuration please provide \"server_range\" with an \"ip\" and \"port\" in configuration - address provided was [" + fmt.Sprintf("%s:%v", proxyConfig["ip"], proxyConfig["port"]) + "]")
	)

	// when
	actualMockRouter, actualError := parseClusterConfig(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, actualMockRouter, nil)
}


func Test_Parse_Cluster_Config_When_ClusterSize_Invalid(testCtx *testing.T) {
	// given
	var (
		proxyConfig     = map[string]interface{}{"ip": "127.0.0.1", "port": 1024, "clusterSize": "invalid"}
		jsonConfig          = map[string]interface{}{"server_range": proxyConfig}
		expectedError      = errors.New("Cluster Size not a valid integer [" + proxyConfig["clusterSize"].(string) + "]")
	)

	// when
	actualMockRouter, actualError := parseClusterConfig(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, actualMockRouter, nil)
}

func Test_Parse_Cluster_When_Server_Range_Nil(testCtx *testing.T) {
	// given
	var (
		jsonConfig = map[string]interface{}{
		"server_range": nil,
	}
		error      = errors.New("Invalid proxy configuration please provide \"proxy\" with an \"ip\" and \"port\" in configuration")
	)
	// when
	tcpProxyLocalAddress, err := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, error)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, tcpProxyLocalAddress, nil)
}

func Test_Parse_Cluster_When_Servers_Nil(testCtx *testing.T) {
	// given
	var (
		jsonConfig = map[string]interface{}{
		"servers": nil,
	}
		error      = errors.New("Invalid proxy configuration please provide \"proxy\" with an \"ip\" and \"port\" in configuration")
	)
	// when
	tcpProxyLocalAddress, err := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, error)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, tcpProxyLocalAddress, nil)
}

// TODO add more for the range



