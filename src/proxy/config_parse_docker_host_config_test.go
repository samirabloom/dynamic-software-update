package proxy

import (
	"testing"
	assertion "util/test/assertion"
	"errors"
	"proxy/docker_client"
)

func Test_Parse_Docker_Host_When_Config_Valid(testCtx *testing.T) {
	// given
	var (
		expectedError error                          = nil
		expectedDockerHost *docker_client.DockerHost = &docker_client.DockerHost{Ip: "123.456.789.012", Port: 1234, Log: true}
		jsonConfig                                   = map[string]interface{}{"dockerHost": map[string]interface{}{"ip": expectedDockerHost.Ip, "port": float64(expectedDockerHost.Port), "log": true}}
	)

	// when
	actualDockerHost, actualErr := parseDockerHost(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Docker Host Port", testCtx, expectedDockerHost, actualDockerHost)
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
}

func Test_Parse_Docker_Host_When_Config_Valid_Default_Port(testCtx *testing.T) {
	// given
	var (
		expectedError error                          = nil
		expectedDockerHost *docker_client.DockerHost = &docker_client.DockerHost{Ip: "123.456.789.012", Port: 2375, Log: true}
		jsonConfig                                   = map[string]interface{}{"dockerHost": map[string]interface{}{"ip": expectedDockerHost.Ip}}
	)

	// when
	actualDockerHost, actualErr := parseDockerHost(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Docker Host Port", testCtx, expectedDockerHost, actualDockerHost)
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
}

func Test_Parse_Docker_Host_When_No_Ip(testCtx *testing.T) {
	// given
	var (
		expectedError error                          = errors.New("Invalid docker host configuration - \"ip\" is missing from \"dockerHost\" config")
		expectedDockerHost *docker_client.DockerHost = nil
		jsonConfig                                   = map[string]interface{}{"dockerHost": map[string]interface{}{"ip": nil}}
	)

	// when
	actualDockerHost, actualErr := parseDockerHost(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Docker Host Port", testCtx, expectedDockerHost, actualDockerHost)
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
}

func Test_Parse_Docker_Host_When_No_DockerHost(testCtx *testing.T) {
	// given
	var (
		expectedError error                          = nil
		expectedDockerHost *docker_client.DockerHost = nil
		jsonConfig                                   = map[string]interface{}{"dockerHost": nil}
	)

	// when
	actualDockerHost, actualErr := parseDockerHost(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Docker Host Port", testCtx, expectedDockerHost, actualDockerHost)
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
}
