package proxy

import (
	"testing"
	assertion "util/test/assertion"
	"errors"
)

func Test_Load_Config_When_File_Valid_Server_list(testCtx *testing.T) {
	// given
	var (
		fileName                          = "does_not_exist.json"
		expectedError                     = errors.New("Error open does_not_exist.json: no such file or directory reading config file [does_not_exist.json]")
		expectedLoadBalance *Proxy = nil
	)

	// when
	actualProxy, actualError := loadConfig(fileName)

	// then
	assertion.AssertDeepEqual("Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Load Balancer", testCtx, expectedLoadBalance, actualProxy)
}


