package docker_client

import (
	"testing"
	"bytes"
	"strings"
	docker "github.com/fsouza/go-dockerclient"
	"time"
	"code.google.com/p/go-uuid/uuid"
	"util/test/vagrant"
)

func createContainer(containerName string) *docker.Container {
	return &docker.Container{
		Created: time.Time{},
		Path: "",
		Args: nil,
		Config: nil,
		State:  docker.State{Running:false, Paused:false, Pid:0, ExitCode:0, StartedAt:time.Time{}, FinishedAt:time.Time{}},
		Image:  "",
		NetworkSettings: nil,
		SysInitPath: "",
		ResolvConfPath: "",
		HostnamePath: "",
		HostsPath: "",
		Name: containerName,
		Driver: "",
		Volumes: nil,
		VolumesRW: nil,
		HostConfig: nil,
	}
}


func Test_Docker_Client_Should_Indicate_Docker_Not_Available(testCtx *testing.T) {
	// given
	var (
		incorrectEndpoint = "http://127.0.0.1:666"
	)

	// when
	client, err := NewDockerClient(incorrectEndpoint)

	// then
	if err == nil {
		testCtx.Fatalf("Expected failure while creating client %s\n", err)
	}
	if client != nil {
		testCtx.Fatalf("Expected client not nil, expected: %v, actual: %v\n", nil, client)
	}
}

func Test_Docker_Client_Should_Pull_Container(testCtx *testing.T) {
	// given
	vagrant.CreateVagrantDockerBox()
	var (
		endpoint   = "http://192.168.50.5:2375"
		repository = "samirabloom/docker-go"
		tag        = "latest"
		actualBuffer bytes.Buffer
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}

	client.PullImage(repository, tag, &actualBuffer)

	// then
	dockerMessage := actualBuffer.String()
	if !strings.Contains(dockerMessage, "Pulling repository "+repository) {
		testCtx.Fatalf("Docker message did not contain \"Pulling repository %s\" message was as follows:\n%s\n", repository, dockerMessage)
	}
	if !strings.Contains(dockerMessage, "Pulling image ("+tag+") from "+repository) {
		testCtx.Fatalf("Docker message did not contain \"Pulling image (%s) from %s\" message was as follows:\n%s\n", tag, repository, dockerMessage)
	}
}

func Test_Docker_Client_Should_Create_Container(testCtx *testing.T) {
	// given
	vagrant.CreateVagrantDockerBox()
	var (
		endpoint      = "http://192.168.50.5:2375"
		imageName     = "samirabloom/docker-go"
		containerName = strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String()
		outputStream bytes.Buffer
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}
	config := &docker.Config{Image: imageName, AttachStdout: true, AttachStdin: true}

	createdContainer, err := client.CreateContainer(config, containerName, &outputStream)

	// then
	if err != nil {
		testCtx.Fatalf("Error while creating container %s\n", err)
	}
	if createdContainer == nil {
		testCtx.Fatalf("Container not created\n", imageName, createdContainer.Image)
	}
	if createdContainer.Image == imageName {
		testCtx.Fatalf("Container does not have the correct image name, expected: [%s], found: [%s]\n", imageName, createdContainer.Image)
	}
	if createdContainer.State.Running {
		testCtx.Fatalf("Container was running, expected: [%t], found: [%t]\n", false, createdContainer.State.Running)
	}
}

func Test_Docker_Client_Should_Inspect_Container(testCtx *testing.T) {
	// given
	vagrant.CreateVagrantDockerBox()
	var (
		endpoint      = "http://192.168.50.5:2375"
		imageName     = "samirabloom/docker-go"
		containerName = strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String()
		outputStream bytes.Buffer
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}
	config := &docker.Config{Image: imageName, AttachStdout: true, AttachStdin: true}

	createdContainer, err := client.CreateContainer(config, containerName, &outputStream)
	if err != nil {
		testCtx.Fatalf("Error while creating container %s\n", err)
	}

	container, err := client.InspectContainer(createdContainer.ID, &outputStream)

	// then
	if err != nil {
		testCtx.Fatalf("Error while inspecting container %s\n", err)
	}
	if container == nil {
		testCtx.Fatalf("Container not found\n", imageName, container.Image)
	}
	if container.Image == imageName {
		testCtx.Fatalf("Container does not have the correct image name, expected: [%s], found: [%s]\n", imageName, container.Image)
	}
	if container.State.Running {
		testCtx.Fatalf("Container was running, expected: [%t], found: [%t]\n", false, container.State.Running)
	}
}

