package stages

import (
	"net"
	"syscall"
	"time"
	"proxy/log"
	"proxy/contexts"
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
	if context.To != nil {
		_, assertion := context.To.(*net.TCPConn)
		if assertion {
			if context.To.(*net.TCPConn) != nil {
				closeReadError := context.To.CloseRead()
				log.LoggerFactory().Debug("Complete Stage closed READ with error %s - %s", closeReadError, context)
			}
		} else {
			closeReadError := context.To.CloseRead()
			log.LoggerFactory().Debug("Complete Stage closed READ with error %s - %s", closeReadError, context)
		}
	}
	context.PipeComplete <- context.TotalWriteSize
	log.LoggerFactory().Debug("Complete Stage END - %s", context)
}

// ==== COMPLETE - END
