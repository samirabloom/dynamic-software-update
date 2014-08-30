package stages

import (
	"proxy/log"
	"proxy/contexts"
)

// ==== CREATE PIPE - START

func CreatePipe(clusters *contexts.Clusters) func(*contexts.ChunkContext) {
	return func(context *contexts.ChunkContext) {
		log.LoggerFactory().Debug("Creating " + contexts.DirectionToDescription[context.Direction] + " START")
		stages := read(
			route(
				write,
				clusters,
				CreatePipe(clusters),
			),
			complete,
		)
		stages(context)
		log.LoggerFactory().Debug("Creating " + contexts.DirectionToDescription[context.Direction] + " END")
	}
}

// ==== CREATE PIPE - END
