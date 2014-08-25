package transition

import (
	"proxy/tcp"
	"proxy/contexts"
)

var ConcurrentMode contexts.TransitionMode = RegisterTransitionMode(contexts.ConcurrentMode, &ConcurrentTransitionRouter{});

type ConcurrentTransitionRouter struct {}

func (router *ConcurrentTransitionRouter) route(clusters *contexts.Clusters, context *contexts.ChunkContext) (err error) {
	cluster := clusters.GetByVersionOrder(0)

	var (
		previousVersionConnection, latestVersionConnection *contexts.TCPConnAndName
	)

	// create dual connection
	latestVersionConnection, err = cluster.NextServer()
	if err == nil {
		previousVersionConnection, err = clusters.GetByVersionOrder(1).NextServer()
		context.To = &tcp.DualTCPConnection{
			Connections:        []tcp.TCPConnection{previousVersionConnection, latestVersionConnection},
			Hosts:              []string{previousVersionConnection.Host, latestVersionConnection.Host},
			Ports:              []string{previousVersionConnection.Port, latestVersionConnection.Port},
			SuccessfulIndex:    -1,
		}
	} else { // fall back to single connection if latest cluster fails on connection
		context.To, err = clusters.GetByVersionOrder(1).NextServer()
	}

	return err
}
