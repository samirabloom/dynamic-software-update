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
	"strconv"
	"hash/fnv"
	"proxy/tcp"
	"strings"
)

// ==== ROUTE - START

var (
	requestUUIDHeaderRegex    = regexp.MustCompile("Cookie: .*dynsoftup=([a-z0-9-]*);.*")
	transitionUUIDHeaderRegex = regexp.MustCompile("Cookie: .*transition=([a-z0-9-]*);.*")
	statusCodeRegex           = regexp.MustCompile("HTTP/[0-9].[0-9] ([a-z0-9-]*) .*")
	expiryRegex               = regexp.MustCompile("Expires: ([0-9-]*)")
	transferEncodingRegex     = regexp.MustCompile("Transfer-Encoding: ([a-z0-9-]*)")
	connectionRegex           = regexp.MustCompile("Connection: ([a-z0-9-]*)")
	contentTypeRegex          = regexp.MustCompile("Content-Type: ([a-z/a-z-]*);.*")

)

type headerMetrics struct {
	contentLength    int64
	statusCode       int
	expire           int64
	headers        map[string]string
}

func hashToPercentage(hash string) int64 {
	hasher := fnv.New64()
	hasher.Write([]byte(hash))
	return int64(hasher.Sum64() % 100)
}

func route(next func(*ChunkContext), clusters *Clusters, createBackPipe func(context *ChunkContext)) func(*ChunkContext) {
	return func(context *ChunkContext) {
		defer log.Trace("route", time.Now())
		log.LoggerFactory().Debug("Route Stage START - %s", context)
		if context.firstChunk {
			if context.clientToServer {  // on the request

				var err error
				cluster := clusters.GetByVersionOrder(0)

				switch {
				case cluster.Mode == SessionMode || cluster.Mode == GradualMode: {

					// find uuid cookie
					submatchs := requestUUIDHeaderRegex.FindSubmatch(context.data)
					var requestUUID uuid.UUID
					if len(submatchs) >= 2 {
						requestUUID = uuid.Parse(string(submatchs[1]))
						log.LoggerFactory().Debug("Route Stage found request UUID %s", context)
					}

					switch {
					case cluster.Mode == SessionMode: {
						// load cluster using uuid cookie
						if (requestUUID != nil && clusters.ContextsByID[requestUUID.String()] != nil) {
							cluster = clusters.ContextsByID[requestUUID.String()]
						}

						context.routingContext = &RoutingContext{headers: make([]string, 1)}
						context.routingContext.headers[0] = fmt.Sprintf("Set-Cookie: dynsoftup=%s; Expires=%s;\n", cluster.Uuid.String(), time.Now().Add(time.Second*time.Duration(cluster.SessionTimeout)).Format(time.RFC1123))
					}
					case cluster.Mode == GradualMode: {
						// find transition uuid cookie
						submatchs := transitionUUIDHeaderRegex.FindSubmatch(context.data)
						var transitionUUID uuid.UUID
						if len(submatchs) >= 2 {
							transitionUUID = uuid.Parse(string(submatchs[1]))
							log.LoggerFactory().Debug("Route Stage found transition UUID %s", context)
						}

						// load cluster using uuid cookie
						percentage := hashToPercentage(transitionUUID.String())

						if percentage >= cluster.RequestCounter {
							// do not latest cluster
							if (requestUUID != nil && clusters.ContextsByID[requestUUID.String()] != nil) {
								cluster = clusters.ContextsByID[requestUUID.String()]
							}
						}

						if transitionUUID != nil {
							transitionUUID = uuid.NewUUID()
						}

						context.routingContext = &RoutingContext{headers: make([]string, 2)}
						context.routingContext.headers[0] = fmt.Sprintf("Set-Cookie: dynsoftup=%s;", cluster.Uuid.String())
						context.routingContext.headers[1] = fmt.Sprintf("Set-Cookie: transition=%s;", transitionUUID.String())
					}
					}

					// create connection
					context.to, err = net.DialTCP("tcp", nil, cluster.NextServer())
				}
				case cluster.Mode == ConcurrentMode: {
					var (
						previousVersionConnection, latestVersionConnection tcp.TCPConnection
					)

					// create dual connection
					latestVersionConnection, err = net.DialTCP("tcp", nil, cluster.NextServer())
					if err == nil {
						previousVersionConnection, err = net.DialTCP("tcp", nil, clusters.GetByVersionOrder(1).NextServer())
						context.to = &tcp.DualTCPConnection{
							ExpectedStatusCode: 200,
							Connections:        []tcp.TCPConnection{previousVersionConnection, latestVersionConnection},
							SuccessfulIndex:    -1,
						}
					} else { // fall back to single connection if latest cluster fails on connection
						context.to, err = net.DialTCP("tcp", nil, clusters.GetByVersionOrder(1).NextServer())
					}
				}
				default: {
					// handle instant mode
					context.to, err = net.DialTCP("tcp", nil, cluster.NextServer())
				}
				}

				if err != nil {
					log.LoggerFactory().Error("Can't forward traffic to %v - %s\n", context.to, err)
					if isConnectionRefused(err) {
						// no such device or address
						context.err = &net.OpError{Op: "dial", Addr: context.to.RemoteAddr(), Err: syscall.ENXIO}
						context.pipeComplete <- 0
					}
					return
				}

				go createBackPipe(NewBackPipeChunkContext(context))

			} else { // on the response

				var parsedHeader = &headerMetrics{}
				parsedHeader.headers = make(map[string]string)
				parseHeader(parsedHeader, context.data)

				if context.routingContext != nil && len(context.routingContext.headers) > 0 { // if any headers to add
					insertLocation := bytes.Index(context.data, []byte("\n"))
					if insertLocation > 0 {
						for _, header := range context.routingContext.headers {
							context.data = byteutil.Insert(context.data, insertLocation+len("\n"), []byte(header))
							context.totalReadSize += int64(len(header))
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

func parseHeader(parsedHeader *headerMetrics, data []byte) {

	// checking for the contentLength in the http response
	parsedHeader.contentLength = int64(len(data))
	//	fmt.Printf("\nthe Content-Length found is: %d\n", parsedHeader.contentLength)
	log.LoggerFactory().Debug("Content-Length found is: %s", parsedHeader.contentLength)

	// checking for the status code in the http response
	statusCodeMatches := statusCodeRegex.FindSubmatch(data)
	if len(statusCodeMatches) >= 2 {
		statusCodeMatch := string(statusCodeMatches[1])
		parsedHeader.statusCode, _ = strconv.Atoi(statusCodeMatch)
		//		fmt.Printf("\nthe statusCode found is: %d\n", parsedHeader.statusCode)
		log.LoggerFactory().Debug("statusCode found is: %s", parsedHeader.statusCode)
	}

	// checking for "Expires" in the http response
	expiryMatches := expiryRegex.FindSubmatch(data)
	if len(expiryMatches) >= 2 {
		expiryMatch := string(expiryMatches[1])
		expiryHeader, _ := strconv.Atoi(expiryMatch)
		parsedHeader.expire = int64(expiryHeader)
		//		fmt.Printf("\nthe \"Expires\" found is: %d\n", parsedHeader.expire)
		log.LoggerFactory().Debug("\"Expires\" found is: %d", parsedHeader.expire)
	}

	// checking for Transfer-Encoding in the http response
	transferEncodingMatches := transferEncodingRegex.FindSubmatch(data)
	if len(transferEncodingMatches) >= 2 {
		parsedHeader.headers["Transfer-Encoding"] = string(transferEncodingMatches[1])
		//		fmt.Printf("\nthe Transfer-Encoding found is: %s\n", parsedHeader.headers["Transfer-Encoding"])
		log.LoggerFactory().Debug("Transfer-Encoding found is: %s", parsedHeader.headers["Transfer-Encoding"])
	}

	// checking for Connection in the http response
	connectionMatches := connectionRegex.FindSubmatch(data)
	if len(connectionMatches) >= 2 {
		parsedHeader.headers["Connection"] = string(connectionMatches[1])
		//		fmt.Printf("\nthe Connection found is: %s\n", parsedHeader.headers["Connection"])
		log.LoggerFactory().Debug("Connection found is: %s", parsedHeader.headers["Connection"])
	}

	// checking for Content-Type in the http response
	contentTypeMatches := contentTypeRegex.FindSubmatch(data)
	if len(contentTypeMatches) >= 2 {
		parsedHeader.headers["Content-Type"] = string(contentTypeMatches[0])
		//		fmt.Printf("\nthe Content-Type found is: %s\n", parsedHeader.headers["Content-Type"])
		log.LoggerFactory().Debug("Content-Type found is: %s\n", parsedHeader.headers["Content-Type"])
	}
}

// ==== ROUTE - END
