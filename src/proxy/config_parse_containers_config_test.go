package proxy

import (
	"testing"
	"encoding/json"
	"proxy/docker_client"
	"util/test/assertion"
	"proxy/contexts"
	"net"
	"errors"
	"fmt"
)

func Test_Docker_Config_JSON_Unmarshal(testCtx *testing.T) {
	// given
	var dockerConfigJson = []byte("{\n" +
			"   \"image\": \"test image\",\n" +
			"   \"tag\": \"test tag\",\n" +
			"   \"alwaysPull\": true,\n" +
			"   \"name\": \"test name\",\n" +
			"   \"workingDir\": \"test workingDir\",\n" +
			"   \"entrypoint\": [\n" +
			"      \"test entrypoint one\",\n" +
			"      \"test entrypoint two\"\n" +
			"   ],\n" +
			"   \"environment\": [\n" +
			"      \"test environment one\",\n" +
			"      \"test environment two\"\n" +
			"   ],\n" +
			"   \"cmd\": [\n" +
			"      \"test cmd one\",\n" +
			"      \"test cmd two\"\n" +
			"   ],\n" +
			"   \"hostname\": \"test hostname\",\n" +
			"   \"volumes\": [\n" +
			"      \"test volumes one\",\n" +
			"      \"test volumes two\"\n" +
			"   ],\n" +
			"   \"volumesFrom\": [\n" +
			"      \"test volumesFrom one\",\n" +
			"      \"test volumesFrom two\"\n" +
			"   ],\n" +
			"   \"portBindings\": {\n" +
			"      \"8080\": [\n" +
			"         {\n" +
			"            \"hostIp\": \"127.0.0.1\",\n" +
			"            \"hostPort\": \"8080\"\n" +
			"         }\n" +
			"      ]\n" +
			"   },\n" +
			"   \"links\": [\n" +
			"      \"test links one\",\n" +
			"      \"test links two\"\n" +
			"   ],\n" +
			"   \"user\": \"test user\",\n" +
			"   \"lxcConf\": [\n" +
			"      {\n" +
			"         \"key\": \"test key\",\n" +
			"         \"value\": \"test value\"\n" +
			"      }\n" +
			"   ],\n" +
			"   \"privileged\": true\n" +
			"}")
	var expectedDockerConfig = &docker_client.DockerConfig{
		Image: "test image",
		Tag: "test tag",
		AlwaysPull: true,
		Name: "test name",
		WorkingDir: "test workingDir",
		Entrypoint: []string{"test entrypoint one", "test entrypoint two"},
		Environment: []string{"test environment one", "test environment two"},
		Cmd: []string{"test cmd one", "test cmd two"},
		Hostname: "test hostname",
		Volumes: []string{"test volumes one", "test volumes two"},
		VolumesFrom: []string{"test volumesFrom one", "test volumesFrom two"},
		PortBindings: map[docker_client.Port][]docker_client.PortBinding{docker_client.Port("8080"): []docker_client.PortBinding{docker_client.PortBinding{HostIp: "127.0.0.1", HostPort: "8080"}}},
		PortToProxy: 0,
		Links: []string{"test links one", "test links two"},
		User: "test user",
		Memory: 0,
		CpuShares: 0,
		LxcConf: []docker_client.KeyValuePair{docker_client.KeyValuePair{Key: "test key", Value: "test value"}},
		Privileged: true,
	}
	var actualDockerConfig = &docker_client.DockerConfig{}

	// when
	json.Unmarshal(dockerConfigJson, actualDockerConfig)

	// then
	assertion.AssertDeepEqual("Correct Docker Container Parsed Correctly", testCtx, expectedDockerConfig, actualDockerConfig)
}

