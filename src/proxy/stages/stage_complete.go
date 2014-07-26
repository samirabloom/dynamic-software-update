package stages

import (
	"net"
	"syscall"
	"time"
	"proxy/log"
)

// ==== COMPLETE - START

func complete(context *ChunkContext) {
	defer log.Trace("complete", time.Now())
	log.LoggerFactory().Debug("Complete Stage START - %s", context)
	if context.err != nil {
		// If the socket we are writing to is shutdown with
		// SHUT_WR, forward it to the other end of the createForwardPipe:
		if err, ok := context.err.(*net.OpError); ok && err.Err == syscall.EPIPE {
			closeWriteError := context.from.CloseWrite()
			log.LoggerFactory().Debug("Complete Stage closed WRITE with error %s - %s", closeWriteError, context)
		}
	}
	if context.to != nil {
		_, assertion := context.to.(*net.TCPConn)
		if assertion {
			if context.to.(*net.TCPConn) != nil {
				closeReadError := context.to.CloseRead()
				log.LoggerFactory().Debug("Complete Stage closed READ with error %s - %s", closeReadError, context)
			}
		} else {
			closeReadError := context.to.CloseRead()
			log.LoggerFactory().Debug("Complete Stage closed READ with error %s - %s", closeReadError, context)
		}
	}
	context.pipeComplete <- context.totalWriteSize
	log.LoggerFactory().Debug("Complete Stage END - %s", context)
}

// ==== COMPLETE - END
