package proxy_c

import (
	"testing"
	assertion "util/test/assertion"
	"fmt"
)

func Test_Read_Config_When_File_Exists(testCtx *testing.T) {
	// given
	var (
		fileName          = new(string)
		expectedByteArray = []byte("{\n    \"proxy\": {\n        \"ip\": \"localhost\",\n        \"port\": 1234\n    },\n    \"server_range\":{\n        \"ip\": \"127.0.0.1\",\n        \"port\": 1024,\n        \"clusterSize\": \"8\"\n    }\n}")
	)
	*fileName = "test_range_config.json"

	// when
	actualByteArray := readConfigFile(fileName)

	// then
	assertion.AssertDeepEqual("Correct Byte Array read from the file", testCtx, actualByteArray, expectedByteArray)
}

func Test_Read_Config_When_File_Not_Exists(testCtx *testing.T) {
	// given
	var (
		fileName          = new(string)
	)
	*fileName = "does_not_exist.json"

	// when
	actualByteArray := readConfigFile(fileName)
	fmt.Printf("actualByteArray %s", actualByteArray)

	// then
	assertion.AssertDeepEqual("Correct Byte Array read from the file", testCtx, actualByteArray, nil)
}



