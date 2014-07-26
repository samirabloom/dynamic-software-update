package proxy

import (
	"testing"
	assertion "util/test/assertion"
)

func Test_Read_Config_When_File_Exists(testCtx *testing.T) {
	// given
	var (
		fileName          = "config_parse_read_config_file_test_config.json"
		expectedByteArray = []byte("{\n    \"proxy\": {\n        \"ip\": \"localhost\",\n        \"port\": 1234\n    },\n    \"cluster\": {\n        \"servers\": [\n            {\n                \"ip\": \"127.0.0.1\",\n                \"port\": 1024\n            }\n        ],\n        \"version\": 1.0,\n        \"upgradeTransition\": {\n            \"sessionTimeout\": 60\n        }\n    }\n}")
	)

	// when
	actualByteArray := readConfigFile(fileName)

	// then
	assertion.AssertDeepEqual("Correct Byte Array read from the file", testCtx, actualByteArray, expectedByteArray)
}

func Test_Read_Config_When_File_Not_Exists(testCtx *testing.T) {
	// given
	var (
		fileName                 = "does_not_exist.json"
		expectedByteArray []byte = nil
	)

	// when
	actualByteArray := readConfigFile(fileName)

	// then
	assertion.AssertDeepEqual("Correct Byte Array read from the file", testCtx, actualByteArray, expectedByteArray)
}



