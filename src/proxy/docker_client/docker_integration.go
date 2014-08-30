package docker_client

import (
	"proxy/log"
	"io"
	docker "github.com/fsouza/go-dockerclient"
	"fmt"
	"code.google.com/p/go-uuid/uuid"
	"strings"
	"net/http"
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
	fmt.Fprintf(outputStream, "ID: %s\n", container.ID)
	fmt.Fprintf(outputStream, "Created: %s\n", container.Created)
	fmt.Fprintf(outputStream, "Path: %s\n", container.Path)
	fmt.Fprintf(outputStream, "Args: %s\n", container.Args)
	if container.Config != nil {
		fmt.Printf("Config: -- \n\t Hostname: %v,\n\t " +
					"Domainname: %v,\n\t " +
					"User: %v,\n\t " +
					"Memory: %v,\n\t " +
					"MemorySwap: %v,\n\t " +
					"CpuShares: %v,\n\t " +
					"AttachStdin: %v,\n\t " +
					"AttachStdout: %v,\n\t " +
					"AttachStderr: %v,\n\t " +
					"PortSpecs: %v,\n\t " +
					"ExposedPorts: %v,\n\t " +
					"Tty: %v,\n\t " +
					"OpenStdin: %v,\n\t " +
					"StdinOnce: %v,\n\t " +
					"Env: %v,\n\t " +
					"Cmd: %v,\n\t " +
					"Dns: %v,\n\t " +
					"Image: %v,\n\t " +
					"Volumes: %v,\n\t " +
					"VolumesFrom: %v,\n\t " +
					"WorkingDir: %v,\n\t " +
					"Entrypoint: %v,\n\t " +
					"NetworkDisabled: %v\n",
			container.Config.Hostname,
			container.Config.Domainname,
			container.Config.User,
			container.Config.Memory,
			container.Config.MemorySwap,
			container.Config.CpuShares,
			container.Config.AttachStdin,
			container.Config.AttachStdout,
			container.Config.AttachStderr,
			container.Config.PortSpecs,
			container.Config.ExposedPorts,
			container.Config.Tty,
			container.Config.OpenStdin,
			container.Config.StdinOnce,
			container.Config.Env,
			container.Config.Cmd,
			container.Config.Dns,
			container.Config.Image,
			container.Config.Volumes,
			container.Config.VolumesFrom,
			container.Config.WorkingDir,
			container.Config.Entrypoint,
			container.Config.NetworkDisabled)
	}
	fmt.Printf("State: -- \n\t Running: %t, \n\t " +
				"Paused: %t, \n\t " +
				"Pid: %d, \n\t " +
				"ExitCode: %d, \n\t " +
				"StartedAt: %v, \n\t " +
				"FinishedAt: %v\n",
		container.State.Running,
		container.State.Paused,
		container.State.Pid,
		container.State.ExitCode,
		container.State.StartedAt,
		container.State.FinishedAt)
	fmt.Fprintf(outputStream, "Image: %s\n", container.Image)
	fmt.Fprintf(outputStream, "NetworkSettings: %s\n", container.NetworkSettings)
	fmt.Fprintf(outputStream, "SysInitPath: %s\n", container.SysInitPath)
	fmt.Fprintf(outputStream, "ResolvConfPath: %s\n", container.ResolvConfPath)
	fmt.Fprintf(outputStream, "HostnamePath: %s\n", container.HostnamePath)
	fmt.Fprintf(outputStream, "HostsPath: %s\n", container.HostsPath)
	fmt.Fprintf(outputStream, "Name: %s\n", container.Name)
	fmt.Fprintf(outputStream, "Driver: %s\n", container.Driver)
	fmt.Fprintf(outputStream, "Volumes: %s\n", container.Volumes)
	fmt.Fprintf(outputStream, "VolumesRW: %s\n", container.VolumesRW)
	fmt.Fprintf(outputStream, "HostConfig: %s\n", container.HostConfig)
	fmt.Fprintf(outputStream, "======================================\n\n")
}

