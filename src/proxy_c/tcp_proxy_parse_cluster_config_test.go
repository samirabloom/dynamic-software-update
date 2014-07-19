package proxy_c

import (
	"testing"
	assertion "util/test/assertion"
	"net"
	"errors"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
)

var uuidGenerator = func(uuidValue uuid.UUID) func() uuid.UUID {
	return func() uuid.UUID {
		return uuidValue
	}
}(uuid.NewUUID())

func Test_Parse_Cluster_Config_When_Config_Valid(testCtx *testing.T) {
	// given
	var (
		expectedError error                 = nil
		serverOne, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		serverTwo, _                        = net.ResolveTCPAddr("tcp", "127.0.0.1:1025")
		expectedRoutingContexts             = &RoutingContexts{all: make(map[string]*RoutingContext)}
//		serverOneConfig         interface{} = map[string]interface{}{"ip": "127.0.0.1", "port": 1024}
//		serverTwoConfig         interface{} = map[string]interface{}{"ip": "127.0.0.1", "port": 1025}
//		serversConfig                       = []map[string]interface{}{{"ip": "127.0.0.1", "port": 1024}, {"ip": "127.0.0.1", "port": 1025}}
		serversConfig                       = []interface {}{map[string]interface {}{"ip":"127.0.0.1", "port":1024}, map[string]interface {}{"ip":"127.0.0.1", "port":1025}}
		jsonConfig                          = map[string]interface{}{"servers": serversConfig}
	)
	expectedRoutingContexts.Add(&RoutingContext{backendAddresses: []*net.TCPAddr{serverOne, serverTwo}, requestCounter: -1, uuid: uuidGenerator()})

	// when
	actualMockRouter, actualError := parseClusterConfig(uuidGenerator)(jsonConfig)
	fmt.Printf("actualMockRouter %#v\n", actualMockRouter)
	fmt.Printf("expectedRoutingContexts %#v\n", expectedRoutingContexts)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, actualMockRouter, expectedRoutingContexts)
}



func Test_Parse_Cluster_When_Server_List_Empty(testCtx *testing.T) {
	// given
	var (
		serversConfig = []interface {}{}
		jsonConfig    = map[string]interface{}{"servers": serversConfig}
		expectedError = errors.New("Invalid cluster configuration - \"servers\" JSON field missing or invalid")
	)

	// when
	actualMockRouter, actualError := parseClusterConfig(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, actualMockRouter, nil)
}


func Test_Parse_Cluster_When_Servers_Nil(testCtx *testing.T) {
	// given
	var (
		jsonConfig    = map[string]interface{}{"servers": nil}
		expectedError = errors.New("Invalid cluster configuration - \"servers\" JSON field missing or invalid")
	)

	// when
	actualRouter, err := parseClusterConfig(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, err, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, actualRouter, nil)
}

func TestParse_Cluster_When_No_IP(testCtx *testing.T) {
	// given
	var (
		serversConfig = []interface {}{map[string]interface {}{"port":1024}}
		jsonConfig    = map[string]interface{}{"server_range": serversConfig}
		expectedError = errors.New("Invalid cluster configuration - \"servers\" JSON field missing or invalid")
	)

	// when
	actualRouter, actualError := parseClusterConfig(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, actualRouter, nil)
}

func TestParse_Cluster_When_IP_Invalid(testCtx *testing.T) {
	// given
	var (
		serversConfig = []interface {}{map[string]interface {}{"ip":"inv@lid", "port":1024}}
		jsonConfig    = map[string]interface{}{"server_range": serversConfig}
		expectedError = errors.New("Invalid cluster configuration - \"servers\" JSON field missing or invalid")
	)

	// when
	actualRouter, actualError := parseClusterConfig(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, actualRouter, nil)
}

func Test_Parse_Cluster_When_No_Port(testCtx *testing.T) {
	// given
	var (
		serversConfig = []interface {}{map[string]interface {}{"ip":"127.0.0.1"}}
		jsonConfig    = map[string]interface{}{"server_range": serversConfig}
		expectedError = errors.New("Invalid cluster configuration - \"servers\" JSON field missing or invalid")
	)

	// when
	actualRouter, actualError := parseClusterConfig(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, actualRouter, nil)
}

func Test_Parse_Cluster_When_Port_Invalid(testCtx *testing.T) {
	// given
	var (
		serversConfig = []interface {}{map[string]interface {}{"ip": "127.0.0.1", "port": "invalid"}}
		jsonConfig    = map[string]interface{}{"server_range": serversConfig}
		expectedError = errors.New("Invalid cluster configuration - \"servers\" JSON field missing or invalid")
	)

	// when
	actualRouter, actualError := parseClusterConfig(uuidGenerator)(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct Routing Contexts", testCtx, actualRouter, nil)
}
