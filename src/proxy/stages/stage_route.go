package stages

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"net"
	"regexp"
	"syscall"
	"time"
	byteutil "util/byte"
	"proxy/log"
	"strings"
)

// ==== ROUTE - START

var (
	cookieHeaderRegex = regexp.MustCompile("Cookie: .*dynsoftup=([a-z0-9-]*);.*")
)

func route(next func(*ChunkContext), router *Clusters, createBackPipe func(context *ChunkContext)) func(*ChunkContext) {
	return func(context *ChunkContext) {
		defer log.Trace("route", time.Now())
		log.LoggerFactory().Debug("Route Stage START - %s", context)
		if context.firstChunk {
			if context.clientToServer {
				var err error
				submatchs := cookieHeaderRegex.FindSubmatch(context.data)
				var requestUUID uuid.UUID
				if len(submatchs) >= 2 {
					requestUUID = uuid.Parse(string(submatchs[1]))
					log.LoggerFactory().Debug("Route Stage found UUID %s", context)
				}

				cluster := router.NextServer(requestUUID)
				context.cluster = cluster

				backendAddr := cluster.NextServer()
				log.LoggerFactory().Info(fmt.Sprintf("Serving response %d from ip: [%s] port: [%d] version: [%.2f]", cluster.RequestCounter, backendAddr.IP, backendAddr.Port, cluster.Version))

				context.to, err = net.DialTCP("tcp", nil, backendAddr)
				if err != nil {
					log.LoggerFactory().Error("Can't forward traffic to server tcp/%v: %s\n", backendAddr, err)
					if isConnectionRefused(err) {
						// no such device or address
						context.err = &net.OpError{Op: "dial", Addr: backendAddr, Err: syscall.ENXIO}
						context.pipeComplete <- 0
					}
					return
				}

				go createBackPipe(NewBackPipeChunkContext(context))
			} else {
				cluster := context.cluster
				setCookieHeader := []byte(fmt.Sprintf("Set-Cookie: dynsoftup=%s; Expires=%s;\n", cluster.Uuid.String(), time.Now().Add(time.Second*time.Duration(cluster.SessionTimeout)).Format(time.RFC1123)))
				insertLocation := bytes.Index(context.data, []byte("\n"))
				if insertLocation > 0 {
					context.data = byteutil.Insert(context.data, insertLocation+len("\n"), setCookieHeader)
					context.totalReadSize += int64(len(setCookieHeader))
				}
			}
		}

		next(context)
		log.LoggerFactory().Debug("Route Stage END - %s", context)
	}
}

func isConnectionRefused(err error) bool {
	// This comparison is ugly, but unfortunately, net.go doesn't export appropriate error code.
	return strings.HasSuffix(err.Error(), "connection refused")
}

// ==== ROUTE - END
