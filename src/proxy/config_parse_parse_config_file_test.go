package proxy

import (
	"testing"
	assertion "util/test/assertion"
	"net"
	"proxy/stages"
	"errors"
)

func mockParseProxy(tcpProxyLocalAddress *net.TCPAddr, proxyParseErr error) func(map[string]interface{}) (*net.TCPAddr, error) {
	return func(map[string]interface{}) (*net.TCPAddr, error) {
		return tcpProxyLocalAddress, proxyParseErr
	}
}

func mockParseConfigService(configServicePort int, parseConfigServiceErr error) func(map[string]interface{}) (int, error) {
	return func(map[string]interface{}) (int, error) {
		return configServicePort, parseConfigServiceErr
	}
}

func mockParseClusters(clusters *stages.Clusters, clusterParseErr error) func(map[string]interface{}) (*stages.Clusters, error) {
	return func(map[string]interface{}) (*stages.Clusters, error) {
		return clusters, clusterParseErr
	}
}

func Test_Parse_Config_File_With_No_Errors(testCtx *testing.T) {
	// given
	var (
		tcpProxyLocalAddress *net.TCPAddr = &net.TCPAddr{}
		proxyParseErr error               = nil
		configServicePort int             = 4321
		parseConfigServiceErr error       = nil
		clusters *stages.Clusters         = &stages.Clusters{}
		clusterParseErr error             = nil
		jsonData                          = []byte("")
		expectedProxy *Proxy              = &Proxy{frontendAddr: tcpProxyLocalAddress, configServicePort: configServicePort, clusters: clusters, stop: make(chan bool), }
		expectedError error               = nil
	)

	// when
	actualProxy, actualError := parseConfigFile(jsonData, mockParseProxy(tcpProxyLocalAddress, proxyParseErr), mockParseConfigService(configServicePort, parseConfigServiceErr), mockParseClusters(clusters, clusterParseErr))

	// then
	assertion.AssertDeepEqual("Correct Proxy", testCtx, actualProxy, expectedProxy)
	assertion.AssertDeepEqual("Correct Error", testCtx, actualError, expectedError)
}

func Test_Parse_Config_File_With_Proxy_Parse_Error(testCtx *testing.T) {
	// given
	var (
		tcpProxyLocalAddress *net.TCPAddr = &net.TCPAddr{}
		proxyParseErr error               = errors.New("Test Proxy Parse Error")
		configServicePort int             = 4321
		parseConfigServiceErr error       = nil
		clusters *stages.Clusters         = &stages.Clusters{}
		clusterParseErr error             = nil
		jsonData                          = []byte("")
		expectedProxy *Proxy              = nil
		expectedError error               = proxyParseErr
	)

	// when
	actualProxy, actualError := parseConfigFile(jsonData, mockParseProxy(tcpProxyLocalAddress, proxyParseErr), mockParseConfigService(configServicePort, parseConfigServiceErr), mockParseClusters(clusters, clusterParseErr))

	// then
	assertion.AssertDeepEqual("Correct Proxy", testCtx, actualProxy, expectedProxy)
	assertion.AssertDeepEqual("Correct Error", testCtx, actualError, expectedError)
}

func Test_Parse_Config_File_With_Config_Service_Parse_Error(testCtx *testing.T) {
	// given
	var (
		tcpProxyLocalAddress *net.TCPAddr = &net.TCPAddr{}
		proxyParseErr error               = nil
		configServicePort int             = 4321
		parseConfigServiceErr error       = errors.New("Test Config Service Parse Error")
		clusters *stages.Clusters         = &stages.Clusters{}
		clusterParseErr error             = nil
		jsonData                          = []byte("")
		expectedProxy *Proxy              = nil
		expectedError error               = parseConfigServiceErr
	)

	// when
	actualProxy, actualError := parseConfigFile(jsonData, mockParseProxy(tcpProxyLocalAddress, proxyParseErr), mockParseConfigService(configServicePort, parseConfigServiceErr), mockParseClusters(clusters, clusterParseErr))

	// then
	assertion.AssertDeepEqual("Correct Proxy", testCtx, actualProxy, expectedProxy)
	assertion.AssertDeepEqual("Correct Error", testCtx, actualError, expectedError)
}

func Test_Parse_Config_File_With_Cluster_Parse_Error(testCtx *testing.T) {
	// given
	var (
		tcpProxyLocalAddress *net.TCPAddr = &net.TCPAddr{}
		proxyParseErr error               = nil
		configServicePort int             = 4321
		parseConfigServiceErr error       = nil
		clusters *stages.Clusters         = &stages.Clusters{}
		clusterParseErr error             = errors.New("Test Config Service Parse Error")
		jsonData                          = []byte("")
		expectedProxy *Proxy              = nil
		expectedError error               = clusterParseErr
	)

	// when
	actualProxy, actualError := parseConfigFile(jsonData, mockParseProxy(tcpProxyLocalAddress, proxyParseErr), mockParseConfigService(configServicePort, parseConfigServiceErr), mockParseClusters(clusters, clusterParseErr))

	// then
	assertion.AssertDeepEqual("Correct Proxy", testCtx, actualProxy, expectedProxy)
	assertion.AssertDeepEqual("Correct Error", testCtx, actualError, expectedError)
}



