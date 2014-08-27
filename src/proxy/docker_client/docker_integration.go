package docker_client

import (
	"proxy/log"
	"io"
	"github.com/fsouza/go-dockerclient/docker"
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
	Env             []string            // Set environment variables
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
	"Env":null,       					// Set environment variables
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

func (dc *DockerClient) CreateContainer(imageName, containerName string) (container *docker.Container, err error) {
	config := docker.Config{Image: imageName, AttachStdout: true, AttachStdin: true}
	opts := docker.CreateContainerOptions{Name: containerName, Config: &config}

	container, err = dc.client.CreateContainer(opts)
	if err != nil {
		log.LoggerFactory().Error("error creating container: %s\n", err)
	}
	return container, err
}

func (dc *DockerClient) InspectContainer(id string) (container *docker.Container, err error) {
	container, err = dc.client.InspectContainer(id)
	if err != nil {
		log.LoggerFactory().Error("error inspecting cotainer: %s\n", err)
	}
	return container, err
}

func (dc *DockerClient) StartContainer(id string) (container *docker.Container, err error) {
	err = dc.client.StartContainer(id, &docker.HostConfig{})
	if err != nil {
		log.LoggerFactory().Error("error starting container: %s\n", err)
	} else {
		container, err = dc.InspectContainer(id)
	}
	return container, err
}

func (dc *DockerClient) StopContainer(id string, timeout uint) (container *docker.Container, err error) {
	err = dc.client.StopContainer(id, timeout)
	if err != nil {
		log.LoggerFactory().Error("error starting container: %s\n", err)
	} else {
		container, err = dc.InspectContainer(id)
	}
	return container, err
}
