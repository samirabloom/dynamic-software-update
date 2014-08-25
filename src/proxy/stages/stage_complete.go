package stages

import (
	"net"
	"syscall"
	"time"
	"proxy/log"
	"proxy/contexts"
	"proxy/tcp"
)

// ==== COMPLETE - START

func complete(context *contexts.ChunkContext) {
	defer log.Trace("complete", time.Now())
	log.LoggerFactory().Debug("Complete Stage START - %s", context)
	if context.Err != nil {
		// If the socket we are writing to is shutdown with
		// SHUT_WR, forward it to the other end of the createForwardPipe:
		if err, ok := context.Err.(*net.OpError); ok && err.Err == syscall.EPIPE {
			closeWriteError := context.From.CloseWrite()
			log.LoggerFactory().Debug("Complete Stage closed WRITE with error %s - %s", closeWriteError, context)
		}
	}
	contexts.AllowForNilConnection(context.To, func(connection tcp.TCPConnection) {
		closeReadError := connection.CloseRead()
		log.LoggerFactory().Debug("Complete Stage closed READ with error %s - %s", closeReadError, context)
	});
	context.PipeComplete <- context.TotalWriteSize
	log.LoggerFactory().Debug("Complete Stage END - %s", context)
}

// ==== COMPLETE - END
