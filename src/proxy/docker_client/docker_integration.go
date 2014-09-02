package docker_client

import (
	"proxy/log"
	"io"
	docker "github.com/fsouza/go-dockerclient"
	"fmt"
	"code.google.com/p/go-uuid/uuid"
	"strings"
	"net/http"
	"encoding/json"
	"time"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type DockerClient struct {
	client *docker.Client
}

/*
For this test to pass its need a running instance of docker, for mac use vagrant up docker

sudo vi /etc/default/docker.io
DOCKER_OPTS="-H unix:// -H tcp://0.0.0.0:2375"
 */

func timeout(job func() bool) bool {
	timeout := make(chan bool, 1)
	work := make(chan bool, 1)
	go func() {
		time.Sleep(5 * time.Second)
		timeout <- true
	}()
	go func() {
		work <- job()
	}()
	select {
	case <-work:
		// the job completed
		return true
	case <-timeout:
		// the job timed out
		return false
	}
}

func NewDockerClient(endpoint string) (*DockerClient, error) {
	var (
		client *docker.Client
		err error
	)
	success := timeout(func() bool {
		client, err = docker.NewClient(endpoint)
		if err != nil {
			log.LoggerFactory().Error("Error creating client: %s\n", err)
			return false
		}
		_, err = client.Version()
		return true
	})
	if !success || err != nil {
		if err == nil {
			err = errors.New("Client took too long to reach docker host")
		}
		log.LoggerFactory().Error("%s is the server running at %s?\n", err, endpoint)
		return nil, err
	} else {
		return &DockerClient{client: client}, nil
	}
}

func (dc *DockerClient) PullImage(repository, tag string, outputStream io.Writer) (err error) {
	err = dc.client.PullImage(docker.PullImageOptions{Repository: repository, Tag: tag, OutputStream: outputStream}, docker.AuthConfiguration{})
	if err != nil {
		log.LoggerFactory().Error("Error pulling image: %s\n", err)
	}
	fmt.Fprintf(outputStream, "Pull Complete for [%s:%s]\n\n", repository, tag)
	return err
}

func (dc *DockerClient) CreateContainer(config *docker.Config, containerName string, outputStream io.Writer) (container *docker.Container, err error) {
	container, err = dc.client.CreateContainer(docker.CreateContainerOptions{Name: containerName, Config: config})
	if err != nil && err != docker.ErrNoSuchImage {
		fmt.Fprintf(outputStream, "error creating container: %s\n", err)
		log.LoggerFactory().Error("Error creating container: %s\n", err)
	} else {
		fmt.Fprintf(outputStream, "Created container [%s] for image [%s]\n", containerName, config.Image)
	}
	return container, err
}

func (dc *DockerClient) InspectContainer(id string, outputStream io.Writer) (container *docker.Container, err error) {
	container, err = dc.client.InspectContainer(id)
	if err != nil {
		fmt.Fprintf(outputStream, "error inspecting container: %s\n", err)
		log.LoggerFactory().Error("Error inspecting cotainer: %s\n", err)
	} else {
		streamContainer(container, outputStream)
	}
	return container, err
}

func streamContainer(container *docker.Container, outputStream io.Writer) {
	fmt.Fprintf(outputStream, "\n======================================\n")
	fmt.Fprintf(outputStream, "==========CONTAINER DETAILS===========\n")
	fmt.Fprintf(outputStream, "======================================\n")
	fmt.Fprintf(outputStream, ConvertToJson(container))
	fmt.Fprintf(outputStream, "\n======================================\n\n")
}

func ConvertToJson(object interface{}) string {
	json, _ := json.MarshalIndent(object, "", "   ")
	return string(json)
}

func (dc *DockerClient) StartContainer(id string, hostConfig *docker.HostConfig, outputStream io.Writer) (container *docker.Container, err error) {
	err = dc.client.StartContainer(id, hostConfig)
	if err != nil {
		fmt.Fprintf(outputStream, "error starting container: %s\n", err)
		log.LoggerFactory().Error("Error starting container: %s\n", err)
	} else {
		container, err = dc.InspectContainer(id, outputStream)
		time.Sleep(3 * time.Second)
		fmt.Fprintf(outputStream, "Container Log (first 3 seconds):\n")
		attachOpts := docker.AttachToContainerOptions{
			Container:    container.ID,
			OutputStream: outputStream,
			ErrorStream:  outputStream,
			Stdout:       true,
			Stderr:       true,
			Logs:         true,
			Stream:          false,
		}
		dc.client.AttachToContainer(attachOpts)
		fmt.Fprintf(outputStream, "\n")
	}
	return container, err
}

func (dc *DockerClient) StopContainer(id string, timeout uint, outputStream io.Writer) (string, error) {
	err := dc.client.StopContainer(id, timeout)
	if err != nil {
		fmt.Fprintf(outputStream, "error stopping container: %s\n", err)
		log.LoggerFactory().Error("Error stopping container: %s\n", err)
		return "", err
	} else {
		fmt.Fprintf(outputStream, "Stopped container id: %s\n", id)
		return id, nil
	}
}

func (dc *DockerClient) RemoveContainer(name string, timeout uint, outputStream io.Writer) (err error) {
	var (
		containerList []docker.APIContainers
		container docker.APIContainers
	)
	containerList, err = dc.client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fmt.Fprintf(outputStream, "error listing containers: %s\n", err)
		log.LoggerFactory().Error("Error listing containers: %s\n", err)
	} else {
		searchName := "/" + strings.Replace(name, ":", "/", 2)
		for _, container = range containerList {
			for _, containerName := range container.Names {
				if containerName == searchName || containerName == name {
					fmt.Fprintf(outputStream, "Stopping container id: %s name: \"%s\"\n", container.ID, containerName)
					dc.StopContainer(container.ID, timeout, outputStream)
					fmt.Fprintf(outputStream, "Removing container id: %s name: \"%s\"\n", container.ID, containerName)
					dc.client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true})
				}
			}
		}
	}
	return err
}

