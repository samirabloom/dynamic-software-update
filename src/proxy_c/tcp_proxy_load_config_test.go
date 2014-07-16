package proxy_c

import (
	"testing"
	"net"
	assertion "util/test/assertion"
)

func Test_Read_Config_When_File_Valid_Range(testCtx *testing.T) {
	// given
	var (
		fileName                        = new(string)
		expectedBackendBaseAddr, _      = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		expectedRouter                  = &RangeRoutingContext{backendBaseAddr: expectedBackendBaseAddr, clusterSize: 8}
		expectedTcpProxyLocalAddress, _ = net.ResolveTCPAddr("tcp", "localhost:1234")
		expectedLoadBalancer            = &LoadBalancer{frontendAddr: expectedTcpProxyLocalAddress, router: expectedRouter, stop: make(chan bool)}
		expectedError error             = nil
	)
	*fileName = "test_range_config.json"

	// when
	actualLoadBalancer, actualError := loadConfig(fileName)

	// then
	assertion.AssertDeepEqual("Correct Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct Load Balancer", testCtx, actualLoadBalancer, expectedLoadBalancer)
}

func Test_Read_Config_When_File_Valid_Server_list(testCtx *testing.T) {
	// given
	var (
		fileName = "does_not_exist.json"
	)

	// when
	println(fileName)

	// then
}


