package transition

import (
	"net"
	"proxy/tcp"
	"proxy/contexts"
)

var ConcurrentMode contexts.TransitionMode = RegisterTransitionMode(contexts.ConcurrentMode, &ConcurrentTransitionRouter{});

type ConcurrentTransitionRouter struct {}

func (router *ConcurrentTransitionRouter) route(clusters *contexts.Clusters, context *contexts.ChunkContext) (err error) {
	cluster := clusters.GetByVersionOrder(0)

	var (
		previousVersionConnection, latestVersionConnection tcp.TCPConnection
	)

	// create dual connection
	latestVersionConnection, err = net.DialTCP("tcp", nil, cluster.NextServer())
	if err == nil {
		previousVersionConnection, err = net.DialTCP("tcp", nil, clusters.GetByVersionOrder(1).NextServer())
		context.To = &tcp.DualTCPConnection{
			ExpectedStatusCode: 200,
			Connections:        []tcp.TCPConnection{previousVersionConnection, latestVersionConnection},
			SuccessfulIndex:    -1,
		}
	} else { // fall back to single connection if latest cluster fails on connection
		context.To, err = net.DialTCP("tcp", nil, clusters.GetByVersionOrder(1).NextServer())
	}

	return err
}
