package proxy

import (
	"testing"
	assertion "util/test/assertion"
	"errors"
)

func Test_Parse_Config_Service_When_Config_Valid(testCtx *testing.T) {
	// given
	var (
		expectedError error           = nil
		expectedConfigServicePort int = 1234
		jsonConfig                    = map[string]interface{}{"configService": map[string]interface{}{"port": float64(expectedConfigServicePort)}}
	)

	// when
	actualConfigServicePort, actualErr := parseConfigService(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Config Service Port", testCtx, expectedConfigServicePort, actualConfigServicePort)
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
}

func Test_Parse_Config_Service_When_No_Port(testCtx *testing.T) {
	// given
	var (
		expectedError error = errors.New("Invalid config service configuration - \"port\" is missing from \"configService\" config")
		jsonConfig          = map[string]interface{}{"configService": map[string]interface{}{"port": nil}}
	)

	// when
	actualConfigServicePort, actualErr := parseConfigService(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Config Service Port", testCtx, 0, actualConfigServicePort)
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
}

func Test_Parse_Config_Service_When_No_ConfigService(testCtx *testing.T) {
	// given
	var (
		expectedError error = errors.New("Invalid proxy configuration - \"configService\" config missing")
		jsonConfig          = map[string]interface{}{"configService": nil}
	)

	// when
	actualConfigServicePort, actualErr := parseConfigService(jsonConfig)

	// then
	assertion.AssertDeepEqual("Correct Config Service Port", testCtx, 0, actualConfigServicePort)
	assertion.AssertDeepEqual("Correct Proxy Error", testCtx, expectedError, actualErr)
}