func (dc *DockerClient) StartContainer(id string, hostConfig *docker.HostConfig, outputStream io.Writer) (container *docker.Container, err error) {
	err = dc.client.StartContainer(id, hostConfig)
	if err != nil {
		fmt.Fprintf(outputStream, "error starting container: %s\n", err)
		log.LoggerFactory().Error("error starting container: %s\n", err)
	} else {
		container, err = dc.InspectContainer(id, outputStream)
	}
	return container, err
}

func (dc *DockerClient) StopContainer(id string, timeout uint, outputStream io.Writer) (container *docker.Container, err error) {
	err = dc.client.StopContainer(id, timeout)
	if err != nil {
		fmt.Fprintf(outputStream, "error stopping container: %s\n", err)
		log.LoggerFactory().Error("error stopping container: %s\n", err)
	} else {
		container, err = dc.InspectContainer(id, outputStream)
	}
	return container, err
}

type DockerConfig struct {
	Image           string
	Tag             string
	Name            string
	WorkingDir      string
	Entrypoint      []string
	Env             []string
	Cmd             []string
	Hostname        string
	Volumes         []string
	VolumesFrom     []string
	ExposedPorts    map[string]struct{}
	PublishAllPorts bool
	PortBindings    map[string][]map[string]string
	PortToProxy     int64
	Links           []string
	User            string
	Memory          int64
	CpuShares       int64
	LxcConf         []KeyValuePair
	Privileged      bool
}

type KeyValuePair struct {
	Key   string
	Value string
}

func Flush(outputStream io.Writer) {
	flusher, isFlusher := outputStream.(http.Flusher)
	if isFlusher {
		flusher.Flush()
	}
}

func (dc *DockerClient) CreateServerFromContainer(config *DockerConfig, outputStream io.Writer) (container *docker.Container, err error) {
	err = dc.PullImage(config.Image, config.Tag, outputStream)
	fmt.Printf("outputStream: %#v\n", outputStream)
	Flush(outputStream)
	if err == nil {

		dockerConfig := &docker.Config{
			Image: fmt.Sprintf("%s:%s", config.Image, config.Tag),
			Hostname: config.Hostname,
			User: config.User,
			Memory: config.Memory,
			CpuShares: config.CpuShares,
			Env: config.Env,
			Cmd: config.Cmd,
			WorkingDir: config.WorkingDir,
			Entrypoint: config.Entrypoint,
		}
		dockerConfig.ExposedPorts = make(map[docker.Port]struct{}, len(config.ExposedPorts))
		for index := range config.ExposedPorts {
			dockerConfig.ExposedPorts[docker.Port(index)] = config.ExposedPorts[index]
		}

		containerName := config.Name
		if len(containerName) == 0 {
			containerName = strings.Replace(config.Image, "/", "_", 2)+"_"+uuid.NewUUID().String()
		}
		container, err = dc.CreateContainer(dockerConfig, containerName, outputStream)
		Flush(outputStream)
		if err == nil {
			dockerHostConfig := &docker.HostConfig{
				Binds:           config.Volumes,
				Privileged:      config.Privileged,
				Links:           config.Links,
				PublishAllPorts: config.PublishAllPorts,
				VolumesFrom:     config.VolumesFrom,
			}
			dockerHostConfig.PortBindings = make(map[docker.Port][]docker.PortBinding, len(config.PortBindings))
			for port := range config.PortBindings {
				dockerHostConfig.PortBindings[docker.Port(port)] = make([]docker.PortBinding, len(config.PortBindings[port]))
				for index := range config.PortBindings[port] {
					dockerHostConfig.PortBindings[docker.Port(port)][index] = docker.PortBinding{
						HostIp: config.PortBindings[port][index]["HostIp"],
						HostPort: config.PortBindings[port][index]["HostPort"],
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