func Test_Docker_Config_JSON_Marshal(testCtx *testing.T) {
	// given
	var dockerConfig = &docker_client.DockerConfig{
		Image: "test image",
		Tag: "test tag",
		AlwaysPull: true,
		Name: "test name",
		WorkingDir: "test workingDir",
		Entrypoint: []string{"test entrypoint one", "test entrypoint two"},
		Environment: []string{"test environment one", "test environment two"},
		Cmd: []string{"test cmd one", "test cmd two"},
		Hostname: "test hostname",
		Volumes: []string{"test volumes one", "test volumes two"},
		VolumesFrom: []string{"test volumesFrom one", "test volumesFrom two"},
		PortBindings: map[docker_client.Port][]docker_client.PortBinding{docker_client.Port("8080"): []docker_client.PortBinding{docker_client.PortBinding{HostIp: "127.0.0.1", HostPort: "8080"}}},
		PortToProxy: 0,
		Links: []string{"test links one", "test links two"},
		User: "test user",
		Memory: 0,
		CpuShares: 0,
		LxcConf: []docker_client.KeyValuePair{docker_client.KeyValuePair{Key: "test key", Value: "test value"}},
		Privileged: true,
	}
	var expectDockerConfigJson = []byte("{\n" +
			"   \"image\": \"test image\",\n" +
			"   \"tag\": \"test tag\",\n" +
			"   \"alwaysPull\": true,\n" +
			"   \"name\": \"test name\",\n" +
			"   \"workingDir\": \"test workingDir\",\n" +
			"   \"entrypoint\": [\n" +
			"      \"test entrypoint one\",\n" +
			"      \"test entrypoint two\"\n" +
			"   ],\n" +
			"   \"environment\": [\n" +
			"      \"test environment one\",\n" +
			"      \"test environment two\"\n" +
			"   ],\n" +
			"   \"cmd\": [\n" +
			"      \"test cmd one\",\n" +
			"      \"test cmd two\"\n" +
			"   ],\n" +
			"   \"hostname\": \"test hostname\",\n" +
			"   \"volumes\": [\n" +
			"      \"test volumes one\",\n" +
			"      \"test volumes two\"\n" +
			"   ],\n" +
			"   \"volumesFrom\": [\n" +
			"      \"test volumesFrom one\",\n" +
			"      \"test volumesFrom two\"\n" +
			"   ],\n" +
			"   \"portBindings\": {\n" +
			"      \"8080\": [\n" +
			"         {\n" +
			"            \"hostIp\": \"127.0.0.1\",\n" +
			"            \"hostPort\": \"8080\"\n" +
			"         }\n" +
			"      ]\n" +
			"   },\n" +
			"   \"links\": [\n" +
			"      \"test links one\",\n" +
			"      \"test links two\"\n" +
			"   ],\n" +
			"   \"user\": \"test user\",\n" +
			"   \"lxcConf\": [\n" +
			"      {\n" +
			"         \"key\": \"test key\",\n" +
			"         \"value\": \"test value\"\n" +
			"      }\n" +
			"   ],\n" +
			"   \"privileged\": true\n" +
			"}")

	// when
	actualDockerConfigJson, _ := json.MarshalIndent(dockerConfig, "", "   ")

	// then
	assertion.AssertDeepEqual("Correct Docker Container Serialized Correctly", testCtx, string(expectDockerConfigJson), string(actualDockerConfigJson))
}