func Test_Docker_Client_Should_Start_Container(testCtx *testing.T) {
	// given
	vagrant.CreateVagrantDockerBox()
	var (
		endpoint      = "http://192.168.50.5:2375"
		imageName     = "samirabloom/docker-go"
		containerName = strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String()
		outputStream bytes.Buffer
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}
	config := &docker.Config{Image: imageName, AttachStdout: true, AttachStdin: true}

	createdContainer, err := client.CreateContainer(config, containerName, &outputStream)
	if err != nil {
		testCtx.Fatalf("Error while creating container %s\n", err)
	}

	container, err := client.StartContainer(createdContainer.ID, &docker.HostConfig{}, &outputStream)

	// then
	if err != nil {
		testCtx.Fatalf("Error while starting container %s\n", err)
	}
	if container == nil {
		testCtx.Fatalf("Container not started\n", imageName, container.Image)
	}
	if container.Image == imageName {
		testCtx.Fatalf("Container does not have the correct image name, expected: [%s], found: [%s]\n", imageName, container.Image)
	}
	if !container.State.Running {
		testCtx.Fatalf("Container was not running, expected: [%t], found: [%t]\n", true, container.State.Running)
	}
}

func Test_Docker_Client_Should_Stop_Container(testCtx *testing.T) {
	// given
	vagrant.CreateVagrantDockerBox()
	var (
		endpoint      = "http://192.168.50.5:2375"
		imageName     = "samirabloom/docker-go"
		containerName = strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String()
		outputStream bytes.Buffer
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}
	config := &docker.Config{Image: imageName, AttachStdout: true, AttachStdin: true}

	createdContainer, err := client.CreateContainer(config, containerName, &outputStream)
	if err != nil {
		testCtx.Fatalf("Error while creating container %s\n", err)
	}

	id, err := client.StopContainer(createdContainer.ID, uint(10), &outputStream)

	// then
	if id != createdContainer.ID {
		testCtx.Fatalf("Incorrect stopped container id, expected: [%s], found: [%s]\n", createdContainer.ID, id)
	}

	// and
	container, err := client.InspectContainer(createdContainer.ID, &outputStream)
	if err != nil {
		testCtx.Fatalf("Error while stopping container %s\n", err)
	}
	if container == nil {
		testCtx.Fatalf("Container not found\n", imageName, container.Image)
	}
	if container.Image == imageName {
		testCtx.Fatalf("Container does not have the correct image name, expected: [%s], found: [%s]\n", imageName, container.Image)
	}
	if container.State.Running {
		testCtx.Fatalf("Container was not running, expected: [%t], found: [%t]\n", true, container.State.Running)
	}
}

func Test_Docker_Client_Should_Remove_Container(testCtx *testing.T) {
	// given
	vagrant.CreateVagrantDockerBox()
	var (
		endpoint      = "http://192.168.50.5:2375"
		imageName     = "samirabloom/docker-go"
		containerName = strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String() // "samirabloom_docker-go"
		outputStream bytes.Buffer
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}
	config := &docker.Config{Image: imageName, AttachStdout: true, AttachStdin: true}

	createdContainer, err := client.CreateContainer(config, containerName, &outputStream)
	if err != nil {
		testCtx.Fatalf("Error while creating container %s\n", err)
	}

	err = client.RemoveContainer(containerName, uint(10), &outputStream)

	// then
	if err != nil {
		testCtx.Fatalf("Error returned, expected: [%s], found: [%s]\n", (error)(nil), err)
	}

	// and
	container, err := client.InspectContainer(createdContainer.ID, &outputStream)
	expectedErrorMessage := "No such container: " + createdContainer.ID
	if err == nil || err.Error() != expectedErrorMessage {
		testCtx.Fatalf("Wrong error while stopping container, expected: [%s], found: [%s]\n", expectedErrorMessage, err)
	}
	if container != nil {
		testCtx.Fatalf("Container still exists, expected: [%s], found: [%s]\n", (*docker.Container)(nil), container)
	}
}
