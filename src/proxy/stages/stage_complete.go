package stages

import (
	"net"
	"syscall"
	"proxy/log"
	"proxy/contexts"
	"proxy/tcp"
	"fmt"
)

// ==== COMPLETE - START

func complete(context *contexts.ChunkContext) {
	log.LoggerFactory().Debug("Complete Stage START - %s", context)
	if context.Err != nil {
		// If the socket we are writing to is shutdown with
		// SHUT_WR, forward it to the other end of the createForwardPipe:
		if err, ok := context.Err.(*net.OpError); ok && err.Err == syscall.EPIPE {
			closeWriteError := context.From.CloseWrite()
			log.LoggerFactory().Debug("Complete Stage closed WRITE with error %s - %s", closeWriteError, context)
		}
	}
	tcp.AllowForNilConnection(context.To, func(connection tcp.TCPConnection) {
		closeReadError := connection.CloseRead()
		log.LoggerFactory().Debug("Complete Stage closed READ with error %s - %s", closeReadError, context)
	});
	context.PipeComplete <- context.TotalWriteSize

	if context.Direction == contexts.ClientToServer {
		var connectionClosedMessage = ""
		tcp.AllowForNilConnection(context.From, func(connection tcp.TCPConnection) {
			connectionClosedMessage += fmt.Sprintf("%s", connection.LocalAddr())
		});
		tcp.AllowForNilConnection(context.To, func(connection tcp.TCPConnection) {
			if len(connectionClosedMessage) > 0 {
				connectionClosedMessage += fmt.Sprintf(" -> %s", connection.RemoteAddr())
				log.LoggerFactory().Info(connectionClosedMessage + " -- Connection Closed")
			}
		});
	}
	log.LoggerFactory().Debug("Complete Stage END - %s", context)
}

// ==== COMPLETE - END
