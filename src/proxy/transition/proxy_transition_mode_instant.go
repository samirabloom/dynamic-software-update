package transition

import (
	"net"
	"proxy/contexts"
)

var InstantMode contexts.TransitionMode = RegisterTransitionMode(contexts.InstantMode, &InstantTransitionRouter{});

type InstantTransitionRouter struct {}

func (router *InstantTransitionRouter) route(clusters *contexts.Clusters, context *contexts.ChunkContext) (err error) {
	cluster := clusters.GetByVersionOrder(0)

	context.To, err = net.DialTCP("tcp", nil, cluster.NextServer())

	return err
}
