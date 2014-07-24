package proxy_c

import (
	"testing"
	assertion "util/test/assertion"
	"errors"
)

func Test_Load_Config_When_File_Valid_Server_list(testCtx *testing.T) {
	// given
	var (
		fileName                          = "does_not_exist.json"
		expectedError                     = errors.New("Invalid proxy configuration - \"proxy\" config missing")
		expectedLoadBalance *LoadBalancer = nil
	)

	// when
	actualLoadBalancer, actualError := loadConfig(fileName)

	// then
	assertion.AssertDeepEqual("Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Load Balancer", testCtx, expectedLoadBalance, actualLoadBalancer)
}


