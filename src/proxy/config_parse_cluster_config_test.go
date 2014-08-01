package proxy

import (
	"testing"
	assertion "util/test/assertion"
	"net"
	"errors"
	"code.google.com/p/go-uuid/uuid"
	"proxy/stages"
)

func Test_Parse_Cluster_Config_When_Default_Version_And_UpgradeTransition(testCtx *testing.T) {
	// given
	var (
		expectedError error               = nil
		serverOne, _                      = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                      = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *stages.Clusters = &stages.Clusters{}
	serversConfig = []interface{}{map[string]interface{}{"ip":"127.0.0.1", "port":1024}, map[string]interface{}{"ip":"127.0.0.1", "port":1025}}
	jsonConfig = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
	)
	expectedClusters.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuidGenerator(), SessionTimeout: 0, Mode: stages.InstantMode, Version: 0.0})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Default_Mode(testCtx *testing.T) {
	// given
	var (
		expectedError error               = nil
		serverOne, _                      = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                      = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *stages.Clusters = &stages.Clusters{}
	serversConfig = []interface{}{map[string]interface{}{"ip":"127.0.0.1", "port":1024}, map[string]interface{}{"ip":"127.0.0.1", "port":1025}}
	jsonConfig = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"sessionTimeout": float64(60)}, "version": 1.0}}
	)
	expectedClusters.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuidGenerator(), SessionTimeout: 60, Mode: stages.SessionMode, Version: 1.0})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Config_Valid_No_Defaults(testCtx *testing.T) {
	// given
	var (
		expectedError error               = nil
		serverOne, _                      = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                      = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *stages.Clusters = &stages.Clusters{}
	serversConfig = []interface{}{map[string]interface{}{"ip":"127.0.0.1", "port":1024}, map[string]interface{}{"ip":"127.0.0.1", "port":1025}}
	jsonConfig = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "INSTANT"}, "version": 1.0}}
	)
	expectedClusters.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuidGenerator(), Mode: stages.InstantMode, Version: 1.0})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Config_Valid_With_UUID(testCtx *testing.T) {
	// given
	var (
		expectedError error               = nil
		serverOne, _                      = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                      = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedClusters *stages.Clusters = &stages.Clusters{}
	serversConfig = []interface{}{map[string]interface{}{"ip":"127.0.0.1", "port":1024}, map[string]interface{}{"ip":"127.0.0.1", "port":1025}}
	jsonConfig = map[string]interface{}{"cluster": map[string]interface{}{"uuid": "1027596f-1034-11e4-8334-600308a82410", "servers": serversConfig, "version": 1.0}}
	)
	expectedClusters.Add(&stages.Cluster{BackendAddresses: []*net.TCPAddr{serverOne, serverTwo}, RequestCounter: -1, Uuid: uuid.Parse("1027596f-1034-11e4-8334-600308a82410"), Mode: stages.InstantMode, Version: 1.0})

	// when
	actualClusters, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_When_Cluster_Nil(testCtx *testing.T) {
	// given
	var (
		jsonConfig                        = map[string]interface{}{"cluster": nil}
		expectedError                     = errors.New("Invalid cluster configuration - \"cluster\" config missing")
		expectedClusters *stages.Clusters = nil
	)

	// when
	actualRouter, err := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_When_Server_List_Empty(testCtx *testing.T) {
	// given
	var (
		serversConfig                     = []interface{}{}
		jsonConfig                        = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                     = errors.New("Invalid cluster configuration - \"servers\" list must contain at least one entry")
		expectedClusters *stages.Clusters = nil
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_When_No_IP(testCtx *testing.T) {
	// given
	var (
		serversConfig                     = []interface{}{map[string]interface{}{"port":1024}}
		jsonConfig                        = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                     = errors.New("Invalid server address [%!s(<nil>):1024] - missing brackets in address %!s(<nil>):1024")
		expectedClusters *stages.Clusters = nil
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_When_IP_Invalid(testCtx *testing.T) {
	// given
	var (
		serversConfig                     = []interface{}{map[string]interface{}{"ip":"inv@lid", "port":1024}}
		jsonConfig                        = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                     = errors.New("Invalid server address [inv@lid:1024] - lookup inv@lid: no such host")
		expectedClusters *stages.Clusters = nil
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_When_No_Port(testCtx *testing.T) {
	// given
	var (
		serversConfig                     = []interface{}{map[string]interface{}{"ip":"127.0.0.1"}}
		jsonConfig                        = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                     = errors.New("Invalid server address [127.0.0.1:<nil>] - unknown port tcp/<nil>")
		expectedClusters *stages.Clusters = nil
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_When_Port_Invalid(testCtx *testing.T) {
	// given
	var (
		serversConfig                     = []interface{}{map[string]interface{}{"ip": "127.0.0.1", "port": "invalid"}}
		jsonConfig                        = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                     = errors.New("Invalid server address [127.0.0.1:invalid] - unknown port tcp/invalid")
		expectedClusters *stages.Clusters = nil
	)

	// when
	actualRouter, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualRouter)
}

func Test_Parse_Cluster_Config_When_Invalid_Mode(testCtx *testing.T) {
	// given
	var (
		serversConfig                     = []interface{}{map[string]interface{}{"ip":"127.0.0.1", "port":1024}, map[string]interface{}{"ip":"127.0.0.1", "port":1025}}
		jsonConfig                        = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "INVALID"}, "version": 1.0}}
		expectedError error               = errors.New("Invalid cluster configuration - \"upgradeTransition.mode\" should be \"INSTANT\" or \"SESSION\"")
		expectedClusters *stages.Clusters = nil
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}

func Test_Parse_Cluster_Config_When_Invalid_Mode_Timeout_Combination(testCtx *testing.T) {
	// given
	var (
		serversConfig                     = []interface{}{map[string]interface{}{"ip":"127.0.0.1", "port":1024}, map[string]interface{}{"ip":"127.0.0.1", "port":1025}}
		jsonConfig                        = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "upgradeTransition": map[string]interface{}{"mode": "INSTANT", "sessionTimeout": float64(60)}, "version": 1.0}}
		expectedError error               = errors.New("Invalid cluster configuration - \"sessionTimeout\" should not be specified when \"mode\" is \"INSTANT\"")
		expectedClusters *stages.Clusters = nil
	)

	// when
	actualClusters, actualError := parseClusters(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedClusters, actualClusters)
}