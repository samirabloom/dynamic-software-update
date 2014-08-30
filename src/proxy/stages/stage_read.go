package stages

import (
	"io"
	"proxy/log"
	"proxy/contexts"
)

// ==== READ - START

func read(next func(*contexts.ChunkContext), complete func(*contexts.ChunkContext)) func(*contexts.ChunkContext) {
	return func(context *contexts.ChunkContext) {
		log.LoggerFactory().Debug("Read Stage START - %s", context)
		var loopCounter = 0
		for {
			log.LoggerFactory().Debug("Read Loop START - %d - %s", loopCounter, context)
			context.Data = context.Data[0:cap(context.Data)]
			readSize, readError := context.From.Read(context.Data)
			context.Data = context.Data[0:readSize]

			if readSize > 0 {
				context.TotalReadSize += int64(readSize)
				next(context)
				if context.FirstChunk {
					context.FirstChunk = false
				}
			}

			if context.Err != nil {
				log.LoggerFactory().Debug("Error reading connection %s - %s", context.Err, context)
				break
			}

			if readError == io.EOF {
				log.LoggerFactory().Debug("Read Loop EOF - %s", context)
				break
			}

			if readError != nil {
				log.LoggerFactory().Debug("Read Loop error %s - %s", readError, context)
				context.Err = readError
				break
			}

			log.LoggerFactory().Debug("Read Loop END - %d - %s", loopCounter, context)
			loopCounter++
		}
		complete(context)
		log.LoggerFactory().Debug("Read Stage END - %s", context)
	}
}

// ==== READ - END
