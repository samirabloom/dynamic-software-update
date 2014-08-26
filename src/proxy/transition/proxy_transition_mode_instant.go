package transition

import (
	"proxy/contexts"
	"fmt"
)

var InstantMode contexts.TransitionMode = RegisterTransitionMode(contexts.InstantMode, &InstantTransitionRouter{});

type InstantTransitionRouter struct {}

func (router *InstantTransitionRouter) route(clusters *contexts.Clusters, context *contexts.ChunkContext) (err error) {
	cluster := clusters.GetByVersionOrder(0)

	context.To, err = cluster.NextServer()
	// add uuid cookie for cluster
	context.RoutingContext = &contexts.RoutingContext{Headers: make([]string, 1)}
	context.RoutingContext.Headers[0] = fmt.Sprintf("Set-Cookie: dynsoftup=%s;\n", cluster.Uuid.String())

	return err
}
