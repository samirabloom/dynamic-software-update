package transition

import (
	"proxy/contexts"
	"code.google.com/p/go-uuid/uuid"
	"hash/fnv"
)


func RegisterTransitionMode(mode contexts.TransitionMode, router TransitionModeRouter) contexts.TransitionMode {
	contexts.ModesModeToRouteFunction[mode] = router.route
	return mode
}

type TransitionModeRouter interface {
	route(clusters *contexts.Clusters, context *contexts.ChunkContext) (err error)
}

func randomInteger() uint64 {
	hasher := fnv.New64()
	hasher.Write([]byte(uuid.NewUUID()))
	return hasher.Sum64()
}