func Test_Parse_Docker_Config_With_No_Error(testCtx *testing.T) {
	// given
	var dockerConfigJson = []byte("{\n" +
			"    \"containers\": [\n" +
			"        {\n" +
			"            \"image\": \"test image\", \n" +
			"            \"tag\": \"test tag\", \n" +
			"            \"alwaysPull\": true,\n" +
			"            \"name\": \"test name\", \n" +
			"            \"workingDir\": \"test workingDir\", \n" +
			"            \"entrypoint\": [\n" +
			"                \"test entrypoint one\", \n" +
			"                \"test entrypoint two\"\n" +
			"            ], \n" +
			"            \"environment\": [\n" +
			"                \"test environment one\", \n" +
			"                \"test environment two\"\n" +
			"            ], \n" +
			"            \"cmd\": [\n" +
			"                \"test cmd one\", \n" +
			"                \"test cmd two\"\n" +
			"            ], \n" +
			"            \"hostname\": \"test hostname\", \n" +
			"            \"volumes\": [\n" +
			"                \"test volumes one\", \n" +
			"                \"test volumes two\"\n" +
			"            ], \n" +
			"            \"volumesFrom\": [\n" +
			"                \"test volumesFrom one\", \n" +
			"                \"test volumesFrom two\"\n" +
			"            ], \n" +
			"            \"portToProxy\": 1024," +
			"            \"portBindings\": {\n" +
			"                \"80\": [\n" +
			"                    {\n" +
			"                        \"hostIp\": \"0.0.0.0\", \n" +
			"                        \"hostPort\": \"1024\"\n" +
			"                    }\n" +
			"                ]\n" +
			"            }, \n" +
			"            \"links\": [\n" +
			"                \"test links one\", \n" +
			"                \"test links two\"\n" +
			"            ], \n" +
			"            \"user\": \"test user\", \n" +
			"            \"lxcConf\": [\n" +
			"                {\n" +
			"                    \"key\": \"test key\", \n" +
			"                    \"value\": \"test value\"\n" +
			"                }\n" +
			"            ], \n" +
			"            \"privileged\": true\n" +
			"        }\n" +
			"    ]\n" +
			"}")
	var dockerHost = &docker_client.DockerHost{Ip: "127.0.0.1", Port: 1234}
	var expectedDockerConfigs = []*docker_client.DockerConfig{
		&docker_client.DockerConfig{
			Image: "test image",
			Tag: "test tag",
			AlwaysPull: true,
			Name: "test name",
			WorkingDir: "test workingDir",
			Entrypoint: []string{"test entrypoint one", "test entrypoint two"},
			Environment: []string{"test environment one", "test environment two"},
			Cmd: []string{"test cmd one", "test cmd two"},
			Hostname: "test hostname",
			Volumes: []string{"test volumes one", "test volumes two"},
			VolumesFrom: []string{"test volumesFrom one", "test volumesFrom two"},
			PortBindings: map[docker_client.Port][]docker_client.PortBinding{docker_client.Port("80"): []docker_client.PortBinding{docker_client.PortBinding{HostIp: "0.0.0.0", HostPort: "1024"}}},
			PortToProxy: 1024,
			Links: []string{"test links one", "test links two"},
			User: "test user",
			Memory: 0,
			CpuShares: 0,
			LxcConf: []docker_client.KeyValuePair{docker_client.KeyValuePair{Key: "test key", Value: "test value"}},
			Privileged: true,
		},
	}
	var containerAddress, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:1024")
	var expectedBackendAddresses = []*contexts.BackendAddress{
		&contexts.BackendAddress{
			Address: containerAddress,
			Host: dockerHost.Ip,
			Port: "1024",
		},
	}
	var expectedError error = nil

	// when
	var jsonConfig = make(map[string]interface{})
	err := json.Unmarshal(dockerConfigJson, &jsonConfig)
	if err != nil {
		testCtx.Fatalf("Error while parsing JSON %s\n", err)
	}
	actualDockerConfigs, actualBackendAddresses, actualError := parseContainers(jsonConfig["containers"], dockerHost)

	// then
	assertion.AssertDeepEqual("Correct Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Docker Containers Parsed Correctly", testCtx, expectedDockerConfigs, actualDockerConfigs)
	assertion.AssertDeepEqual("Correct Docker Backend Addresses", testCtx, expectedBackendAddresses, actualBackendAddresses)
}

func Test_Parse_Docker_Config_With_Minimum_Config(testCtx *testing.T) {
	// given
	var dockerConfigJson = []byte("{\n" +
			"    \"containers\": [\n" +
			"        {\n" +
			"            \"image\": \"test image\"\n" +
			"        }\n" +
			"    ]\n" +
			"}")
	var dockerHost = &docker_client.DockerHost{Ip: "127.0.0.1", Port: 1234}
	var expectedDockerConfigs = []*docker_client.DockerConfig{
		&docker_client.DockerConfig{
			Image: "test image",
		},
	}
	var expectedBackendAddresses = []*contexts.BackendAddress{}
	var expectedError error = nil

	// when
	var jsonConfig = make(map[string]interface{})
	err := json.Unmarshal(dockerConfigJson, &jsonConfig)
	if err != nil {
		testCtx.Fatalf("Error while parsing JSON %s\n", err)
	}
	actualDockerConfigs, actualBackendAddresses, actualError := parseContainers(jsonConfig["containers"], dockerHost)

	// then
	assertion.AssertDeepEqual("Correct Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Docker Containers Parsed Correctly", testCtx, expectedDockerConfigs, actualDockerConfigs)
	assertion.AssertDeepEqual("Correct Docker Backend Addresses", testCtx, expectedBackendAddresses, actualBackendAddresses)
}