func GenerateRandomName(prefix string, size int) (string, error) {
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return "", err
	}
	return prefix+hex.EncodeToString(id)[:size], nil
}

type DockerHost struct {
	Ip   string `json:"ip,omitempty"`
	Port int `json:"port,omitempty"`
	Log  bool `json:"log,omitempty"`
}

func (dh *DockerHost) Endpoint() string {
	return fmt.Sprintf("http://%s:%d", dh.Ip, dh.Port)
}

type DockerConfig struct {
	// name of the image to create each container from
	Image           string `json:"image,omitempty"`
	// the tag for the image, defaulting to `latest`
	Tag             string `json:"tag,omitempty"`
	// override the docker host for this container
	DockerHost	    *DockerHost `json:"dockerHost,omitempty"`
	// whether to automatically check for new image versions
	AlwaysPull      bool `json:"alwaysPull,omitempty"`
	// the port the proxy routes HTTP request to
	PortToProxy     int64 `json:"portToProxy,omitempty"`
	// the name of the container, defaulting to an auto-generated unique name
	Name            string `json:"name,omitempty"`
	// the current working directory inside the container's root file system
	WorkingDir      string `json:"workingDir,omitempty"`
	// the container's entrypoint
	Entrypoint      []string `json:"entrypoint,omitempty"`
	// the environment variables that are set in the running container
	Environment     []string `json:"environment,omitempty"`
	// the command to be executed when running the container
	Cmd             []string `json:"cmd,omitempty"`
	// the container's hostname
	Hostname        string `json:"hostname,omitempty"`
	// mount volumes from either the host or from another docker container
	Volumes         []string `json:"volumes,omitempty"`
	// mount all the volumes defined for one or more other docker containers, including control over whether the volumes are mounted in read-write or read-only mode
	VolumesFrom     []string `json:"volumesFrom,omitempty"`
	// bind ports from within the container to one or more host port and IP combinations
	PortBindings    map[Port][]PortBinding `json:"portBindings,omitempty"`
	// link the container to another Docker container
	Links           []string `json:"links,omitempty"`
	// set the user user id and group id of the executing process running inside the container
	User            string `json:"user,omitempty"`
	// restrict the container's memory (in bytes)
	Memory          int64 `json:"memory,omitempty"`
	// restrict the container's cpu share using relative weight compared to other containers
	CpuShares       int64 `json:"cpuShares,omitempty"`
	// add custom LXC options for the container
	LxcConf         []KeyValuePair `json:"lxcConf,omitempty"`
	// give extended privileges to the container
	Privileged      bool `json:"privileged,omitempty"`
}

