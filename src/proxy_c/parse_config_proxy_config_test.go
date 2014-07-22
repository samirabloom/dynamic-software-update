package proxy_c

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
		proxyConfig                     = map[string]interface{}{"ip": "localhost", "port":   1234}
		jsonConfig                      = map[string]interface{}{"proxy": proxyConfig}
		expectedError error             = nil
		expectedTcpProxyLocalAddress, _ = net.ResolveTCPAddr("tcp", "localhost:1234")
	)

	// when
	tcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, tcpProxyLocalAddress)
}

func Test_Parse_Proxy_When_No_IP(testCtx *testing.T) {
	// given
	var (
		proxyConfig                               = map[string]interface{}{"port":   "1234"}
		jsonConfig                                = map[string]interface{}{"proxy": proxyConfig}
		expectedError                             = errors.New("Invalid proxy address [%!s(<nil>):1234] - missing brackets in address %!s(<nil>):1234")
		expectedTcpProxyLocalAddress *net.TCPAddr = nil
	)
	// when
	tcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, tcpProxyLocalAddress)
}

func Test_Parse_Proxy_When_IP_Invalid(testCtx *testing.T) {
	// given
	var (
		proxyConfig                               = map[string]interface{}{"ip": "inv@lid", "port": "1234"}
		jsonConfig                                = map[string]interface{}{"proxy": proxyConfig}
		expectedError                             = errors.New("Invalid proxy address [inv@lid:1234] - lookup inv@lid: no such host")
		expectedTcpProxyLocalAddress *net.TCPAddr = nil
	)
	// when
	tcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, tcpProxyLocalAddress)
}

func Test_Parse_Proxy_When_No_Port(testCtx *testing.T) {
	// given
	var (
		proxyConfig                               = map[string]interface{}{"ip": "localhost"}
		jsonConfig                                = map[string]interface{}{"proxy": proxyConfig}
		expectedError                             = errors.New("Invalid proxy address [localhost:<nil>] - unknown port tcp/<nil>")
		expectedTcpProxyLocalAddress *net.TCPAddr = nil
	)
	// when
	tcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, tcpProxyLocalAddress)
}

func Test_Parse_Proxy_When_Port_Invalid(testCtx *testing.T) {
	// given
	var (
		proxyConfig                               = map[string]interface{}{"ip": "localhost", "port":   "not valid port"}
		jsonConfig                                = map[string]interface{}{"proxy": proxyConfig}
		expectedError                             = errors.New("Invalid proxy address [localhost:not valid port] - unknown port tcp/not valid port")
		expectedTcpProxyLocalAddress *net.TCPAddr = nil
	)
	// when
	tcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, tcpProxyLocalAddress)
}

func Test_Parse_Proxy_When_No_IP_Or_Port(testCtx *testing.T) {
	// given
	var (
		jsonConfig                                = map[string]interface{}{"proxy": nil}
		expectedError                             = errors.New("Invalid proxy configuration - \"proxy\" config missing")
		expectedTcpProxyLocalAddress *net.TCPAddr = nil
	)
	// when
	tcpProxyLocalAddress, actualErr := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
	assertion.AssertDeepEqual("Correct Local Proxy Address", testCtx, expectedTcpProxyLocalAddress, tcpProxyLocalAddress)
}





