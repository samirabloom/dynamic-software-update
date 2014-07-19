package proxy_c

import (
	"bytes"
	uuid "code.google.com/p/go-uuid/uuid"
	"encoding/csv"
	"flag"
	"fmt"
	logging "github.com/op/go-logging"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	byteutil "util/byte"
	"io/ioutil"
	"encoding/json"
	"errors"
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

	go Server(1024, 1025, 1026, 1027, 1028, 1029, 1030, 1031)

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

// ==== PARSE CONFIG - START

func loadConfig(configFile *string) (*LoadBalancer, error) {
	return parseConfigFile(readConfigFile(configFile), parseProxy, parseClusterConfig(func() uuid.UUID {
			return uuid.NewUUID()
		}))
}

func readConfigFile(configFile *string) []byte {
	jsonConfig, err := ioutil.ReadFile(*configFile)
	if err != nil {
		loggerFactory().Error("Error %s reading config file [%s]", err, *configFile)
	}
	return jsonConfig
}

func parseConfigFile(jsonData []byte, parseProxy func(map[string]interface{}) (*net.TCPAddr, error), parseClusterConfig func(map[string]interface{}) (*RoutingContexts, error)) (loadBalancer *LoadBalancer, err error) {
	// parse json object
	var jsonConfig = make(map[string]interface{})
	err = json.Unmarshal(jsonData, &jsonConfig)
	if err != nil {
		loggerFactory().Error("Error %s parsing config file:\n%s", err.Error(), jsonData)
	}

	tcpProxyLocalAddress, proxyParseErr := parseProxy(jsonConfig)
	if proxyParseErr == nil {
		router, clusterParseErr := parseClusterConfig(jsonConfig)
		if clusterParseErr == nil {
			// create load balancer
			loadBalancer = &LoadBalancer{
				frontendAddr: tcpProxyLocalAddress,
				router: router,
				stop: make(chan bool),
			}
			loggerFactory().Info("Parsed config file:\n%s\nas:\n%s", jsonData, loadBalancer)

			return loadBalancer, nil
		} else {
			return nil, clusterParseErr
		}
	} else {
		return nil, proxyParseErr
	}
}

func parseProxy(jsonConfig map[string]interface{}) (tcpProxyLocalAddress *net.TCPAddr, err error) {
	if jsonConfig["proxy"] != nil {
		var proxyConfig map[string]interface{} = jsonConfig["proxy"].(map[string]interface{})
		tcpProxyLocalAddress, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%v", proxyConfig["ip"], proxyConfig["port"]))
		if err != nil {
			loggerFactory().Error("Invalid proxy address [" + fmt.Sprintf("%s:%v", proxyConfig["ip"], proxyConfig["port"]) + "]")
		}
	}
	if tcpProxyLocalAddress == nil {
		errorMessage := "Invalid proxy configuration - \"proxy\" JSON field missing or invalid"
		loggerFactory().Error(errorMessage)
		return nil, errors.New(errorMessage)
	}
	return tcpProxyLocalAddress, err
}

func parseClusterConfig(uuidGenerator func() uuid.UUID) func(map[string]interface{}) (*RoutingContexts, error) {
	return func(jsonConfig map[string]interface{}) (router *RoutingContexts, err error) {
		if jsonConfig["servers"] != nil {
			var servers = jsonConfig["servers"].([]interface {})
			var backendAddresses = make([]*net.TCPAddr, len(servers))
			for index := range servers {
				server := servers[index].(map[string]interface {})
				backendAddresses[index], err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%v", server["ip"], server["port"]))
				if err != nil {
					loggerFactory().Error("Invalid server address [" + fmt.Sprintf("%s:%v", server["ip"], server["port"]) + "]")
					return nil, err
				}
			}
			router = &RoutingContexts{
				all: make(map[string]*RoutingContext),
			}
			router.Add(&RoutingContext{backendAddresses: backendAddresses, requestCounter: -1, uuid: uuidGenerator()})
		}

		if router == nil {
			errorMessage := "Invalid cluster configuration - \"servers\" JSON field missing or invalid"
			loggerFactory().Error(errorMessage)
			return nil, errors.New(errorMessage)
		}
		return router, nil
	}
}

// ==== PARSE CONFIG - END

// ==== READ - START

func read(next func(*chunkContext), complete func(*chunkContext)) func(*chunkContext) {
	return func(context *chunkContext) {
		defer trace(time.Now(), context.performance.read)
		loggerFactory().Info("Read Stage START - %s", context)
		var loopCounter = 0
		for {
			loggerFactory().Info("Read Loop START - %d - %s", loopCounter, context)
			context.data = context.data[0:cap(context.data)]
			readSize, readError := context.from.Read(context.data)
			context.data = context.data[0:readSize]

			if readSize > 0 {
				context.totalReadSize += int64(readSize)
				next(context)
				loggerFactory().Info("Error routing connection %s - %s", context.err, context)
				if context.firstChunk {
					context.firstChunk = false
				}
			}

			if context.err != nil {
				loggerFactory().Info("Error routing connection %s - %s", context.err, context)
				break
			}

			if readError == io.EOF {
				loggerFactory().Info("Read Loop EOF - %s", context)
				break
			}

			if readError != nil {
				loggerFactory().Info("Read Loop error %s - %s", readError, context)
				context.err = readError
				break
			}

			loggerFactory().Info("Read Loop END - %d - %s", loopCounter, context)
			loopCounter++
		}
		complete(context)
		loggerFactory().Info("Read Stage END - %s", context)
	}
}

// ==== READ - END

// ==== ROUTE - START

var (
	cookieHeaderRegex = regexp.MustCompile("Cookie: .*dynsoftup=([a-z0-9-]*);.*")
)

func route(next func(*chunkContext), uuidGenerator func() uuid.UUID, router Router, createBackPipe func(context *chunkContext)) func(*chunkContext) {
	return func(context *chunkContext) {
		defer trace(time.Now(), context.performance.route)
		loggerFactory().Info("Route Stage START - %s", context)
		if context.firstChunk {
			if context.clientToServer {
				var err error

				backendAddr := router.NextServer(nil)
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

				submatchs := cookieHeaderRegex.FindSubmatch(context.data)
				if len(submatchs) >= 2 {
					context.requestUUID = uuid.Parse(string(submatchs[1]))
					loggerFactory().Info("Route Stage found UUID %s", context)
				}
				go createBackPipe(NewBackPipeChunkContext(context))
			} else {
				uuidCookieValue := context.requestUUID
				if uuidCookieValue == nil {
					uuidCookieValue = uuidGenerator()
				}
				setCookieHeader := []byte(fmt.Sprintf("Set-Cookie: dynsoftup=%s;\n", uuidCookieValue.String()))
				insertLocation := bytes.Index(context.data, []byte("\n"))
				if insertLocation > 0 {
					context.data = byteutil.Insert(context.data, insertLocation+len("\n"), setCookieHeader)
					context.totalReadSize += int64(len(setCookieHeader))
				}
			}
		}

		next(context)
		loggerFactory().Info("Route Stage END - %s", context)
	}
}

// ==== ROUTE - END

// ==== WRITE - START

func write(context *chunkContext) {
	defer trace(time.Now(), context.performance.write)
	loggerFactory().Info("Write Stage START - %s", context)
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
	loggerFactory().Info("Write Stage END - %s", context)
}

// ==== WRITE - END

// ==== COMPLETE - START

func complete(context *chunkContext) {
	defer trace(time.Now(), context.performance.complete)
	loggerFactory().Info("Complete Stage START - %s", context)
	if context.err != nil {
		// If the socket we are writing to is shutdown with
		// SHUT_WR, forward it to the other end of the createForwardPipe:
		if err, ok := context.err.(*net.OpError); ok && err.Err == syscall.EPIPE {
			closeWriteError := context.from.CloseWrite()
			loggerFactory().Info("Complete Stage closed WRITE with error %s - %s", closeWriteError, context)
		}
	}
	if context.to != nil {
		_, assertion := context.to.(*net.TCPConn)
		if assertion {
			if context.to.(*net.TCPConn) != nil {
				closeReadError := context.to.CloseRead()
				loggerFactory().Info("Complete Stage closed READ with error %s - %s", closeReadError, context)
			}
		} else {
			closeReadError := context.to.CloseRead()
			loggerFactory().Info("Complete Stage closed READ with error %s - %s", closeReadError, context)
		}
	}
	context.pipeComplete <- context.totalWriteSize
	loggerFactory().Info("Complete Stage END - %s", context)
}

// ==== COMPLETE - END

// ==== CREATE PIPE - START

func createPipe(router Router) func(*chunkContext) {
	return func(context *chunkContext) {
		loggerFactory().Info("Creating " + context.description + " START")
		stages := read(
			route(
				write,
				func() uuid.UUID {
					return uuid.NewUUID()
				},
				router,
				createPipe(router),
			),
			complete,
		)
		stages(context)
		writePerformanceLogEntry(context)
		loggerFactory().Info("Creating " + context.description + " END")
	}
}

// ==== CREATE PIPE - END

// ==== LOAD BALANCER - START

type Router interface {
	NextServer(uuid uuid.UUID) *net.TCPAddr
}

type RoutingContexts  struct {
	current *RoutingContext
	all     map[string]*RoutingContext
	mode    string
	timeout int
}

func (routingContexts *RoutingContexts) NextServer(uuid uuid.UUID) *net.TCPAddr {
	routingContext := routingContexts.current
	if (uuid != nil && routingContexts.all[uuid.String()] != nil) {
		routingContext = routingContexts.all[uuid.String()]
	}
	return routingContext.NextServer(nil)
}

func (routingContexts *RoutingContexts) Add(routingContext *RoutingContext) {
	routingContexts.current = routingContext
	routingContexts.all[routingContext.uuid.String()] = routingContext
}

func (routingContexts *RoutingContexts) Delete(uuid uuid.UUID) {
	delete(routingContexts.all, uuid.String())
}

func (routingContexts *RoutingContexts) Get(uuid uuid.UUID) *RoutingContext {
	return routingContexts.all[uuid.String()]
}

func (routingContexts *RoutingContexts) String() string {
	return routingContexts.current.String()
}

type RoutingContext struct {
	backendAddresses []*net.TCPAddr
	requestCounter   int64
	uuid             uuid.UUID
}

func (routingContext *RoutingContext) NextServer(uuid uuid.UUID) *net.TCPAddr {
	routingContext.requestCounter++
	return routingContext.backendAddresses[int(routingContext.requestCounter) % len(routingContext.backendAddresses)]
}

func (routingContext *RoutingContext) String() string {
	var result string = ""
	for index, address := range routingContext.backendAddresses {
		if index > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%s", address)
	}
	return result
}

type LoadBalancer struct {
	frontendAddr   *net.TCPAddr
	router         Router
	stop           chan bool
}

func (proxy *LoadBalancer) String() string {
	return fmt.Sprintf("LoadBalancer{\n\tProxy Address:   %s\n\tProxied Servers: %s\n}", proxy.frontendAddr, proxy.router)
}

func (proxy *LoadBalancer) Start() {
	var started = make(chan bool)
	go proxy.acceptLoop(started)
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
		loggerFactory().Info("Accept loop IN LOOP - START")

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
				go createPipe(proxy.router)(forwardContext)

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

		loggerFactory().Info("Accept loop IN LOOP - END")
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
	requestUUID            uuid.UUID
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
	if context.requestUUID != nil {
		output += fmt.Sprintf("\t requestUUID: %s\n", context.requestUUID)
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
		requestUUID:    nil,
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
		requestUUID:    forwardContext.requestUUID,
		clientToServer: false,
	}
}

// ==== CHUNK_CONTEXT - END

// ==== PERFORMANCE - START

type performance struct {
	read     *int64
	route    *int64
	write    *int64
	complete *int64
}

func trace(startTime time.Time, result *int64) {
	*result = int64(time.Since(startTime))
}

var performanceLog = func() *csv.Writer {
	file, error := os.Create("performance_log.csv")

	if error != nil {
		panic(error)
	}

	// New Csv writer
	writer := csv.NewWriter(file)

	// Headers
	var new_headers = []string{"count", "read", "route", "write", "complete"}
	returnError := writer.Write(new_headers)
	if returnError != nil {
		loggerFactory().Error("Error writing headers into performance log - %s", returnError)
	}
	writer.Flush()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGKILL)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGHUP)
	signal.Notify(c, syscall.SIGQUIT)
	go func() {
		<-c
		file.Close()
		os.Exit(1)
	}()

	return writer
}()

func writePerformanceLogEntry(context *chunkContext) {
	performanceLog.Write([]string{
	strconv.FormatInt(*context.performance.read, 10),
	strconv.FormatInt(*context.performance.route, 10),
	strconv.FormatInt(*context.performance.write, 10),
	strconv.FormatInt(*context.performance.complete, 10)})
	performanceLog.Flush()
}

// ==== PERFORMANCE - END

// ==== LOGGER - START

var logLevel *string

var loggerFactory = func() func() *logging.Logger {
	var logg *logging.Logger = nil

	return func() *logging.Logger {
		if logg == nil {
			logg = logging.MustGetLogger("main")

			// Customize the output format
			logging.SetFormatter(logging.MustStringFormatter("%{level:8s} - %{message}"))

			// Setup one stdout and one syslog backend
			logBackend := logging.NewLogBackend(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
			logBackend.Color = true

			// Combine them both into one logging backend
			logging.SetBackend(logBackend)

			// set log level
			level, _ := logging.LogLevel("WARN")
			if logLevel != nil {
				level, _ = logging.LogLevel(*logLevel)
			}
			logging.SetLevel(level, "main")
		}
		return logg
	}
}()

// ==== LOGGER - END

