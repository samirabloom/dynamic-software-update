package contexts

import (
	"container/list"
	"code.google.com/p/go-uuid/uuid"
	"net"
	"fmt"
	"proxy/log"
)

type Clusters  struct {
	ContextsByVersion *list.List
	ContextsByID      map[string]*Cluster
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

func (clusters *Clusters) Delete(uuidValue uuid.UUID) {
	clusterToDelete := clusters.ContextsByID[uuidValue.String()]
	if clusterToDelete != nil {
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
	var element *list.Element
	for element = clusters.ContextsByVersion.Front(); element != nil; element = element.Next() {
		if age == 0 {
			break
		}
		age--
	}
	return element.Value.(*Cluster)
}

func (clusters *Clusters) String() string {
	return clusters.ContextsByVersion.Front().Value.(*Cluster).String()
}

type Cluster struct {
	BackendAddresses                []*net.TCPAddr
	RequestCounter                  int64
	TransitionCounter               float64
	PercentageTransitionPerRequest  float64
	Uuid                            uuid.UUID
	SessionTimeout                  int64
	Mode                            TransitionMode
	Version                         float64
}

func (cluster *Cluster) NextServer() *net.TCPAddr {
	cluster.RequestCounter++
	server := cluster.BackendAddresses[int(cluster.RequestCounter) % len(cluster.BackendAddresses)]
	message := fmt.Sprintf("Serving response %d from ip: [%s] port: [%d] version: [%.2f] mode: [%s]", cluster.RequestCounter, server.IP, server.Port, cluster.Version, ModesModeToCode[cluster.Mode])
	if cluster.PercentageTransitionPerRequest > 0 {
		message += fmt.Sprintf(" transition counter [%.2f] percentage transition per request [%.2f]", cluster.TransitionCounter, cluster.PercentageTransitionPerRequest)
	}
	if cluster.SessionTimeout > 0 {
		message += fmt.Sprintf(" session timeout [%d] uuid [%s]", cluster.SessionTimeout, cluster.Uuid)
	}
	log.LoggerFactory().Info(message)
	return server
}

func (cluster *Cluster) String() string {
	var result string = fmt.Sprintf("version: %.2f [", cluster.Version)
	for index, address := range cluster.BackendAddresses {
		if index > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%s", address)
	}
	result += "]"
	return result
}