func (dc *DockerConfig) HasPortExposed(portToTest string) bool {
	for _, hostBindings := range dc.PortBindings {
		for _, hostBinding := range hostBindings {
			if hostBinding.HostPort == portToTest {
				return true;
			}
		}
	}
	return false;
}

func (dc *DockerConfig) String() string {
	data, _ := json.MarshalIndent(dc, "", "   ")
	return string(data)
}

type Port string

type PortBinding struct {
	// host ip address used in port binding
	HostIp   string `json:"hostIp,omitempty"`
	// host port used in port binding
	HostPort string `json:"hostPort,omitempty"`
}


type KeyValuePair struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func Flush(outputStream io.Writer) {
	flusher, isFlusher := outputStream.(http.Flusher)
	if isFlusher {
		flusher.Flush()
	}
}

func (dc *DockerClient) CreateServerFromContainer(config *DockerConfig, outputStream io.Writer) (container *docker.Container, err error) {
	if config.AlwaysPull {
		err = dc.PullImage(config.Image, config.Tag, outputStream)
		Flush(outputStream)
	}
	if err == nil {

		dockerConfig := &docker.Config{
			Image: fmt.Sprintf("%s:%s", config.Image, config.Tag),
			Hostname: config.Hostname,
			User: config.User,
			Memory: config.Memory,
			CpuShares: config.CpuShares,
			Env: config.Environment,
			Cmd: config.Cmd,
			WorkingDir: config.WorkingDir,
			Entrypoint: config.Entrypoint,
		}

		containerName := config.Name
		if len(containerName) == 0 {
			containerName = strings.Replace(config.Image, "/", "_", 2)+"_"+uuid.NewUUID().String()
		} else {
			fmt.Fprintf(outputStream, "Checking for any existing containers with name \"%s\"\n", containerName)
			dc.RemoveContainer(containerName, 60, outputStream)
			Flush(outputStream)
		}
		container, err = dc.CreateContainer(dockerConfig, containerName, outputStream)
		Flush(outputStream)
		if err == docker.ErrNoSuchImage {
			err = dc.PullImage(config.Image, config.Tag, outputStream)
			Flush(outputStream)
			container, err = dc.CreateContainer(dockerConfig, containerName, outputStream)
			if err == docker.ErrNoSuchImage {
				fmt.Fprintf(outputStream, "error creating container: %s\n", err)
				log.LoggerFactory().Error("Error creating container: %s\n", err)
			}
			Flush(outputStream)
		}
		if err == nil {
			dockerHostConfig := &docker.HostConfig{
				Binds:           config.Volumes,
				Privileged:      config.Privileged,
				Links:           config.Links,
				VolumesFrom:     config.VolumesFrom,
			}
			dockerHostConfig.PortBindings = make(map[docker.Port][]docker.PortBinding, len(config.PortBindings))
			for port := range config.PortBindings {
				dockerHostConfig.PortBindings[docker.Port(port)] = make([]docker.PortBinding, len(config.PortBindings[port]))
				for index := range config.PortBindings[port] {
					dockerHostConfig.PortBindings[docker.Port(port)][index] = docker.PortBinding{
						HostIp: config.PortBindings[port][index].HostIp,
						HostPort: config.PortBindings[port][index].HostPort,
					}
				}
			}
			dockerHostConfig.LxcConf = make([]docker.KeyValuePair, len(config.LxcConf))
			for index := range config.LxcConf {
				dockerHostConfig.LxcConf[index] = docker.KeyValuePair{
					Key: config.LxcConf[index].Key,
					Value: config.LxcConf[index].Value,
				}
			}

			container, err = dc.StartContainer(container.ID, dockerHostConfig, outputStream)
			Flush(outputStream)
		}
	}

	return container, err
}
