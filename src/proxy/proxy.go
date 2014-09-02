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
	"proxy/transition"
	"proxy/contexts"
	"proxy/docker_client"
)

// this is a trick to ensure that go loads the proxy/transition package
var registerAllModes = transition.InstantMode

var uuidGenerator = func(uuidValue uuid.UUID) func() uuid.UUID {
	return func() uuid.UUID {
		return uuidValue
	}
}(uuid.NewUUID())

// ==== LOAD BALANCER - START

type Proxy struct {
	frontendAddr      *net.TCPAddr
	configServicePort int
	dockerHost        *docker_client.DockerHost
	clusters          *contexts.Clusters
	stop              chan bool
}

func NewProxy(configFile string) *Proxy {
	proxy, err := LoadConfig(configFile, os.Stdout)
	if err != nil {
		log.LoggerFactory().Error("%s", err.Error())
		os.Exit(1)
	}
	return proxy
}

func CLI() {
	log.LogLevel = flag.String("logLevel", "WARN", "Set the log level as \"CRITICAL\", \"ERROR\", \"WARNING\", \"NOTICE\", \"INFO\" or \"DEBUG\"\n")

	var configFile = flag.String("configFile", "./config.json", "Set the location of the configuration file that should contain configuration to start the proxy," +
				"\n                               for example:" +
				"\n                                           {" +
				"\n                                               \"proxy\": {" +
				"\n                                                   \"port\": 1235" +
				"\n                                               }," +
				"\n                                               \"configService\": {" +
				"\n                                                   \"port\": 9090" +
				"\n                                               }," +
				"\n                                               \"dockerHost\": {" +
				"\n                                                   \"ip\": \"127.0.0.1\"," +
				"\n                                                   \"port\": 2375" +
				"\n                                               }," +
				"\n                                               \"cluster\": {" +
				"\n                                                   \"containers\":[" +
				"\n                                                       {" +
				"\n                                                           \"image\": \"mysql\"," +
				"\n                                                           \"name\": \"some-mysql\"," +
				"\n                                                           \"environment\": [" +
				"\n                                                               \"MYSQL_ROOT_PASSWORD=mysecretpassword\"" +
				"\n                                                           ]," +
				"\n                                                           \"volumes\": [" +
				"\n                                                               \"/var/lib/mysql:/var/lib/mysql\"" +
				"\n                                                           ]" +
				"\n                                                       }," +
				"\n                                                       {" +
				"\n                                                           \"image\": \"wordpress\"," +
				"\n                                                           \"tag\": \"3.9.1\"," +
				"\n                                                           \"portToProxy\": 8080," +
				"\n                                                           \"name\": \"some-wordpress\"," +
				"\n                                                           \"links\": [" +
				"\n                                                               \"some-mysql:mysql\"" +
				"\n                                                           ]," +
				"\n                                                           \"portBindings\": {" +
				"\n                                                               \"80/tcp\": [" +
				"\n                                                                   {" +
				"\n                                                                       \"HostIp\": \"0.0.0.0\"," +
				"\n                                                                       \"HostPort\": \"8080\"" +
				"\n                                                                   }" +
				"\n                                                               ]" +
				"\n                                                           }" +
				"\n                                                       }" +
				"\n                                                   ]," +
				"\n                                                   \"version\": \"3.9.1\"" +
				"\n                                               }" +
				"\n                                           }\n")


	flag.Parse()

	NewProxy(*configFile).Start(true)
}

func (proxy *Proxy) String() string {
	proxyAddress := proxy.frontendAddr.String()
	if len(proxy.frontendAddr.IP) == 0 {
		proxyAddress = fmt.Sprintf("0.0.0.0:%d", proxy.frontendAddr.Port)
	}
	return fmt.Sprintf("Proxy{\n\tProxy Address:      %s\n\tConfigService Port: %v\n\tProxied Servers:    %s\n}", proxyAddress, proxy.configServicePort, proxy.clusters)
}

func (proxy *Proxy) Start(blocking bool) {
	var started = make(chan bool)
	go proxy.acceptLoop(started)
	go ConfigServer(proxy.configServicePort, proxy.clusters, proxy.dockerHost)
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
		log.LoggerFactory().Error("Error opening socket - %s", err)
		os.Exit(1)
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
				forwardContext := contexts.NewForwardPipeChunkContext(client.(*net.TCPConn), pipesComplete)
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
