package contexts

import (
	"container/list"
	"code.google.com/p/go-uuid/uuid"
	"net"
	"fmt"
	"proxy/log"
	"proxy/tcp"
	"errors"
	"proxy/docker_client"
	"io"
)

type Clusters  struct {
	ContextsByVersion    *list.List
	ContextsByID         map[string]*Cluster
	DockerHostEndpoint   string
}

func (clusters *Clusters) Add(cluster *Cluster) {
	if clusters.ContextsByVersion == nil {
		clusters.ContextsByVersion = list.New()
	}
	if clusters.ContextsByID == nil {
		clusters.ContextsByID = make(map[string]*Cluster)
	}
	clusterToAdd := clusters.ContextsByID[cluster.Uuid.String()]
	if clusterToAdd == nil {
		insertOrderedByVersion(clusters.ContextsByVersion, cluster)
		clusters.ContextsByID[cluster.Uuid.String()] = cluster
	}
}

func insertOrderedByVersion(orderedList *list.List, cluster *Cluster) {
	if orderedList.Front() == nil {
		orderedList.PushFront(cluster)
	} else {
		inserted := false
		for element := orderedList.Front(); element != nil && !inserted; element = element.Next() {
			if element.Value.(*Cluster).Version <= cluster.Version {
				orderedList.InsertBefore(cluster, element)
				inserted = true
			}
		}
		if !inserted {
			orderedList.PushBack(cluster)
		}
	}
}

func (clusters *Clusters) Delete(uuidValue uuid.UUID, outputStream io.Writer) {
	clusterToDelete := clusters.ContextsByID[uuidValue.String()]
	if clusterToDelete != nil {
		if len(clusters.DockerHostEndpoint) > 0 {
			for _, container := range clusterToDelete.DockerConfigurations {
				dockerHost := clusters.DockerHostEndpoint
				if container.DockerHost != nil && len(container.DockerHost.Endpoint()) > 0 {
					dockerHost = container.DockerHost.Endpoint()
				}
				dockerClient, err := docker_client.NewDockerClient(dockerHost)
				if err == nil {
					if err = dockerClient.RemoveContainer(container.Name, 60, outputStream); err != nil {
						fmt.Fprintf(outputStream, "Error deleting docker container for name [%s]: %s\n", container.Name, err)
					}
				} else {
					fmt.Fprintf(outputStream, "Error creating docker client: %s\n", err)
				}
			}
		}
		deleteFromList(clusters.ContextsByVersion, uuidValue)
		delete(clusters.ContextsByID, uuidValue.String())
	}
}

func deleteFromList(orderedList *list.List, uuidValue uuid.UUID) {
	for element := orderedList.Front(); element != nil; element = element.Next() {
		if element.Value.(*Cluster).Uuid.String() == uuidValue.String() {
			orderedList.Remove(element)
			break;
		}
	}
}

func (clusters *Clusters) Get(uuidValue uuid.UUID) *Cluster {
	return clusters.ContextsByID[uuidValue.String()]
}

func (clusters *Clusters) GetByVersionOrder(age int) *Cluster {
	if clusters.ContextsByVersion != nil {
		var element *list.Element
		for element = clusters.ContextsByVersion.Front(); element != nil; element = element.Next() {
			if age == 0 {
				break
			}
			age--
		}
		if element != nil {
			return element.Value.(*Cluster)
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (clusters *Clusters) String() string {
	return clusters.ContextsByVersion.Front().Value.(*Cluster).String()
}

type BackendAddress struct {
	Address  *net.TCPAddr
	Host string
	Port string
}

func (this *BackendAddress) String() string {
	return fmt.Sprintf("%s:%s", this.Host, this.Port)
}

func (this *BackendAddress) Equals(other *BackendAddress) bool {
	return compare(this.Address.IP, other.Address.IP) && (this.Address.Port == other.Address.Port)
}

func compare(a, b []byte) bool {
	if &a == &b {
		return true
	}
	if len(a) != len(b) || len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

type Cluster struct {
	BackendAddresses                []*BackendAddress
	DockerConfigurations            []*docker_client.DockerConfig
	RequestCounter                  int64
	TransitionCounter               float64
	PercentageTransitionPerRequest  float64
	Uuid                            uuid.UUID
	SessionTimeout                  int64
	Mode                            TransitionMode
	Version                         string
}

func (cluster *Cluster) NextServer() (*tcp.TCPConnAndName, error) {
	cluster.RequestCounter++
	if len(cluster.BackendAddresses) > 0 {
		server := cluster.BackendAddresses[int(cluster.RequestCounter) % len(cluster.BackendAddresses)]
		message := fmt.Sprintf("Serving response %d from ip: [%s] port: [%d] version: [%s] mode: [%s]", cluster.RequestCounter, server.Address.IP, server.Address.Port, cluster.Version, ModesModeToCode[cluster.Mode])
		if cluster.PercentageTransitionPerRequest > 0 {
			message += fmt.Sprintf(" transition counter [%.2f] percentage transition per request [%.2f]", cluster.TransitionCounter, cluster.PercentageTransitionPerRequest)
		}
		if cluster.SessionTimeout > 0 {
			message += fmt.Sprintf(" session timeout [%d] uuid [%s]", cluster.SessionTimeout, cluster.Uuid)
		}
		log.LoggerFactory().Info(message)
		connection, err := net.DialTCP("tcp", nil, server.Address)
		return &tcp.TCPConnAndName{connection, server.Host, server.Port}, err
	} else {
		return nil, errors.New("Error no backend addresses found")
	}
}

func (cluster *Cluster) String() string {
	var result string = fmt.Sprintf("version: %s [", cluster.Version)
	for index, address := range cluster.BackendAddresses {
		if index > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%s:%s", address.Host, address.Port)
	}
	result += "]"
	return result
}
