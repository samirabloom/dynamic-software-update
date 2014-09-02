package proxy

import (
	"testing"
	assertion "util/test/assertion"
	"net"
	"bytes"
	"proxy/contexts"
	"errors"
	"io"
	"proxy/docker_client"
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

func mockParseDockerHost(dockerHost *docker_client.DockerHost, parseDockerHostErr error) func(map[string]interface{}) (*docker_client.DockerHost, error) {
	return func(map[string]interface{}) (*docker_client.DockerHost, error) {
		return dockerHost, parseDockerHostErr
	}
}

func mockParseClusters(clusters *contexts.Clusters, clusterParseErr error) func(map[string]interface{}, *docker_client.DockerHost, io.Writer) (*contexts.Clusters, error) {
	return func(map[string]interface{}, *docker_client.DockerHost, io.Writer) (*contexts.Clusters, error) {
		return clusters, clusterParseErr
	}
}

func Test_Parse_Config_File_With_No_Errors(testCtx *testing.T) {
	// given
	var (
		tcpProxyLocalAddress *net.TCPAddr    = &net.TCPAddr{}
		proxyParseErr error                  = nil
		configServicePort int                = 4321
		parseConfigServiceErr error          = nil
		dockerHost *docker_client.DockerHost = &docker_client.DockerHost{}
		parseDockerHostErr error             = nil
		clusters *contexts.Clusters          = &contexts.Clusters{}
		clusterParseErr error                = nil
		jsonData                             = []byte("{}")
		expectedProxy *Proxy                 = &Proxy{frontendAddr: tcpProxyLocalAddress, configServicePort: configServicePort, dockerHost: dockerHost, clusters: clusters, stop: make(chan bool), }
		expectedError error                  = nil
		outputStream bytes.Buffer
	)

	// when
	actualProxy, actualError := parseConfigFile(jsonData, mockParseProxy(tcpProxyLocalAddress, proxyParseErr), mockParseConfigService(configServicePort, parseConfigServiceErr), mockParseDockerHost(dockerHost, parseDockerHostErr), mockParseClusters(clusters, clusterParseErr), &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Proxy", testCtx, expectedProxy, actualProxy)
}

func Test_Parse_Config_File_With_Proxy_Parse_Error(testCtx *testing.T) {
	// given
	var (
		tcpProxyLocalAddress *net.TCPAddr    = &net.TCPAddr{}
		proxyParseErr error                  = errors.New("Test Proxy Parse Error")
		configServicePort int                = 4321
		parseConfigServiceErr error          = nil
		dockerHost *docker_client.DockerHost = &docker_client.DockerHost{}
		parseDockerHostErr error             = nil
		clusters *contexts.Clusters          = &contexts.Clusters{}
		clusterParseErr error                = nil
		jsonData                             = []byte("{}")
		expectedProxy *Proxy                 = nil
		expectedError error                  = proxyParseErr
		outputStream bytes.Buffer
	)

	// when
	actualProxy, actualError := parseConfigFile(jsonData, mockParseProxy(tcpProxyLocalAddress, proxyParseErr), mockParseConfigService(configServicePort, parseConfigServiceErr), mockParseDockerHost(dockerHost, parseDockerHostErr), mockParseClusters(clusters, clusterParseErr), &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Proxy", testCtx, expectedProxy, actualProxy)
}

func Test_Parse_Config_File_With_Config_Service_Parse_Error(testCtx *testing.T) {
	// given
	var (
		tcpProxyLocalAddress *net.TCPAddr    = &net.TCPAddr{}
		proxyParseErr error                  = nil
		configServicePort int                = 4321
		parseConfigServiceErr error          = errors.New("Test Config Service Parse Error")
		dockerHost *docker_client.DockerHost = &docker_client.DockerHost{}
		parseDockerHostErr error             = nil
		clusters *contexts.Clusters          = &contexts.Clusters{}
		clusterParseErr error                = nil
		jsonData                             = []byte("{}")
		expectedProxy *Proxy                 = nil
		expectedError error                  = parseConfigServiceErr
		outputStream bytes.Buffer
	)

	// when
	actualProxy, actualError := parseConfigFile(jsonData, mockParseProxy(tcpProxyLocalAddress, proxyParseErr), mockParseConfigService(configServicePort, parseConfigServiceErr), mockParseDockerHost(dockerHost, parseDockerHostErr), mockParseClusters(clusters, clusterParseErr), &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Proxy", testCtx, expectedProxy, actualProxy)
}

func Test_Parse_Config_File_With_Docker_Host_Parse_Error(testCtx *testing.T) {
	// given
	var (
		tcpProxyLocalAddress *net.TCPAddr    = &net.TCPAddr{}
		proxyParseErr error                  = nil
		configServicePort int                = 4321
		parseConfigServiceErr error          = nil
		dockerHost *docker_client.DockerHost = &docker_client.DockerHost{}
		parseDockerHostErr error             = errors.New("Test Docker Host Parse Error")
		clusters *contexts.Clusters          = &contexts.Clusters{}
		clusterParseErr error                = nil
		jsonData                             = []byte("{}")
		expectedProxy *Proxy                 = nil
		expectedError error                  = parseDockerHostErr
		outputStream bytes.Buffer
	)

	// when
	actualProxy, actualError := parseConfigFile(jsonData, mockParseProxy(tcpProxyLocalAddress, proxyParseErr), mockParseConfigService(configServicePort, parseConfigServiceErr), mockParseDockerHost(dockerHost, parseDockerHostErr), mockParseClusters(clusters, clusterParseErr), &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Proxy", testCtx, expectedProxy, actualProxy)
}

func Test_Parse_Config_File_With_Cluster_Parse_Error(testCtx *testing.T) {
	// given
	var (
		tcpProxyLocalAddress *net.TCPAddr    = &net.TCPAddr{}
		proxyParseErr error                  = nil
		configServicePort int                = 4321
		parseConfigServiceErr error          = nil
		dockerHost *docker_client.DockerHost = &docker_client.DockerHost{}
		parseDockerHostErr error             = nil
		clusters *contexts.Clusters          = &contexts.Clusters{}
		clusterParseErr error                = errors.New("Test Config Service Parse Error")
		jsonData                             = []byte("{}")
		expectedProxy *Proxy                 = nil
		expectedError error                  = clusterParseErr
		outputStream bytes.Buffer
	)

	// when
	actualProxy, actualError := parseConfigFile(jsonData, mockParseProxy(tcpProxyLocalAddress, proxyParseErr), mockParseConfigService(configServicePort, parseConfigServiceErr), mockParseDockerHost(dockerHost, parseDockerHostErr), mockParseClusters(clusters, clusterParseErr), &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Proxy", testCtx, expectedProxy, actualProxy)
}



