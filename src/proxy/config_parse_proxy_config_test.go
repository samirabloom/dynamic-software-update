package proxy

import (
	"testing"
	assertion "util/test/assertion"
	"net"
	"errors"
)

// calling ParseProxy with
// Proxy not nil
// Err nil
func Test_Parse_Proxy_When_Config_Valid(testCtx *testing.T) {
	// given
	var (
		proxyConfig                     = map[string]interface{}{"port":   1234}
		jsonConfig                      = map[string]interface{}{"proxy": proxyConfig}
		expectedError error             = nil
		expectedTcpProxyLocalAddress, _ = net.ResolveTCPAddr("tcp", ":1234")
	)

	// when
	actualTcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, actualTcpProxyLocalAddress)
}

func Test_Parse_Proxy_When_No_Port(testCtx *testing.T) {
	// given
	var (
		proxyConfig                               = map[string]interface{}{}
		jsonConfig                                = map[string]interface{}{"proxy": proxyConfig}
		expectedError                             = errors.New("Invalid proxy address [:<nil>] - unknown port tcp/<nil>")
		expectedTcpProxyLocalAddress *net.TCPAddr = nil
	)
	// when
	actualTcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, actualTcpProxyLocalAddress)
}

func Test_Parse_Proxy_When_Port_Invalid(testCtx *testing.T) {
	// given
	var (
		proxyConfig                               = map[string]interface{}{"port":   "not valid port"}
		jsonConfig                                = map[string]interface{}{"proxy": proxyConfig}
		expectedError                             = errors.New("Invalid proxy address [:not valid port] - unknown port tcp/not valid port")
		expectedTcpProxyLocalAddress *net.TCPAddr = nil
	)
	// when
	actualTcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, actualTcpProxyLocalAddress)
}

func Test_Parse_Proxy_When_No_IP_Or_Port(testCtx *testing.T) {
	// given
	var (
		jsonConfig                                = map[string]interface{}{"proxy": nil}
		expectedError                             = errors.New("Invalid proxy configuration - \"proxy\" config missing")
		expectedTcpProxyLocalAddress *net.TCPAddr = nil
	)
	// when
	actualTcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, actualTcpProxyLocalAddress)
}





