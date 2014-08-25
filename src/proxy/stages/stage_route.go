package stages

import (
	"bytes"
	"regexp"
	"time"
	byteutil "util/byte"
	"proxy/log"
	"strconv"
	"proxy/contexts"
	"syscall"
	"net"
	"strings"
)

// ==== ROUTE - START

var (
	statusCodeRegex = regexp.MustCompile("HTTP/[0-9].[0-9] ([a-z0-9-]*) .*")
)

type headerMetrics struct {
	contentLength    int64
	statusCode       int
	headers          map[string]string
}

func route(next func(*contexts.ChunkContext), clusters *contexts.Clusters, createBackPipe func(context *contexts.ChunkContext)) func(*contexts.ChunkContext) {
	return func(context *contexts.ChunkContext) {
		defer log.Trace("route", time.Now())
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
							clusters.Delete(cluster.Uuid)
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
				var parsedHeader = &headerMetrics{}
				parsedHeader.headers = make(map[string]string)
				parseMetrics(parsedHeader, context.Data)

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

func isConnectionRefused(err error) bool {
	// This comparison is ugly, but unfortunately, net.go doesn't export appropriate error code.
	return strings.HasSuffix(err.Error(), "connection refused")
}

func parseMetrics(parsedHeader *headerMetrics, response []byte) {

	// checking for the contentLength in the http response
	contentLengthHeader := parseHeader("Content-Length", response)
	if len(contentLengthHeader) > 0 {
		parsedHeader.contentLength, _ = strconv.ParseInt(contentLengthHeader, 10, 64)
	}

	// checking for the status code in the http response
	statusCodeMatches := statusCodeRegex.FindSubmatch(response)
	if len(statusCodeMatches) >= 2 {
		parsedHeader.statusCode, _ = strconv.Atoi(string(statusCodeMatches[1]))
	}

	parsedHeader.headers["Expires"] = parseHeader("Expires", response)
	parsedHeader.headers["Transfer-Encoding"] = parseHeader("Transfer-Encoding", response)
	parsedHeader.headers["Connection"] = parseHeader("Connection", response)
	parsedHeader.headers["Content-Type"] = parseHeader("Content-Type", response)
}

func parseHeader(headerName string, response []byte) string {
	contentTypeMatches := regexp.MustCompile(headerName + ": ([a-z/a-z-; =0-9]*)").FindSubmatch(response)
	if len(contentTypeMatches) >= 2 {
		return string(contentTypeMatches[1])
	}
	return ""
}

// ==== ROUTE - END
