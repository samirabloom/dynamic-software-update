package stages

import (
	"bytes"
	byteutil "util/byte"
	"proxy/log"
	"proxy/contexts"
	"syscall"
	"net"
	"os"
)

// ==== ROUTE - START

func route(next func(*contexts.ChunkContext), clusters *contexts.Clusters, createBackPipe func(context *contexts.ChunkContext)) func(*contexts.ChunkContext) {
	return func(context *contexts.ChunkContext) {
		log.LoggerFactory().Debug("Route Stage START - %s", context)
		if context.FirstChunk {

			if context.Direction == contexts.ClientToServer {  // on the request

				for {
					cluster := clusters.GetByVersionOrder(0)
					if cluster != nil {
						err := cluster.Mode.Route(clusters, context)
						if err != nil {
							log.LoggerFactory().Error("Error communicating with server - %s\n", err)
							log.LoggerFactory().Warning("Removing cluster from configuration - %s\n", cluster)
							clusters.Delete(cluster.Uuid, os.Stdout)
							continue;
						} else {
							go createBackPipe(contexts.NewBackPipeChunkContext(context))
							break;
						}
					} else {
						log.LoggerFactory().Error("No clusters in configuration\n")
						context.Err = &net.OpError{Op: "dial", Err: syscall.ENXIO}
						context.PipeComplete <- 0
						return
					}
				}

			} else { // on the response
				if context.RoutingContext != nil && len(context.RoutingContext.Headers) > 0 { // if any headers to add
					insertLocation := bytes.Index(context.Data, []byte("\n"))
					if insertLocation > 0 {
						for _, header := range context.RoutingContext.Headers {
							context.Data = byteutil.Insert(context.Data, insertLocation+len("\n"), []byte(header))
							context.TotalReadSize += int64(len(header))
						}
					}
				}
			}
		}
		next(context)
		log.LoggerFactory().Debug("Route Stage END - %s", context)
	}
}

// ==== ROUTE - END
