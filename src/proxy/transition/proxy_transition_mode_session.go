package transition

import (
	"proxy/log"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"time"
	"proxy/contexts"
	"regexp"
)

var requestUUIDHeaderRegex = regexp.MustCompile("Cookie:.*dynsoftup=([a-z0-9-]*).*")

var SessionMode contexts.TransitionMode = RegisterTransitionMode(contexts.SessionMode, &SessionTransitionRouter{});

type SessionTransitionRouter struct {}

func (router *SessionTransitionRouter) 	route(clusters *contexts.Clusters, context *contexts.ChunkContext) (err error) {
	cluster := clusters.GetByVersionOrder(0)

	// find uuid cookie
	submatchs := requestUUIDHeaderRegex.FindSubmatch(context.Data)
	var requestUUID uuid.UUID
	if len(submatchs) >= 2 {
		requestUUID = uuid.Parse(string(submatchs[1]))
		log.LoggerFactory().Debug("Route Stage found request UUID %s", context)
	}

	// load cluster using uuid cookie
	if (requestUUID != nil && clusters.ContextsByID[requestUUID.String()] != nil) {
		cluster = clusters.ContextsByID[requestUUID.String()]
	}

	// add uuid cookie for cluster with expiry time
	context.RoutingContext = &contexts.RoutingContext{Headers: make([]string, 1)}
	context.RoutingContext.Headers[0] = fmt.Sprintf("Set-Cookie: dynsoftup=%s; Expires=%s;\n", cluster.Uuid.String(), time.Now().Add(time.Second*time.Duration(cluster.SessionTimeout)).Format(time.RFC1123))

	// create connection
	context.To, err = cluster.NextServer()

	return err
}
