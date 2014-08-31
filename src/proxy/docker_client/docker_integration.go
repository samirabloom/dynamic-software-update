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
)

type DockerClient struct {
	client *docker.Client
}

/*
For this test to pass its need a running instance of docker, for mac use vagrant up docker

sudo vi /etc/default/docker.io
DOCKER_OPTS="-H unix:// -H tcp://0.0.0.0:2375"
 */

func NewDockerClient(endpoint string) (*DockerClient, error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.LoggerFactory().Error("error creating client: %s\n", err)
	}
	_, err = client.Version()
	if err != nil {
		log.LoggerFactory().Error("%s is the server running at %s?\n", err, endpoint)
		return nil, err
	} else {
		return &DockerClient{client: client}, nil
	}
}

func (dc *DockerClient) PullImage(repository, tag string, outputStream io.Writer) (err error) {
	err = dc.client.PullImage(docker.PullImageOptions{Repository: repository, Tag: tag, OutputStream: outputStream}, docker.AuthConfiguration{})
	if err != nil {
		log.LoggerFactory().Error("error pulling image: %s\n", err)
	}
	fmt.Fprintf(outputStream, "Pull Complete for [%s:%s]\n\n", repository, tag)
	return err
}

/*

type Config struct {
	Hostname        string              // Container host name
	Domainname      string
	User            string              // Username or UID
	Memory          int64               // Memory limit (format: <number><optional unit>, where unit = b, k, m or g)
	MemorySwap      int64
	CpuShares       int64               // CPU shares (relative weight)
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	PortSpecs       []string
	ExposedPorts    map[Port]struct{} 	// Expose a port from the container without publishing it to your host
	Tty             bool                // Allocate a pseudo-TTY
	OpenStdin       bool
	StdinOnce       bool
	Environment             []string            // Set environment variables
	Cmd             []string
	Image           string
	Volumes         map[string]struct{} // Bind mount a volume (e.g., from the host: -v /host:/container, from Docker: -v /container)
	VolumesFrom     string              // Mount volumes from the specified container(s)
	WorkingDir      string              // Working directory inside the container
	Entrypoint      []string 			// Overwrite the default ENTRYPOINT of the image
	NetworkDisabled bool
}

type CreateContainerOptions struct {
	Name   string                       // Assign a name to the container
	Config *Config `qs:"-"`
}

type HostConfig struct {
	Binds           []string            // Bind mount a volume (e.g., from the host: -v /host:/container, from Docker: -v /container)
	ContainerIDFile string
	LxcConf         []KeyValuePair      // (lxc exec-driver only) Add custom lxc options --lxc-conf="lxc.cgroup.cpuset.cpus = 0,1"
	Privileged      bool                // Give extended privileges to this container
	PortBindings    map[Port][]PortBinding
	Links           []string            // Add link to another container in the form of name:alias
	PublishAllPorts bool                // Publish all exposed ports to the host interfaces
	Dns             []string            // Set custom DNS servers
	DnsSearch       []string
	VolumesFrom     []string            // Mount volumes from the specified container(s)
	NetworkMode     string              // Set the Network mode for the container
                                               'bridge': creates a new network stack for the container on the docker bridge
                                               'none': no networking for this container
                                               'container:<name|id>': reuses another container network stack
                                               'host': use the host network stack inside the container.  Note: the host mode gives the container full access to local system services such as D-bus and is therefore considered insecure.
	RestartPolicy   RestartPolicy       //  Restart policy to apply when a container exits (no, on-failure, always)
}


{
*	"Image":"",                         // Image for container

	"WorkingDir":"",  					// Working directory inside the container
	"Entrypoint":"",  					// Overwrite the default ENTRYPOINT of the image
	"Environment":null,       					// Set environment variables
	"Cmd":[                             // Set command executed when the container runs
		 ""
	],

	"Hostname":"",   					// Container host name
	"Volumes":{       					// Bind mount a volume (e.g., from the host: -v /host:/container, from Docker: -v /container)
		 "/tmp": {}
	},
	"VolumesFrom":[                 	// Mount volumes from the specified container(s)
		 "parent",
		 "other:ro"
	],
	"ExposedPorts":{  					// Expose a port from the container without publishing it to your host
		 "22/tcp": {}
	},
	"PublishAllPorts":false,            // Publish all exposed ports to the host interfaces
*	"PortBindings":{ "22/tcp": [{ "HostPort": "11022" }] },
*   "PortToProxy":
	"Links":["redis3:redis"],       	// Add link to another container in the form of name:alias

	"User":"",       					// Username or UID
	"Memory":0,      					// Memory limit (format: <number><optional unit>, where unit = b, k, m or g)
	"CpuShares":0                   	// CPU shares (relative weight)
	"LxcConf":{"lxc.utsname":"docker"}  // (lxc exec-driver only) Add custom lxc options --lxc-conf="lxc.cgroup.cpuset.cpus = 0,1"
	"Privileged":false,                 // Give extended privileges to this container
	"CapAdd: [""],             			// Add Linux capabilities
	"CapDrop: [""]                 		// Drop Linux capabilities
}
 */