func Test_Parse_Docker_Config_With_Port_As_Integer(testCtx *testing.T) {
	// given
	var dockerConfigJson = []byte("{\n" +
			"    \"containers\": [\n" +
			"        {\n" +
			"            \"image\": \"test image\", \n" +
			"            \"portToProxy\": 9090, \n" +
			"            \"portBindings\": {\n" +
			"                \"80\": [\n" +
			"                    {\n" +
			"                        \"hostIp\": \"0.0.0.0\", \n" +
			"                        \"hostPort\": \"9090\"\n" +
			"                    }\n" +
			"                ]\n" +
			"            }\n" +
			"        }\n" +
			"    ]\n" +
			"}")
	var dockerHost = &docker_client.DockerHost{Ip: "127.0.0.1", Port: 1234}
	var expectedDockerConfigs = []*docker_client.DockerConfig{
		&docker_client.DockerConfig{
			Image: "test image",
			PortToProxy: 9090,
			PortBindings: map[docker_client.Port][]docker_client.PortBinding{docker_client.Port("80"): []docker_client.PortBinding{docker_client.PortBinding{HostIp: "0.0.0.0", HostPort: "9090"}}},
		},
	}
	var containerAddress, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:9090", dockerHost.Ip))
	var expectedBackendAddresses = []*contexts.BackendAddress{
		&contexts.BackendAddress{
			Address: containerAddress,
			Host: dockerHost.Ip,
			Port: "9090",
		},
	}
	var expectedError error = nil

	// when
	var jsonConfig = make(map[string]interface{})
	err := json.Unmarshal(dockerConfigJson, &jsonConfig)
	if err != nil {
		testCtx.Fatalf("Error while parsing JSON %s\n", err)
	}
	actualDockerConfigs, actualBackendAddresses, actualError := parseContainers(jsonConfig["containers"], dockerHost)

	// then
	assertion.AssertDeepEqual("Correct Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Docker Containers Parsed Correctly", testCtx, expectedDockerConfigs, actualDockerConfigs)
	assertion.AssertDeepEqual("Correct Docker Backend Addresses", testCtx, expectedBackendAddresses, actualBackendAddresses)
}

func Test_Parse_Docker_Config_With_Port_Not_Exposed(testCtx *testing.T) {
	// given
	var dockerConfigJson = []byte("{\n" +
			"    \"containers\": [\n" +
			"        {\n" +
			"            \"image\": \"test image\", \n" +
			"            \"portToProxy\": 9090, \n" +
			"            \"portBindings\": {\n" +
			"                \"80\": [\n" +
			"                    {\n" +
			"                        \"hostIp\": \"0.0.0.0\", \n" +
			"                        \"hostPort\": \"1024\"\n" +
			"                    }\n" +
			"                ]\n" +
			"            }\n" +
			"        }\n" +
			"    ]\n" +
			"}")
	var dockerHost = &docker_client.DockerHost{Ip: "127.0.0.1", Port: 1234}
	var expectedDockerConfigs = []*docker_client.DockerConfig{
		&docker_client.DockerConfig{
			Image: "test image",
			PortToProxy: 9090,
			PortBindings: map[docker_client.Port][]docker_client.PortBinding{docker_client.Port("80"): []docker_client.PortBinding{docker_client.PortBinding{HostIp: "0.0.0.0", HostPort: "1024"}}},
		},
	}
	var expectedBackendAddresses = []*contexts.BackendAddress{}
	var expectedError error = errors.New("Invalid container configuration - port specified in \"portToProxy\" must be exposed by container in \"portBindings\"")

	// when
	var jsonConfig = make(map[string]interface{})
	err := json.Unmarshal(dockerConfigJson, &jsonConfig)
	if err != nil {
		testCtx.Fatalf("Error while parsing JSON %s\n", err)
	}
	actualDockerConfigs, actualBackendAddresses, actualError := parseContainers(jsonConfig["containers"], dockerHost)

	// then
	assertion.AssertDeepEqual("Correct Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Docker Containers Parsed Correctly", testCtx, expectedDockerConfigs, actualDockerConfigs)
	assertion.AssertDeepEqual("Correct Docker Backend Addresses", testCtx, expectedBackendAddresses, actualBackendAddresses)
}

