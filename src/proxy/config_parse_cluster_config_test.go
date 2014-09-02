package proxy

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"net"
	"bytes"
	"proxy/contexts"
	"testing"
	assertion "util/test/assertion"
	"proxy/docker_client"
)

func Test_Parse_Cluster_Config_When_Default_Version_And_UpgradeTransition(testCtx *testing.T) {
	// given
	var (
		expectedError    error              = nil
		serverOne, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *contexts.Clusters = &contexts.Clusters{DockerHostEndpoint: "http://127.0.0.1:-1"}
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		outputStream bytes.Buffer
	)
	expectedClusters.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, RequestCounter: -1, Uuid: uuidGenerator(), SessionTimeout: 0, Mode: contexts.InstantMode, Version: "0.0"})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Default_Mode(testCtx *testing.T) {
	// given
	var (
		expectedError    error              = nil
		serverOne, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *contexts.Clusters = &contexts.Clusters{DockerHostEndpoint: "http://127.0.0.1:-1"}
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"sessionTimeout": float64(60)}, "version": "1.0"}}
		outputStream bytes.Buffer
	)
	expectedClusters.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, RequestCounter: -1, Uuid: uuidGenerator(), SessionTimeout: 60, Mode: contexts.SessionMode, Version: "1.0"})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Config_Valid_No_Defaults(testCtx *testing.T) {
	// given
	var (
		expectedError    error              = nil
		serverOne, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *contexts.Clusters = &contexts.Clusters{DockerHostEndpoint: "http://127.0.0.1:-1"}
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "INSTANT"}, "version": "1.0"}}
		outputStream bytes.Buffer
	)
	expectedClusters.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, RequestCounter: -1, Uuid: uuidGenerator(), Mode: contexts.InstantMode, Version: "1.0"})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Concurrent_Transition_Mode(testCtx *testing.T) {
	// given
	var (
		expectedError    error              = nil
		serverOne, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *contexts.Clusters = &contexts.Clusters{DockerHostEndpoint: "http://127.0.0.1:-1"}
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "CONCURRENT"}, "version": "1.0"}}
		outputStream bytes.Buffer
	)
	expectedClusters.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, RequestCounter: -1, Uuid: uuidGenerator(), Mode: contexts.ConcurrentMode, Version: "1.0"})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Gradual_Transition_Mode(testCtx *testing.T) {
	// given
	var (
		expectedError    error              = nil
		serverOne, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *contexts.Clusters = &contexts.Clusters{DockerHostEndpoint: "http://127.0.0.1:-1"}
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "GRADUAL", "percentageTransitionPerRequest": float64(0.01)}, "version": "1.0"}}
		outputStream bytes.Buffer
	)
	expectedClusters.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, RequestCounter: -1, Uuid: uuidGenerator(), Mode: contexts.GradualMode, PercentageTransitionPerRequest: float64(0.01), Version: "1.0"})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Config_Valid_With_UUID(testCtx *testing.T) {
	// given
	var (
		expectedError    error              = nil
		serverOne, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *contexts.Clusters = &contexts.Clusters{DockerHostEndpoint: "http://127.0.0.1:-1"}
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"uuid": "1027596f-1034-11e4-8334-600308a82410", "servers": serversConfig, "version": "1.0"}}
		outputStream bytes.Buffer
	)
	expectedClusters.Add(&contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, RequestCounter: -1, Uuid: uuid.Parse("1027596f-1034-11e4-8334-600308a82410"), Mode: contexts.InstantMode, Version: "1.0"})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_When_Cluster_Nil(testCtx *testing.T) {
	// given
	var (
		jsonConfig                          = map[string]interface{}{"cluster": nil}
		expectedError                       = errors.New("Invalid cluster configuration - \"cluster\" config missing")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualRouter, err := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_When_Server_List_Empty(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                       = errors.New("Invalid cluster configuration - \"servers\" list must contain at least one entry")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_Config_When_Gradual_Transition_Mode_And_No_PercentageTransition(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "GRADUAL"}, "version": "1.0"}}
		expectedError    error              = errors.New("Invalid cluster configuration - \"percentageTransitionPerRequest\" must be specified in \"upgradeTransition\" for mode \"GRADUAL\"")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Session_Mode_And_No_Timeout(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "SESSION"}, "version": "1.0"}}
		expectedError    error              = errors.New("Invalid cluster configuration - \"sessionTimeout\" must be specified in \"upgradeTransition\" for mode \"SESSION\"")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Servers_List_Missing(testCtx *testing.T) {
	// given
	var (
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"upgradeTransition": map[string]interface{}{"mode": "INSTANT"}, "version": "1.0"}}
		expectedError    error              = errors.New("Invalid cluster configuration - \"cluster\" must contain \"servers\" or \"containers\" list")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_When_No_IP(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"port": 1024}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                       = errors.New("Invalid server address [%!s(<nil>):1024] - missing brackets in address %!s(<nil>):1024")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_When_IP_Invalid(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "inv@lid", "port": 1024}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                       = errors.New("Invalid server address [inv@lid:1024] - lookup inv@lid: no such host")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_When_No_Port(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1"}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                       = errors.New("Invalid server address [127.0.0.1:<nil>] - unknown port tcp/<nil>")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_When_Port_Invalid(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": "invalid"}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                       = errors.New("Invalid server address [127.0.0.1:invalid] - unknown port tcp/invalid")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_Config_When_Invalid_Mode(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "INVALID"}, "version": "1.0"}}
		expectedError    error              = errors.New("Invalid cluster configuration - \"upgradeTransition.mode\" should be \"INSTANT\", \"SESSION\", \"GRADUAL\" or \"CONCURRENT\"")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Invalid_Instance_Mode_Timeout_Combination(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "INSTANT", "sessionTimeout": float64(60)}, "version": "1.0"}}
		expectedError    error              = errors.New("Invalid cluster configuration - \"sessionTimeout\" should not be specified when \"mode\" is not \"SESSION\"")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Invalid_Concurrent_Mode_Timeout_Combination(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "CONCURRENT", "sessionTimeout": float64(60)}, "version": "1.0"}}
		expectedError    error              = errors.New("Invalid cluster configuration - \"sessionTimeout\" should not be specified when \"mode\" is not \"SESSION\"")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Invalid_Instance_Mode_TransitionPerRequest_Combination(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "INSTANT", "percentageTransitionPerRequest": float64(0.01)}, "version": "1.0"}}
		expectedError    error              = errors.New("Invalid cluster configuration - \"percentageTransitionPerRequest\" should not be specified when \"mode\" is not \"GRADUAL\"")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Invalid_Concurrent_Mode_TansitionPerRequest_Combination(testCtx *testing.T) {
	// given
	var (
		serversConfig                       = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		jsonConfig                          = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "CONCURRENT", "percentageTransitionPerRequest": float64(0.01)}, "version": "1.0"}}
		expectedError    error              = errors.New("Invalid cluster configuration - \"percentageTransitionPerRequest\" should not be specified when \"mode\" is not \"GRADUAL\"")
		expectedClusters *contexts.Clusters = nil
		outputStream bytes.Buffer
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator, false)(jsonConfig, &docker_client.DockerHost{Ip: "127.0.0.1", Port: -1}, &outputStream)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Serialise_Cluster_When_Instant_Mode(testCtx *testing.T) {
	// given
	var (
		serverOne, _                         = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                         = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		uuidValue                            = uuidGenerator()
		serversConfig                        = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		expectedJsonConfig                   = map[string]interface{}{"cluster": map[string]interface{}{"uuid": uuidValue.String(), "servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "INSTANT"}, "version": "1.0"}}
		cluster            *contexts.Cluster = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, Uuid: uuidValue, Mode: contexts.InstantMode, Version: "1.0"}
		
	)

	// when
	actualJsonConfig := serialiseCluster(cluster)

	// then
	assertion.AssertDeepEqual("Correct JSON config", testCtx, expectedJsonConfig, actualJsonConfig)
}

func Test_Serialise_Cluster_When_Session_Mode(testCtx *testing.T) {
	// given
	var (
		serverOne, _                         = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                         = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		uuidValue                            = uuidGenerator()
		serversConfig                        = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		expectedJsonConfig                   = map[string]interface{}{"cluster": map[string]interface{}{"uuid": uuidValue.String(), "servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "SESSION", "sessionTimeout": int64(10)}, "version": "1.0"}}
		cluster            *contexts.Cluster = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, Uuid: uuidValue, SessionTimeout: int64(10), Mode: contexts.SessionMode, Version: "1.0"}
	)

	// when
	actualJsonConfig := serialiseCluster(cluster)

	// then
	assertion.AssertDeepEqual("Correct JSON config", testCtx, expectedJsonConfig, actualJsonConfig)
}

func Test_Serialise_Cluster_When_Gradual_Mode(testCtx *testing.T) {
	// given
	var (
		serverOne, _                         = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                         = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		uuidValue                            = uuidGenerator()
		serversConfig                        = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		expectedJsonConfig                   = map[string]interface{}{"cluster": map[string]interface{}{"uuid": uuidValue.String(), "servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "GRADUAL", "percentageTransitionPerRequest": float64(0.01)}, "version": "1.0"}}
		cluster            *contexts.Cluster = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, Uuid: uuidValue, PercentageTransitionPerRequest: float64(0.01), Mode: contexts.GradualMode, Version: "1.0"}
	)

	// when
	actualJsonConfig := serialiseCluster(cluster)

	// then
	assertion.AssertDeepEqual("Correct JSON config", testCtx, expectedJsonConfig, actualJsonConfig)
}

func Test_Serialise_Cluster_When_Concurrent_Mode(testCtx *testing.T) {
	// given
	var (
		serverOne, _                         = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                         = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		uuidValue                            = uuidGenerator()
		serversConfig                        = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}, map[string]interface{}{"hostname": "127.0.0.1", "port": 1025}}
		expectedJsonConfig                   = map[string]interface{}{"cluster": map[string]interface{}{"uuid": uuidValue.String(), "servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "CONCURRENT"}, "version": "1.0"}}
		cluster            *contexts.Cluster = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverOne, Host: "127.0.0.1", Port: "1024"}, &contexts.BackendAddress{Address: serverTwo, Host: "127.0.0.1", Port: "1025"}}, Uuid: uuidValue, Mode: contexts.ConcurrentMode, Version: "1.0"}
	)

	// when
	actualJsonConfig := serialiseCluster(cluster)

	// then
	assertion.AssertDeepEqual("Correct JSON config", testCtx, expectedJsonConfig, actualJsonConfig)
}

func Test_Serialise_Containers(testCtx *testing.T) {
	// given
	var (
		serverAddress, _                     = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		uuidValue                            = uuidGenerator()
		serversConfig                        = []interface{}{map[string]interface{}{"hostname": "127.0.0.1", "port": 1024}}
		containersConfig                     = []*docker_client.DockerConfig{&docker_client.DockerConfig{Image: "test image"}}
		expectedJsonConfig                   = map[string]interface{}{"cluster": map[string]interface{}{"uuid": uuidValue.String(), "servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "INSTANT"}, "version": "1.0", "containers": containersConfig}}
		cluster            *contexts.Cluster = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: serverAddress, Host: "127.0.0.1", Port: "1024"}}, DockerConfigurations: containersConfig, Uuid: uuidValue, Mode: contexts.InstantMode, Version: "1.0"}

	)

	// when
	actualJsonConfig := serialiseCluster(cluster)

	// then
	assertion.AssertDeepEqual("Correct JSON config", testCtx, expectedJsonConfig, actualJsonConfig)
}
