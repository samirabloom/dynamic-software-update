package proxy_c

import (
	"bytes"
	uuid "code.google.com/p/go-uuid/uuid"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	byteutil "util/byte"
	"server"
	"container/list"
)

// ==== MAIN - START

func Proxy() {
	logLevel = flag.String("logLevel", "WARN", "Set the log level as \"CRITICAL\", \"ERROR\", \"WARNING\", \"NOTICE\", \"INFO\" or \"DEBUG\"")

	var cmd, _ = os.Getwd()
	if !strings.HasSuffix(cmd, "/") {
		cmd = cmd+"/"
	}
	var configFile = flag.String("configFile", cmd+"config.json", "Set the location of the configuration file")

	flag.Parse()

	go server.Server(1024, 1025, 1026, 1027, 1028, 1029, 1030, 1031)

	time.Sleep(1000 * time.Millisecond)

	loadBalancer, err := loadConfig(configFile)
	if err == nil {
		loadBalancer.Start()
	} else {
		loggerFactory().Error("Error parsing config %v", err)
	}
	var blocking = make(chan bool)
	<-blocking
}

// ==== MAIN - END

// ==== READ - START

func read(next func(*chunkContext), complete func(*chunkContext)) func(*chunkContext) {
	return func(context *chunkContext) {
		defer trace(time.Now(), context.performance.read)
		loggerFactory().Debug("Read Stage START - %s", context)
		var loopCounter = 0
		for {
			loggerFactory().Debug("Read Loop START - %d - %s", loopCounter, context)
			context.data = context.data[0:cap(context.data)]
			readSize, readError := context.from.Read(context.data)
			context.data = context.data[0:readSize]

			if readSize > 0 {
				context.totalReadSize += int64(readSize)
				next(context)
				loggerFactory().Debug("Error routing connection %s - %s", context.err, context)
				if context.firstChunk {
					context.firstChunk = false
				}
			}

			if context.err != nil {
				loggerFactory().Debug("Error routing connection %s - %s", context.err, context)
				break
			}

			if readError == io.EOF {
				loggerFactory().Debug("Read Loop EOF - %s", context)
				break
			}

			if readError != nil {
				loggerFactory().Debug("Read Loop error %s - %s", readError, context)
				context.err = readError
				break
			}

			loggerFactory().Debug("Read Loop END - %d - %s", loopCounter, context)
			loopCounter++
		}
		complete(context)
		loggerFactory().Debug("Read Stage END - %s", context)
	}
}

// ==== READ - END

// ==== ROUTE - START

var (
	cookieHeaderRegex = regexp.MustCompile("Cookie: .*dynsoftup=([a-z0-9-]*);.*")
)

func route(next func(*chunkContext), router *RoutingContexts, createBackPipe func(context *chunkContext)) func(*chunkContext) {
	return func(context *chunkContext) {
		defer trace(time.Now(), context.performance.route)
		loggerFactory().Debug("Route Stage START - %s", context)
		if context.firstChunk {
			if context.clientToServer {
				var err error
				submatchs := cookieHeaderRegex.FindSubmatch(context.data)
				var requestUUID uuid.UUID
				if len(submatchs) >= 2 {
					requestUUID = uuid.Parse(string(submatchs[1]))
					loggerFactory().Debug("Route Stage found UUID %s", context)
				}

				routingContext := router.NextServer(requestUUID)
				context.routingContext = routingContext

				backendAddr := routingContext.NextServer()
				loggerFactory().Info(fmt.Sprintf("Serving response %d from ip: [%s] port: [%d] version: [%.2f]", routingContext.requestCounter, backendAddr.IP, backendAddr.Port, routingContext.version))

				context.to, err = net.DialTCP("tcp", nil, backendAddr)
				if err != nil {
					loggerFactory().Error("Can't forward traffic to server tcp/%v: %s\n", backendAddr, err)
					if isConnectionRefused(err) {
						// no such device or address
						context.err = &net.OpError{Op: "dial", Addr: backendAddr, Err: syscall.ENXIO}
						context.pipeComplete <- 0
					}
					return
				}

				go createBackPipe(NewBackPipeChunkContext(context))
			} else {
				routingContext := context.routingContext
				setCookieHeader := []byte(fmt.Sprintf("Set-Cookie: dynsoftup=%s; Expires=%s;\n", routingContext.uuid.String(), time.Now().Add(time.Second * time.Duration(routingContext.sessionTimeout)).Format(time.RFC1123)))
				insertLocation := bytes.Index(context.data, []byte("\n"))
				if insertLocation > 0 {
					context.data = byteutil.Insert(context.data, insertLocation+len("\n"), setCookieHeader)
					context.totalReadSize += int64(len(setCookieHeader))
				}
			}
		}

		next(context)
		loggerFactory().Debug("Route Stage END - %s", context)
	}
}

