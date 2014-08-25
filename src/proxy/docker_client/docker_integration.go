package docker_client

import (
	"proxy/log"
	docker "github.com/fsouza/go-dockerclient"
	"io"
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
