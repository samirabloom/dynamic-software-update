package transition

import (
	"proxy/log"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"proxy/contexts"
	"hash/fnv"
	"regexp"
)

var transitionUUIDHeaderRegex = regexp.MustCompile("Cookie:.*transition=([a-z0-9-]*).*")

var GradualMode contexts.TransitionMode = RegisterTransitionMode(contexts.GradualMode, &GradualTransitionRouter{});

type GradualTransitionRouter struct {}

func (router *GradualTransitionRouter) route(clusters *contexts.Clusters, context *contexts.ChunkContext) (err error) {
	cluster := clusters.GetByVersionOrder(0)

	// find uuid cookie
	requestUUIDSubMatches := requestUUIDHeaderRegex.FindSubmatch(context.Data)
	var requestUUID uuid.UUID
	if len(requestUUIDSubMatches) >= 2 {
		requestUUID = uuid.Parse(string(requestUUIDSubMatches[1]))
		log.LoggerFactory().Debug("Route Stage found request UUID %s", context)
	}

	// find transition uuid cookie
	transitionUUIDSubMatches := transitionUUIDHeaderRegex.FindSubmatch(context.Data)
	var transitionUUID uuid.UUID
	if len(transitionUUIDSubMatches) >= 2 {
		transitionUUID = uuid.Parse(string(transitionUUIDSubMatches[1]))
		log.LoggerFactory().Debug("Route Stage found transition UUID %s", context)
	}

	if transitionUUID == nil {
		transitionUUID = uuid.NewUUID()
	}

	// determine transition percentage for request
	percentage := hashToPercentage(transitionUUID.String())

	cluster.TransitionCounter += cluster.PercentageTransitionPerRequest
	if float64(percentage) >= cluster.TransitionCounter {
		// do not latest cluster
		if (requestUUID != nil && clusters.ContextsByID[requestUUID.String()] != nil) {
			cluster = clusters.ContextsByID[requestUUID.String()]
		} else {
			cluster = clusters.GetByVersionOrder(1)
		}
	}

	// add uuid and transition cookies
	context.RoutingContext = &contexts.RoutingContext{Headers: make([]string, 2)}
	context.RoutingContext.Headers[0] = fmt.Sprintf("Set-Cookie: dynsoftup=%s;\n", cluster.Uuid.String())
	context.RoutingContext.Headers[1] = fmt.Sprintf("Set-Cookie: transition=%s;\n", transitionUUID.String())

	// create connection
	context.To, err = cluster.NextServer()

	return err
}

func hashToPercentage(hash string) int64 {
	hasher := fnv.New64()
	hasher.Write([]byte(hash))
	return int64(hasher.Sum64() % 100)
}
