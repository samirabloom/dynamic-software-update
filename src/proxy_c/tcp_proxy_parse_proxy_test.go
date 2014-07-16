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
		expectedTcpProxyLocalAddress, _ = net.ResolveTCPAddr("tcp", "localhost:1234")
	)

	// when
	tcpProxyLocalAddress, err := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, nil)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, tcpProxyLocalAddress, expectedTcpProxyLocalAddress)
}

func Test_Parse_Proxy_When_No_IP(testCtx *testing.T) {
	// given
	var (
		mockProxyConfig = map[string]interface{}{"port":   "1234"}
		jsonConfig      = map[string]interface{}{"proxy": mockProxyConfig}
		error           = errors.New("Invalid proxy configuration please provide \"proxy\" with an \"ip\" and \"port\" in configuration")
	)
	// when
	tcpProxyLocalAddress, err := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, error)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, tcpProxyLocalAddress, nil)
}

func Test_Parse_Proxy_When_IP_Invalid(testCtx *testing.T) {
	// given
	var (
		mockProxyConfig = map[string]interface{}{"ip": "", "port": "1234"}
		jsonConfig      = map[string]interface{}{"proxy": mockProxyConfig}
		error           = errors.New("Invalid proxy configuration please provide \"proxy\" with an \"ip\" and \"port\" in configuration")
	)
	// when
	tcpProxyLocalAddress, err := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, error)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, tcpProxyLocalAddress, nil)
}

func Test_Parse_Proxy_When_No_Port(testCtx *testing.T) {
	// given
	var (
		mockProxyConfig = map[string]interface{}{"ip": "localhost"}
		jsonConfig      = map[string]interface{}{"proxy": mockProxyConfig}
		error           = errors.New("Invalid proxy configuration please provide \"proxy\" with an \"ip\" and \"port\" in configuration")
	)
	// when
	tcpProxyLocalAddress, err := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, error)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, tcpProxyLocalAddress, nil)
}

func Test_Parse_Proxy_When_Port_Invalid(testCtx *testing.T) {
	// given
	var (
		mockProxyConfig = map[string]interface{}{
		"ip": "localhost",
		"port":   "not valid port",
	}
		jsonConfig      = map[string]interface{}{
		"proxy": mockProxyConfig,
	}
		error           = errors.New("Invalid proxy configuration please provide \"proxy\" with an \"ip\" and \"port\" in configuration")
	)
	// when
	tcpProxyLocalAddress, err := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, error)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, tcpProxyLocalAddress, nil)
}

func Test_Parse_Proxy_When_No_IP_Or_Port(testCtx *testing.T) {
	// given
	var (
		jsonConfig = map[string]interface{}{
		"proxy": nil,
	}
		error      = errors.New("Invalid proxy configuration please provide \"proxy\" with an \"ip\" and \"port\" in configuration")
	)
	// when
	tcpProxyLocalAddress, err := parseProxy(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, error)
	assertion.AssertDeepEqual("Correct tcpProxy Local Address", testCtx, tcpProxyLocalAddress, nil)
}