func (dc *DockerClient) CreateContainer(config *docker.Config, containerName string, outputStream io.Writer) (container *docker.Container, err error) {
	container, err = dc.client.CreateContainer(docker.CreateContainerOptions{Name: containerName, Config: config})
	if err != nil {
		fmt.Fprintf(outputStream, "error creating container: %s\n", err)
		log.LoggerFactory().Error("error creating container: %s\n", err)
	} else {
		fmt.Fprintf(outputStream, "Created container [%s] for image [%s]\n", containerName, config.Image)
	}
	return container, err
}

func (dc *DockerClient) InspectContainer(id string, outputStream io.Writer) (container *docker.Container, err error) {
	container, err = dc.client.InspectContainer(id)
	if err != nil {
		fmt.Fprintf(outputStream, "error inspecting container: %s\n", err)
		log.LoggerFactory().Error("error inspecting cotainer: %s\n", err)
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
		log.LoggerFactory().Error("error starting container: %s\n", err)
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
		log.LoggerFactory().Error("error stopping container: %s\n", err)
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
		log.LoggerFactory().Error("error listing containers: %s\n", err)
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

/*
&docker_client.DockerConfig{
		Image: "string",
		Tag: "string",
		Name: "string",
		WorkingDir: "string",
		Entrypoint: []string{"string", "string"},
		Environment: []string{"string", "string"},
		Cmd: []string{"string", "string"},
		Hostname: "string",
		Volumes: []string{"string", "string"},
		VolumesFrom: []string{"string", "string"},
		PortBindings: map[docker_client.Port][]docker_client.PortBinding{docker_client.Port("8080"): []docker_client.PortBinding{docker_client.PortBinding{HostIp: "127.0.0.1", HostPort: "8080"}}},
		PortToProxy: 0,
		PortSpecs: []string{"string", "string"},
		Links: []string{"string", "string"},
		User: "string",
		Memory: 0,
		CpuShares: 0,
		LxcConf: []docker_client.KeyValuePair{docker_client.KeyValuePair{Key: "key", Value: "value"}},
		Privileged: false,
	}
 */

type DockerConfig struct {
	// name of the image to create the container from
	Image           string `json:"image,omitempty"`
	// the tag for the image, if no value is provided `latest` is used
	Tag             string `json:"tag,omitempty"`
	// always check for new image version
	AlwaysPull      bool `json:"alwaysPull,omitempty"`
	// the port the proxy routes HTTP request to
	PortToProxy     int64 `json:"portToProxy,omitempty"`
	// the name of the container, if no value is provided a unique name is autogenerated
	Name            string `json:"name,omitempty"`
	// sets the processes current working directory inside the container's rootfs
	WorkingDir      string `json:"workingDir,omitempty"`
	// configures a the whole container runs as if it was just that executable, the entrypoint will be executed in a /bin/sh -c
	Entrypoint      []string `json:"entrypoint,omitempty"`
	// a list of environment variables that are set in the running container
	Environment    []string `json:"environment,omitempty"`
	// sets the command to be executed when running the image, the command will be executed in a /bin/sh -c
	Cmd             []string `json:"cmd,omitempty"`
	// Hostname optionally sets the container's hostname if provided
	Hostname        string `json:"hostname,omitempty"`
	// mount volumes from either the host `/host:/container` or from docker `/container`
	Volumes         []string `json:"volumes,omitempty"`
	// mount all the volumes defined for each listed container id or name; an optional suffix of :ro or :rw can be used to mount the volumes in read-only or read-write mode
	VolumesFrom     []string `json:"volumesFrom,omitempty"`
	// bind ports from within
	PortBindings    map[Port][]PortBinding `json:"portBindings,omitempty"`
	// add links to another container in the form of `name:alias`
	Links           []string `json:"links,omitempty"`
	// set the user uid and gid of the executing process running inside the container
	User            string `json:"user,omitempty"`
	// memory limit (in bytes)
	Memory          int64 `json:"memory,omitempty"`
	// cpu shares (relative weight vs. other containers)
	CpuShares       int64 `json:"cpuShares,omitempty"`
	// add custom lxc options as key-value pairs
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
