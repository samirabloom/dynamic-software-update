package proxy_c

import (
	"testing"
	assertion "util/test/assertion"
	"net"
	"errors"
	"code.google.com/p/go-uuid/uuid"
)

func Test_Parse_Cluster_Config_When_Config_Valid(testCtx *testing.T) {
	// given
	var (
		expectedError error                      = nil
		serverOne, _                             = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                             = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedRoutingContexts *RoutingContexts = &RoutingContexts{}
		serversConfig                            = []interface{}{map[string]interface{}{"ip":"127.0.0.1", "port":1024}, map[string]interface{}{"ip":"127.0.0.1", "port":1025}}
		jsonConfig                               = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig, "version": 1.0}}
	)
	expectedRoutingContexts.Add(&RoutingContext{backendAddresses: []*net.TCPAddr{serverOne, serverTwo}, requestCounter: -1, uuid: uuidGenerator(), version: 1.0})

	// when
	actualMockRouter, actualError := parseRoutingContexts(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedRoutingContexts, actualMockRouter)
}

func Test_Parse_Cluster_Config_When_Config_Valid_With_UUID(testCtx *testing.T) {
	// given
	var (
		expectedError error                      = nil
		serverOne, _                             = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                             = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedRoutingContexts *RoutingContexts = &RoutingContexts{}
		serversConfig                            = []interface{}{map[string]interface{}{"ip":"127.0.0.1", "port":1024}, map[string]interface{}{"ip":"127.0.0.1", "port":1025}}
		jsonConfig                               = map[string]interface{}{"cluster": map[string]interface{}{"uuid": "1027596f-1034-11e4-8334-600308a82410", "servers": serversConfig, "version": 1.0}}
	)
	expectedRoutingContexts.Add(&RoutingContext{backendAddresses: []*net.TCPAddr{serverOne, serverTwo}, requestCounter: -1, uuid: uuid.Parse("1027596f-1034-11e4-8334-600308a82410"), version: 1.0})

	// when
	actualMockRouter, actualError := parseRoutingContexts(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedRoutingContexts, actualMockRouter)
}

func Test_Parse_Cluster_When_Cluster_Nil(testCtx *testing.T) {
	// given
	var (
		jsonConfig                               = map[string]interface{}{"cluster": nil}
		expectedError                            = errors.New("Invalid cluster configuration - \"cluster\" config missing")
		expectedRoutingContexts *RoutingContexts = nil
	)

	// when
	actualRouter, err := parseRoutingContexts(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedRoutingContexts, actualRouter)
}

func Test_Parse_Cluster_When_Server_List_Empty(testCtx *testing.T) {
	// given
	var (
		serversConfig                            = []interface{}{}
		jsonConfig                               = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                            = errors.New("Invalid cluster configuration - \"servers\" list must contain at least one entry")
		expectedRoutingContexts *RoutingContexts = nil
	)

	// when
	actualRouter, actualError := parseRoutingContexts(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedRoutingContexts, actualRouter)
}

func Test_Parse_Cluster_When_No_IP(testCtx *testing.T) {
	// given
	var (
		serversConfig                            = []interface{}{map[string]interface{}{"port":1024}}
		jsonConfig                               = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                            = errors.New("Invalid server address [%!s(<nil>):1024] - missing brackets in address %!s(<nil>):1024")
		expectedRoutingContexts *RoutingContexts = nil
	)

	// when
	actualRouter, actualError := parseRoutingContexts(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedRoutingContexts, actualRouter)
}

func Test_Parse_Cluster_When_IP_Invalid(testCtx *testing.T) {
	// given
	var (
		serversConfig                            = []interface{}{map[string]interface{}{"ip":"inv@lid", "port":1024}}
		jsonConfig                               = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                            = errors.New("Invalid server address [inv@lid:1024] - lookup inv@lid: no such host")
		expectedRoutingContexts *RoutingContexts = nil
	)

	// when
	actualRouter, actualError := parseRoutingContexts(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedRoutingContexts, actualRouter)
}

func Test_Parse_Cluster_When_No_Port(testCtx *testing.T) {
	// given
	var (
		serversConfig                            = []interface{}{map[string]interface{}{"ip":"127.0.0.1"}}
		jsonConfig                               = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                            = errors.New("Invalid server address [127.0.0.1:<nil>] - unknown port tcp/<nil>")
		expectedRoutingContexts *RoutingContexts = nil
	)

	// when
	actualRouter, actualError := parseRoutingContexts(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedRoutingContexts, actualRouter)
}

func Test_Parse_Cluster_When_Port_Invalid(testCtx *testing.T) {
	// given
	var (
		serversConfig                            = []interface{}{map[string]interface{}{"ip": "127.0.0.1", "port": "invalid"}}
		jsonConfig                               = map[string]interface{}{"cluster": map[string]interface{}{"servers": serversConfig}}
		expectedError                            = errors.New("Invalid server address [127.0.0.1:invalid] - unknown port tcp/invalid")
		expectedRoutingContexts *RoutingContexts = nil
	)

	// when
	actualRouter, actualError := parseRoutingContexts(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, expectedRoutingContexts, actualRouter)
}