// ==== ROUTE - END

// ==== WRITE - START

func write(context *chunkContext) {
	defer trace(time.Now(), context.performance.write)
	loggerFactory().Debug("Write Stage START - %s", context)
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
	loggerFactory().Debug("Write Stage END - %s", context)
}

// ==== WRITE - END

// ==== COMPLETE - START

func complete(context *chunkContext) {
	defer trace(time.Now(), context.performance.complete)
	loggerFactory().Debug("Complete Stage START - %s", context)
	if context.err != nil {
		// If the socket we are writing to is shutdown with
		// SHUT_WR, forward it to the other end of the createForwardPipe:
		if err, ok := context.err.(*net.OpError); ok && err.Err == syscall.EPIPE {
			closeWriteError := context.from.CloseWrite()
			loggerFactory().Debug("Complete Stage closed WRITE with error %s - %s", closeWriteError, context)
		}
	}
	if context.to != nil {
		_, assertion := context.to.(*net.TCPConn)
		if assertion {
			if context.to.(*net.TCPConn) != nil {
				closeReadError := context.to.CloseRead()
				loggerFactory().Debug("Complete Stage closed READ with error %s - %s", closeReadError, context)
			}
		} else {
			closeReadError := context.to.CloseRead()
			loggerFactory().Debug("Complete Stage closed READ with error %s - %s", closeReadError, context)
		}
	}
	context.pipeComplete <- context.totalWriteSize
	loggerFactory().Debug("Complete Stage END - %s", context)
}

// ==== COMPLETE - END

// ==== CREATE PIPE - START

func createPipe(routingContexts *RoutingContexts) func(*chunkContext) {
	return func(context *chunkContext) {
		loggerFactory().Debug("Creating " + context.description + " START")
		stages := read(
			route(
				write,
				routingContexts,
				createPipe(routingContexts),
			),
			complete,
		)
		stages(context)
		writePerformanceLogEntry(context)
		loggerFactory().Debug("Creating " + context.description + " END")
	}
}

// ==== CREATE PIPE - END

// ==== LOAD BALANCER - START

type RoutingContexts  struct {
	contextsByVersion *list.List
	contextsByID      map[string]*RoutingContext
}

func (routingContexts *RoutingContexts) NextServer(uuidValue uuid.UUID) *RoutingContext {
	routingContext := routingContexts.contextsByVersion.Front().Value.(*RoutingContext)
	if routingContext.mode != instantMode {
		if (uuidValue != nil && routingContexts.contextsByID[uuidValue.String()] != nil) {
			routingContext = routingContexts.contextsByID[uuidValue.String()]
		}
	}
	return routingContext
}

func (routingContexts *RoutingContexts) Add(routingContext *RoutingContext) {
	if routingContexts.contextsByVersion == nil {
		routingContexts.contextsByVersion = list.New()
	}
	if routingContexts.contextsByID == nil {
		routingContexts.contextsByID = make(map[string]*RoutingContext)
	}
	routingContextToAdd := routingContexts.contextsByID[routingContext.uuid.String()]
	if routingContextToAdd == nil {
		insertOrderedByVersion(routingContexts.contextsByVersion, routingContext)
		routingContexts.contextsByID[routingContext.uuid.String()] = routingContext
	}
}

