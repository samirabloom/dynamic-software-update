package stages

import "proxy/log"

// ==== CREATE PIPE - START

func CreatePipe(clusters *Clusters) func(*ChunkContext) {
	return func(context *ChunkContext) {
		log.LoggerFactory().Debug("Creating " + context.description + " START")
		stages := read(
			route(
				write,
				clusters,
				CreatePipe(clusters),
			),
			complete,
		)
		stages(context)
		log.EndPerformanceLogEntry()
		log.LoggerFactory().Debug("Creating " + context.description + " END")
	}
}

// ==== CREATE PIPE - END
