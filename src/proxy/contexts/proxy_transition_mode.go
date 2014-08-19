package contexts

import (
	"proxy/log"
)

type TransitionMode uint64

const (
	InstantMode TransitionMode    = 1
	SessionMode TransitionMode    = 2
	GradualMode TransitionMode    = 3
	ConcurrentMode TransitionMode = 4
)

var ModesCodeToMode = map[string]TransitionMode {
	"INSTANT": InstantMode,
	"SESSION": SessionMode,
	"GRADUAL": GradualMode,
	"CONCURRENT": ConcurrentMode,
}

var ModesModeToCode = map[TransitionMode]string {
	InstantMode: "INSTANT",
	SessionMode: "SESSION",
	GradualMode: "GRADUAL",
	ConcurrentMode: "CONCURRENT",
}

var ModesModeToRouteFunction = make(map[TransitionMode]func(*Clusters, *ChunkContext) (error), 10)

func (mode *TransitionMode) Route(clusters *Clusters, context *ChunkContext) (err error) {
	cluster := clusters.GetByVersionOrder(0)

	if ModesModeToRouteFunction[cluster.Mode] == nil {
		var keys string = ""
		for key := range ModesModeToRouteFunction {
			if len(keys) > 0 {
				keys += ", "
			}
			keys += ModesModeToCode[key]
		}
		log.LoggerFactory().Error("Transition Mode %s not configured, only modes [%s] are available", ModesModeToCode[cluster.Mode], keys)
	} else {
		err = ModesModeToRouteFunction[cluster.Mode](clusters, context)
	}

	return err
}