func Test_Parse_Docker_Config_With_No_Image(testCtx *testing.T) {
	// given
	var dockerConfigJson = []byte("{\n" +
			"    \"containers\": [\n" +
			"        {\n" +
			"        }\n" +
			"    ]\n" +
			"}")
	var dockerHost = &docker_client.DockerHost{Ip: "127.0.0.1", Port: 1234}
	var expectedDockerConfigs = []*docker_client.DockerConfig{(*docker_client.DockerConfig)(nil)}
	var expectedBackendAddresses = []*contexts.BackendAddress{}
	var expectedError error = errors.New("Invalid container configuration - no \"image\" specified")

	// when
	var jsonConfig = make(map[string]interface{})
	err := json.Unmarshal(dockerConfigJson, &jsonConfig)
	if err != nil {
		testCtx.Fatalf("Error while parsing JSON %s\n", err)
	}
	actualDockerConfigs, actualBackendAddresses, actualError := parseContainers(jsonConfig["containers"], dockerHost)

	// then
	assertion.AssertDeepEqual("Correct Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Docker Containers Parsed Correctly", testCtx, expectedDockerConfigs, actualDockerConfigs)
	assertion.AssertDeepEqual("Correct Docker Backend Addresses", testCtx, expectedBackendAddresses, actualBackendAddresses)
}


func Test_Parse_Docker_Config_With_Multiple_Containers(testCtx *testing.T) {
	// given
	var dockerConfigJson = []byte("{\n" +
			"    \"containers\": [\n" +
			"        {\n" +
			"            \"image\": \"test image one\", \n" +
			"            \"portToProxy\": 9090, \n" +
			"            \"portBindings\": {\n" +
			"                \"80\": [\n" +
			"                    {\n" +
			"                        \"hostIp\": \"0.0.0.0\", \n" +
			"                        \"hostPort\": \"9090\"\n" +
			"                    }\n" +
			"                ]\n" +
			"            }\n" +
			"        },\n" +
			"        {\n" +
			"            \"image\": \"test image two\", \n" +
			"            \"portToProxy\": 9090, \n" +
			"            \"portBindings\": {\n" +
			"                \"80\": [\n" +
			"                    {\n" +
			"                        \"hostIp\": \"0.0.0.0\", \n" +
			"                        \"hostPort\": \"9090\"\n" +
			"                    }\n" +
			"                ]\n" +
			"            }\n" +
			"        }\n" +
			"    ]\n" +
			"}")
	var dockerHost = &docker_client.DockerHost{Ip: "127.0.0.1", Port: 1234}
	var expectedDockerConfigs = []*docker_client.DockerConfig{
		&docker_client.DockerConfig{
			Image: "test image one",
			PortToProxy: 9090,
			PortBindings: map[docker_client.Port][]docker_client.PortBinding{docker_client.Port("80"): []docker_client.PortBinding{docker_client.PortBinding{HostIp: "0.0.0.0", HostPort: "9090"}}},
		},
		&docker_client.DockerConfig{
			Image: "test image two",
			PortToProxy: 9090,
			PortBindings: map[docker_client.Port][]docker_client.PortBinding{docker_client.Port("80"): []docker_client.PortBinding{docker_client.PortBinding{HostIp: "0.0.0.0", HostPort: "9090"}}},
		},
	}
	var containerAddress, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:9090", dockerHost.Ip))
	var expectedBackendAddresses = []*contexts.BackendAddress{
		&contexts.BackendAddress{
			Address: containerAddress,
			Host: dockerHost.Ip,
			Port: "9090",
		},
		&contexts.BackendAddress{
			Address: containerAddress,
			Host: dockerHost.Ip,
			Port: "9090",
		},
	}
	var expectedError error = nil

	// when
	var jsonConfig = make(map[string]interface{})
	err := json.Unmarshal(dockerConfigJson, &jsonConfig)
	if err != nil {
		testCtx.Fatalf("Error while parsing JSON %s\n", err)
	}
	actualDockerConfigs, actualBackendAddresses, actualError := parseContainers(jsonConfig["containers"], dockerHost)

	// then
	assertion.AssertDeepEqual("Correct Correct Error", testCtx, expectedError, actualError)
	assertion.AssertDeepEqual("Correct Docker Containers Parsed Correctly", testCtx, expectedDockerConfigs, actualDockerConfigs)
	assertion.AssertDeepEqual("Correct Docker Backend Addresses", testCtx, expectedBackendAddresses, actualBackendAddresses)
}