func insertOrderedByVersion(orderedList *list.List, routingContext *RoutingContext) {
	if orderedList.Front() == nil {
		orderedList.PushFront(routingContext)
	} else {
		inserted := false
		for element := orderedList.Front(); element != nil && !inserted; element = element.Next() {
			if element.Value.(*RoutingContext).version <= routingContext.version {
				orderedList.InsertBefore(routingContext, element)
				inserted = true
			}
		}
		if !inserted {
			orderedList.PushBack(routingContext)
		}
	}
}

func (routingContexts *RoutingContexts) Delete(uuidValue uuid.UUID) {
	routingContextToDelete := routingContexts.contextsByID[uuidValue.String()]
	if routingContextToDelete != nil {
		deleteFromList(routingContexts.contextsByVersion, uuidValue)
		delete(routingContexts.contextsByID, uuidValue.String())
	}
}

func deleteFromList(orderedList *list.List, uuidValue uuid.UUID) {
	for element := orderedList.Front(); element != nil; element = element.Next() {
		if element.Value.(*RoutingContext).uuid.String() == uuidValue.String() {
			orderedList.Remove(element)
			break;
		}
	}
}

func (routingContexts *RoutingContexts) Get(uuidValue uuid.UUID) *RoutingContext {
	return routingContexts.contextsByID[uuidValue.String()]
}

func (routingContexts *RoutingContexts) String() string {
	return routingContexts.contextsByVersion.Front().Value.(*RoutingContext).String()
}

type RoutingContext struct {
	backendAddresses    []*net.TCPAddr
	requestCounter      int64
	uuid                uuid.UUID
	sessionTimeout		int64
	mode			    TransitionMode
	version             float64
}

func (routingContext *RoutingContext) NextServer() *net.TCPAddr {
	routingContext.requestCounter++
	return routingContext.backendAddresses[int(routingContext.requestCounter) % len(routingContext.backendAddresses)]
}

func (routingContext *RoutingContext) String() string {
	var result string = fmt.Sprintf("version: %.2f [", routingContext.version)
	for index, address := range routingContext.backendAddresses {
		if index > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%s", address)
	}
	result += "]"
	return result
}

type LoadBalancer struct {
	frontendAddr     *net.TCPAddr
	configServicePort int
	routingContexts  *RoutingContexts
	stop             chan bool
}

func (proxy *LoadBalancer) String() string {
	return fmt.Sprintf("LoadBalancer{\n\tProxy Address:      %s\n\tConfigService Port: %v\n\tProxied Servers:    %s\n}", proxy.frontendAddr, proxy.configServicePort, proxy.routingContexts)
}

func (proxy *LoadBalancer) Start() {
	var started = make(chan bool)
	go proxy.acceptLoop(started)
	go ConfigServer(proxy.configServicePort, proxy.routingContexts)
	<-started
}

func (proxy *LoadBalancer) Stop() {
	close(proxy.stop)
}

func (proxy *LoadBalancer) acceptLoop(started chan bool) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(proxy.frontendAddr.Port))
	if err != nil {
		log.Fatalf("Error opening socket - %s", err)
	}
	started <- true

	// allow for stopping
	running := true
	go func() {
		<-proxy.stop
		running = false
		listener.Close()
	}()

	for running {
		loggerFactory().Debug("Accept loop IN LOOP - START")

		// accept connection
		client, err := listener.Accept()

		if err != nil { // stop proxy

			if !isClosedError(err) {
				log.Printf("Stopping proxy on tcp/%v", listener.Addr().(*net.TCPAddr), err)
			}
			running = false

		} else { // process connection

			go func() {
				pipesComplete := make(chan int64)

				// create forward pipe
				forwardContext := NewForwardPipeChunkContext(client.(*net.TCPConn), pipesComplete)
				go createPipe(proxy.routingContexts)(forwardContext)

				// wait for pipes to complete (or quit early)
				for i := 0; i < 2; i++ {
					select {
					case <-pipesComplete:
					case <-proxy.stop:
						running = false
						listener.Close()
					}
				}

				// close sockets
				forwardContext.from.Close()
				if forwardContext.to != nil && forwardContext.to.(*net.TCPConn) != nil {
					forwardContext.to.Close()
				}
			}()

		}

		loggerFactory().Debug("Accept loop IN LOOP - END")
	}
}

