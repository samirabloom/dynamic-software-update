package stages

import (
	"io"
	"time"
	"proxy/log"
	"proxy/contexts"
)

// ==== WRITE - START

func write(context *contexts.ChunkContext) {
	defer log.Trace("write", time.Now())
	log.LoggerFactory().Debug("Write Stage START - %s", context)
	amountToWrite := len(context.Data)
	if amountToWrite > 0 {
		writeSize, writeError := context.To.Write(context.Data)
		if writeSize > 0 {
			context.TotalWriteSize += int64(writeSize)
		}
		if writeError != nil {
			context.Err = writeError
		} else if amountToWrite != writeSize {
			context.Err = io.ErrShortWrite
		}
	}
	log.LoggerFactory().Debug("Write Stage END - %s", context)
}

// ==== WRITE - END
