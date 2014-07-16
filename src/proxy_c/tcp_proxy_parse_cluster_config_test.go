package proxy_c

import (
	"testing"
	assertion "util/test/assertion"
	"net"
)

// calling ParseProxy with
// Proxy not nil
// Err nil
func Test_Parse_Cluster_Config_Proxy_Not_Nil_Err_Nil(testCtx *testing.T) {
	// given
	var (
		expectedError error = nil
		backendBaseAddr, _  = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
		expectedMockRouter  = &RangeRoutingContext{backendBaseAddr: backendBaseAddr, clusterSize: 8}
		mockProxyConfig     = map[string]interface{}{"ip": "127.0.0.1", "port": 1024, "clusterSize": "8"}
		jsonConfig          = map[string]interface{}{"server_range": mockProxyConfig}
	)

	// when
	actualMockRouter, actualError := parseClusterConfig(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, actualMockRouter, expectedMockRouter)
}

// TODO add more invalid and missing as per proxy

//func Parse_Cluster_Config(jsonConfig map[string]interface{}) (router Router, err error) {
//	if jsonConfig["server_range"] != nil {
//		var serverConfig map[string]interface{} = jsonConfig["server_range"].(map[string]interface{})
//		backendBaseAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%v", serverConfig["ip"], serverConfig["port"]))
//		if err != nil {
//			loggerFactory().Error("Invalid address [" + fmt.Sprintf("%s:%v", serverConfig["ip"], serverConfig["port"]) + "]")
//			return nil, err
//		}
//
//		clusterSize , err := strconv.Atoi(serverConfig["clusterSize"].(string))
//		if err != nil {
//			loggerFactory().Error("Cluster Size not a valid integer [" + serverConfig["clusterSize"].(string) + "]")
//			return nil, err
//		}
//
//		router = &RangeRoutingContext{
//			backendBaseAddr:  backendBaseAddr,
//			clusterSize: clusterSize,
//		}
//	}
//}
