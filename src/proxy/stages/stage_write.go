package stages

import (
	"io"
	"time"
	"proxy/log"
)

// ==== WRITE - START

func write(context *ChunkContext) {
	defer log.Trace("write", time.Now())
	log.LoggerFactory().Debug("Write Stage START - %s", context)
	amountToWrite := len(context.data)
	if amountToWrite > 0 {
		writeSize, writeError := context.to.Write(context.data)
		if writeSize > 0 {
			context.totalWriteSize += int64(writeSize)
		}
		if writeError != nil {
			context.err = writeError
		} else if amountToWrite != writeSize {
			context.err = io.ErrShortWrite
		}
	}
	log.LoggerFactory().Debug("Write Stage END - %s", context)
}

// ==== WRITE - END
