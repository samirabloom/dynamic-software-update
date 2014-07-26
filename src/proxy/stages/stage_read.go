package stages

import (
	"io"
	"time"
	"proxy/log"
)

// ==== READ - START

func read(next func(*ChunkContext), complete func(*ChunkContext)) func(*ChunkContext) {
	return func(context *ChunkContext) {
		defer log.Trace("read", time.Now())
		log.LoggerFactory().Debug("Read Stage START - %s", context)
		var loopCounter = 0
		for {
			log.LoggerFactory().Debug("Read Loop START - %d - %s", loopCounter, context)
			context.data = context.data[0:cap(context.data)]
			readSize, readError := context.from.Read(context.data)
			context.data = context.data[0:readSize]

			if readSize > 0 {
				context.totalReadSize += int64(readSize)
				next(context)
				if context.firstChunk {
					context.firstChunk = false
				}
			}

			if context.err != nil {
				log.LoggerFactory().Debug("Error routing connection %s - %s", context.err, context)
				break
			}

			if readError == io.EOF {
				log.LoggerFactory().Debug("Read Loop EOF - %s", context)
				break
			}

			if readError != nil {
				log.LoggerFactory().Debug("Read Loop error %s - %s", readError, context)
				context.err = readError
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
