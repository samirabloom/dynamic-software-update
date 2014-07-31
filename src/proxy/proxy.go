package proxy

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"proxy/log"
	"code.google.com/p/go-uuid/uuid"
	"proxy/stages"
	"flag"
	"os"
)
var uuidGenerator = func(uuidValue uuid.UUID) func() uuid.UUID {
	return func() uuid.UUID {
		return uuidValue
	}
}(uuid.NewUUID())

// ==== LOAD BALANCER - START

type Proxy struct {
	frontendAddr      *net.TCPAddr
	configServicePort int
	clusters          *stages.Clusters
	stop              chan bool
}

func NewProxy(configFile string) *Proxy {
	proxy, err := loadConfig(configFile)
	if err != nil {
		log.LoggerFactory().Error("Error parsing config %v", err)
	}
	return proxy
}

func CLI() {
	log.LogLevel = flag.String("logLevel", "WARN", "Set the log level as \"CRITICAL\", \"ERROR\", \"WARNING\", \"NOTICE\", \"INFO\" or \"DEBUG\"")

	var cmd, _ = os.Getwd()
	if !strings.HasSuffix(cmd, "/") {
		cmd = cmd+"/"
	}
	var configFile = flag.String("configFile", cmd+"config.json", "Set the location of the configuration file")

	flag.Parse()

	NewProxy(*configFile).Start(true)
}

func (proxy *Proxy) String() string {
	return fmt.Sprintf("Proxy{\n\tProxy Address:      %s\n\tConfigService Port: %v\n\tProxied Servers:    %s\n}", proxy.frontendAddr, proxy.configServicePort, proxy.clusters)
}

func (proxy *Proxy) Start(blocking bool) {
	var started = make(chan bool)
	go proxy.acceptLoop(started)
	go ConfigServer(proxy.configServicePort, proxy.clusters)
	<-started

	if blocking {
		var block = make(chan bool)
		<-block
	}
}

func (proxy *Proxy) Stop() {
	close(proxy.stop)
}

func (proxy *Proxy) acceptLoop(started chan bool) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(proxy.frontendAddr.Port))
	if err != nil {
		panic(fmt.Sprintf("Error opening socket - %s", err))
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
		log.LoggerFactory().Debug("Accept loop IN LOOP - START")

		// accept connection
		client, err := listener.Accept()

		if err != nil { // stop proxy

			if !isClosedError(err) {
				log.LoggerFactory().Notice("Stopping proxy on tcp/%v", listener.Addr().(*net.TCPAddr), err)
			}
			running = false

		} else { // process connection

			go func() {
				pipesComplete := make(chan int64)

				// create forward pipe
				forwardContext := stages.NewForwardPipeChunkContext(client.(*net.TCPConn), pipesComplete)
				go stages.CreatePipe(proxy.clusters)(forwardContext)

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
				forwardContext.Close()
			}()

		}

		log.LoggerFactory().Debug("Accept loop IN LOOP - END")
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

// ==== LOAD BALANCER - END
