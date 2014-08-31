package proxy

import (
	"testing"
	assertion "util/test/assertion"
	"strconv"
	"io/ioutil"
	"errors"
)

func WriteConfigFile(proxyPort int, configPort int, uuid string, serverPorts []int, version string) string {
	fileName := "/tmp/system_test_config.json"
	data := "{\"proxy\":{\"hostname\":\"localhost\",\"port\":" + strconv.Itoa(proxyPort) + "},\"configService\":{\"port\":" + strconv.Itoa(configPort) + "},\"cluster\":{\"servers\":["
	for index, serverPort := range serverPorts {
		if index > 0 {
			data += ","
		}
		data += "{\"hostname\":\"127.0.0.1\",\"port\":"+strconv.Itoa(serverPort)+"}"
	}
	data += "]"
	if len(uuid) > 0 {
		data += ",\"uuid\":\""+uuid+"\""
	}
	if len(version) > 0 {
		data += ",\"version\":\""+version+"\""
	}
	data += "}}"

	ioutil.WriteFile(fileName, []byte(data), 0644)
	return fileName
}

func Test_Read_Config_When_File_Exists(testCtx *testing.T) {
	// given
	var (
		proxyPort int       = 1234
		configPort int      = 4321
		uuid string         = "a37a290f-2088-11e4-b3a6-600308a8245e"
		serverPorts []int   = []int{1024, 1025}
		version string      = "0.5"
		fileName            = WriteConfigFile(proxyPort, configPort, uuid, serverPorts, version)
		expectedByteArray   = []byte("{\"proxy\":{\"hostname\":\"localhost\",\"port\":1234},\"configService\":{\"port\":4321},\"cluster\":{\"servers\":[{\"hostname\":\"127.0.0.1\",\"port\":1024},{\"hostname\":\"127.0.0.1\",\"port\":1025}],\"uuid\":\"a37a290f-2088-11e4-b3a6-600308a8245e\",\"version\":\"0.5\"}}")
		expectedError error = nil
	)

	// when
	actualByteArray, actualError := readConfigFile(fileName)

	// then
	assertion.AssertDeepEqual("Correct Byte Array read from the file", testCtx, expectedByteArray, actualByteArray)
	assertion.AssertDeepEqual("Correct error while reading file", testCtx, expectedError, actualError)
}

func Test_Read_Config_When_File_Not_Exists(testCtx *testing.T) {
	// given
	var (
		fileName                 = "does_not_exist.json"
		expectedByteArray []byte = nil
		expectedError error      = errors.New("Error open does_not_exist.json: no such file or directory reading config file [does_not_exist.json]")
	)

	// when
	actualByteArray, actualError := readConfigFile(fileName)

	// then
	assertion.AssertDeepEqual("Correct Byte Array read from the file", testCtx, expectedByteArray, actualByteArray)
	assertion.AssertDeepEqual("Correct error while reading file", testCtx, expectedError, actualError)
}