func isClosedError(err error) bool {
	/* This comparison is ugly, but unfortunately, net.go doesn't export errClosing.
	 * See:
	 * http://golang.org/src/pkg/net/net.go
	 * https://code.google.com/p/go/issues/detail?id=4337
	 * https://groups.google.com/forum/#!msg/golang-nuts/0_aaCvBmOcM/SptmDyX1XJMJ
	 */
	return strings.HasSuffix(err.Error(), "use of closed network connection")
}

func isConnectionRefused(err error) bool {
	// This comparison is ugly, but unfortunately, net.go doesn't export appropriate error code.
	return strings.HasSuffix(err.Error(), "connection refused")
}

// ==== LOAD BALANCER - END

// ==== CHUNK_CONTEXT - START

type TCPConnection interface {
	Read(readBuffer []byte) (n int, err error)
	Write(writeBuffer []byte) (n int, err error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	ReadFrom(r io.Reader) (int64, error)
	CloseRead() error
	CloseWrite() error
	SetLinger(sec int) error
	SetKeepAlive(keepAlive bool) error
	SetKeepAlivePeriod(d time.Duration) error
	SetNoDelay(noDelay bool) error
}

type chunkContext struct {
	description            string
	data                   []byte
	to                     TCPConnection
	from                   TCPConnection
	err                    error
	totalReadSize          int64
	totalWriteSize         int64
	pipeComplete           chan int64
	firstChunk             bool
	performance            performance
	routingContext         *RoutingContext
	clientToServer         bool
}

func (context *chunkContext) String() string {
	var output string = ""
	output += "\n{\n"
	output += fmt.Sprintf("\t description: %s\n", context.description)
	if context.clientToServer {
		output += "\t direction: client->server\n"
	} else {
		output += "\t direction: server->client\n"
	}
	if len(context.data) > 0 {
		output += "\t data:\n\t\t"+strings.Replace(string(context.data), "\n", "\n\t\t", -1)
	}
	output += "\n"
	if context.from.(*net.TCPConn) != nil && context.from.LocalAddr() != nil && context.from.RemoteAddr() != nil {
		output += fmt.Sprintf("\t from: %s -> %s\n", context.from.LocalAddr(), context.from.RemoteAddr())
	}
	if context.to.(*net.TCPConn) != nil && context.to.LocalAddr() != nil && context.to.RemoteAddr() != nil {
		output += fmt.Sprintf("\t to: %s -> %s\n", context.to.LocalAddr(), context.to.RemoteAddr())
	}
	output += fmt.Sprintf("\t totalReadSize: %d\n", context.totalReadSize)
	output += fmt.Sprintf("\t totalWriteSize: %d\n", context.totalWriteSize)
	if context.routingContext != nil {
		output += fmt.Sprintf("\t routingContext UUID: %s\n", context.routingContext.uuid)
	}
	output += "}\n"
	return output
}

func NewForwardPipeChunkContext(from *net.TCPConn, pipeComplete chan int64) *chunkContext {
	return &chunkContext{
		description:    "forwardpipe",
		data:           make([]byte, 64*1024),
		from:           from,
		pipeComplete:   pipeComplete,
		firstChunk:     true,
		performance:    *&performance{
			read:       new(int64),
			route:      new(int64),
			write:      new(int64),
			complete:   new(int64),
		},
		routingContext: nil,
		clientToServer: true,
	}
}

func NewBackPipeChunkContext(forwardContext *chunkContext) *chunkContext {
	return &chunkContext{
		description:    "backpipe",
		data:           make([]byte, 64*1024),
		from:           forwardContext.to,
		to:             forwardContext.from,
		pipeComplete:   forwardContext.pipeComplete,
		firstChunk:     true,
		performance:    *&performance{
			read:       new(int64),
			route:      new(int64),
			write:      new(int64),
			complete:   new(int64),
		},
		routingContext: forwardContext.routingContext,
		clientToServer: false,
	}
}

// ==== CHUNK_CONTEXT - END

