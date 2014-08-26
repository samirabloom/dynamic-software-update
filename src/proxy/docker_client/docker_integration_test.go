package docker_client

import (
	"testing"
	"bytes"
	"strings"
	docker "github.com/fsouza/go-dockerclient"
	"time"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"os/exec"
	"os"
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

func createVagrantDockerBox() {
	vagrantCommand := exec.Command("/usr/bin/vagrant", "up", "docker")
	fmt.Println("===================================")
	fmt.Println("Launching Vagrant Docker Ubuntu Box")
	fmt.Println("===================================")
	vagrantCommand.Stderr = os.Stderr
	vagrantCommand.Stdout = os.Stdout
	err := vagrantCommand.Run()
	if err != nil {
		fmt.Printf("error occured %s\n", err)
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
	createVagrantDockerBox()
	var (
		endpoint = "http://192.168.50.5:2375"
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
	createVagrantDockerBox()
	var (
		endpoint = "http://192.168.50.5:2375"
		imageName     = "samirabloom/docker-go"
		containerName = strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String()
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}

	container, err := client.CreateContainer(imageName, containerName)

	// then
	if container.Image == imageName {
		testCtx.Fatalf("Container does not have the correct image name, expect: [%s], found: [%s]\n", imageName, container.Image)
	}
	if container.State.Running {
		testCtx.Fatalf("Container was running, expect: [%t], found: [%t]\n", false, container.State.Running)
	}
	if err != nil {
		testCtx.Fatalf("Error while creating container %s\n", err)
	}
}

func Test_Docker_Client_Should_Inspect_Container(testCtx *testing.T) {
	// given
	createVagrantDockerBox()
	var (
		endpoint = "http://192.168.50.5:2375"
		imageName     = "samirabloom/docker-go"
		containerName = strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String()
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}

	createdContainer, err := client.CreateContainer(imageName, containerName)
	if err != nil {
		testCtx.Fatalf("Error while creating container %s\n", err)
	}

	container, err := client.InspectContainer(createdContainer.ID)

	// then
	if container.Image == imageName {
		testCtx.Fatalf("Container does not have the correct image name, expect: [%s], found: [%s]\n", imageName, container.Image)
	}
	if container.State.Running {
		testCtx.Fatalf("Container was running, expect: [%t], found: [%t]\n", false, container.State.Running)
	}
	if err != nil {
		testCtx.Fatalf("Error while inspecting container %s\n", err)
	}
}

func Test_Docker_Client_Should_Start_Container(testCtx *testing.T) {
	// given
	createVagrantDockerBox()
	var (
		endpoint = "http://192.168.50.5:2375"
		imageName                           = "samirabloom/docker-go"
		containerName                       = strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String()
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}

	createdContainer, err := client.CreateContainer(imageName, containerName)
	if err != nil {
		testCtx.Fatalf("Error while creating container %s\n", err)
	}

	container, err := client.StartContainer(createdContainer.ID)

	// then
	if container.Image == imageName {
		testCtx.Fatalf("Container does not have the correct image name, expect: [%s], found: [%s]\n", imageName, container.Image)
	}
	if !container.State.Running {
		testCtx.Fatalf("Container was not running, expect: [%t], found: [%t]\n", true, container.State.Running)
	}
	if err != nil {
		testCtx.Fatalf("Error while starting container %s\n", err)
	}
}

func Test_Docker_Client_Should_Stop_Container(testCtx *testing.T) {
	// given
	createVagrantDockerBox()
	var (
		endpoint = "http://192.168.50.5:2375"
		imageName                           = "samirabloom/docker-go"
		containerName                       = strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String()
	)

	// when
	client, err := NewDockerClient(endpoint)
	if err != nil {
		testCtx.Fatalf("Error while creating client %s\n", err)
	}

	createdContainer, err := client.CreateContainer(imageName, containerName)
	if err != nil {
		testCtx.Fatalf("Error while creating container %s\n", err)
	}

	container, err := client.StopContainer(createdContainer.ID, uint(10))

	// then
	if container.Image == imageName {
		testCtx.Fatalf("Container does not have the correct image name, expect: [%s], found: [%s]\n", imageName, container.Image)
	}
	if container.State.Running {
		testCtx.Fatalf("Container was not running, expect: [%t], found: [%t]\n", true, container.State.Running)
	}
	if err != nil {
		testCtx.Fatalf("Error while stopping container %s\n", err)
	}
}