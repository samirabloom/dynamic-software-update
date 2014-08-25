package main

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"code.google.com/p/go-uuid/uuid"
	"strings"
	"os"
)

func main() {

	incorrectEndPoint := "http://192.168.50.7:2377"
	clientError, err := docker.NewClient(incorrectEndPoint)
	if err != nil {
		fmt.Printf("error creating client: %s\n", err)
	}
	_, err = clientError.Version()
	if err != nil {
		fmt.Printf("%s is the server running at %s\n\n", err, incorrectEndPoint)
	}

	endpoint := "http://192.168.50.7:2375"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		fmt.Printf("error creating client: %s\n", err)
	}
	_, err = client.Version()
	if err != nil {
		fmt.Printf("%s is the server running at %s\n\n", err, endpoint)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		fmt.Printf("error list containers: %s\n", err)
	} else {
		for _, container := range containers {
			fmt.Println("===============CONTANER===============")
			fmt.Println("======================================")
			fmt.Println("ID: ", container.ID)
			fmt.Println("Image: ", container.Image)
			fmt.Println("Command: ", container.Command)
			fmt.Println("Created: ", container.Created)
			fmt.Println("Ports: ", container.Ports)
			fmt.Println("Status: ", container.Status)
			fmt.Println("======================================\n")
		}
	}

	fmt.Println("==============PULL IMAGE==============")
	fmt.Println("======================================")
	err = client.PullImage(docker.PullImageOptions{Repository: "jamesdbloom/couchbase", Tag: "latest", OutputStream: os.Stdout}, docker.AuthConfiguration{})
	if err != nil {
		fmt.Printf("error pulling image: %s\n", err)
	}
	fmt.Println("======================================\n")

	imageName := "jamesdbloom/couchbase"
	config := docker.Config{Image: imageName, AttachStdout: true, AttachStdin: true}
	opts := docker.CreateContainerOptions{Name: strings.Replace(imageName, "/", "_", 2) + "_" + uuid.NewUUID().String(), Config: &config}
	container, err := client.CreateContainer(opts)
	if err != nil {
		fmt.Printf("error creating container: %s\n", err)
	} else {
		fmt.Println("============CREATE CONTIANER==========")
		fmt.Println("======================================")
		fmt.Println("ID: ", container.ID)
		fmt.Println("Created: ", container.Created)
		fmt.Println("Path: ", container.Path)
		fmt.Println("Args: ", container.Args)
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
		fmt.Println("Image: ", container.Image)
		fmt.Println("NetworkSettings: ", container.NetworkSettings)
		fmt.Println("SysInitPath: ", container.SysInitPath)
		fmt.Println("ResolvConfPath: ", container.ResolvConfPath)
		fmt.Println("HostnamePath: ", container.HostnamePath)
		fmt.Println("HostsPath: ", container.HostsPath)
		fmt.Println("Name: ", container.Name)
		fmt.Println("Driver: ", container.Driver)
		fmt.Println("Volumes: ", container.Volumes)
		fmt.Println("VolumesRW: ", container.VolumesRW)
		fmt.Println("HostConfig: ", container.HostConfig)
		fmt.Println("======================================\n")
	}

	fmt.Println("============START CONTAINER===========")
	fmt.Println("======================================")
	err = client.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		fmt.Printf("error starting container: %s\n", err)
	}
	container, err = client.InspectContainer(container.ID)
	if err != nil {
		fmt.Printf("error inspecting cotainer: %s\n", err)
	} else {
		fmt.Println("============STARTED IMAGE=============")
		fmt.Println("======================================")
		fmt.Println("ID: ", container.ID)
		fmt.Println("Created: ", container.Created)
		fmt.Println("Path: ", container.Path)
		fmt.Println("Args: ", container.Args)
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
		fmt.Println("Image: ", container.Image)
		fmt.Println("NetworkSettings: ", container.NetworkSettings)
		fmt.Println("SysInitPath: ", container.SysInitPath)
		fmt.Println("ResolvConfPath: ", container.ResolvConfPath)
		fmt.Println("HostnamePath: ", container.HostnamePath)
		fmt.Println("HostsPath: ", container.HostsPath)
		fmt.Println("Name: ", container.Name)
		fmt.Println("Driver: ", container.Driver)
		fmt.Println("Volumes: ", container.Volumes)
		fmt.Println("VolumesRW: ", container.VolumesRW)
		fmt.Println("HostConfig: ", container.HostConfig)
		fmt.Println("======================================\n")
	}
	fmt.Println("======================================\n")

	fmt.Println("=============STOP CONTAINER===========")
	fmt.Println("======================================")
	err = client.StopContainer(container.ID, uint(10))
	if err != nil {
		fmt.Printf("error stoping container: %s\n", err)
	}
	container, err = client.InspectContainer(container.ID)
	if err != nil {
		fmt.Printf("error inspecting cotainer: %s\n", err)
	} else {
		fmt.Println("============STOPPED IMAGE=============")
		fmt.Println("======================================")
		fmt.Println("ID: ", container.ID)
		fmt.Println("Created: ", container.Created)
		fmt.Println("Path: ", container.Path)
		fmt.Println("Args: ", container.Args)
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
		fmt.Println("Image: ", container.Image)
		fmt.Println("NetworkSettings: ", container.NetworkSettings)
		fmt.Println("SysInitPath: ", container.SysInitPath)
		fmt.Println("ResolvConfPath: ", container.ResolvConfPath)
		fmt.Println("HostnamePath: ", container.HostnamePath)
		fmt.Println("HostsPath: ", container.HostsPath)
		fmt.Println("Name: ", container.Name)
		fmt.Println("Driver: ", container.Driver)
		fmt.Println("Volumes: ", container.Volumes)
		fmt.Println("VolumesRW: ", container.VolumesRW)
		fmt.Println("HostConfig: ", container.HostConfig)
		fmt.Println("======================================\n")
	}
	fmt.Println("======================================\n")
}

