package proxy_c

import (
	"testing"
	assertion "util/test/assertion"
	"errors"
)

func Test_Read_Config_When_File_Valid_Server_list(testCtx *testing.T) {
	// given
	var (
		fileName      = new(string)
		expectedError = errors.New("Invalid proxy configuration - \"proxy\" JSON field missing or invalid")
	)
	*fileName = "does_not_exist.json"

	// when
	actualLoadBalancer, actualError := loadConfig(fileName)

	// then
	assertion.AssertDeepEqual("Correct Error", testCtx, actualError, expectedError)
	assertion.AssertDeepEqual("Correct Load Balancer", testCtx, actualLoadBalancer, nil)
}


