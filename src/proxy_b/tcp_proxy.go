package proxy_b

import (
	"io"
	"log"
	"net"
	"syscall"
	"flag"
	"time"
	"regexp"
	uuid "code.google.com/p/go-uuid/uuid"
	logging "github.com/op/go-logging"
	"os"
	"fmt"
	"bytes"
	"strings"
	byteutil "util/byte"
	"strconv"
	"net/http"
)

// ==== LOGGER - START

var logger *logging.Logger = nil

var loggerFactory = func() func(*string) *logging.Logger {
	var logg *logging.Logger = nil

	return func(logLevel *string) *logging.Logger {
		if logg == nil {
			logg = logging.MustGetLogger("main")

			// Customize the output format
			logging.SetFormatter(logging.MustStringFormatter("%{level:8s} - %{message}"))

			// Setup one stdout and one syslog backend
			logBackend := logging.NewLogBackend(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
			logBackend.Color = true

			// Combine them both into one logging backend
			logging.SetBackend(logBackend)

			// set log level
			level, _ := logging.LogLevel(*logLevel)
			logging.SetLevel(level, "main")
		}
		return logg
	}
}()

// ==== LOGGER - END

// ==== SERVER - START

func Server(ports ...int) {
	for _, port := range ports {
		logger.Info("Starting server " + strconv.Itoa(port) + " ....")
		go http.ListenAndServe(":"+strconv.Itoa(port), &handle{port})
	}
}

type handle struct {
	Port int
}

func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "ONE, ")
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "TWO, ")
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "THREE, ")
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "FOUR, ")
	time.Sleep(10 * time.Millisecond)
	fmt.Fprintf(w, "FIVE from %d\n", h.Port)
}

// ==== SERVER - END

// ==== MAIN - START

func Proxy() {
	logLevel := flag.String("logLevel", "WARN", "Set the log level as \"CRITICAL\", \"ERROR\", \"WARNING\", \"NOTICE\", \"INFO\" or \"DEBUG\"")
	flag.Parse()

	logger = loggerFactory(logLevel)

	go Server(1024, 1025, 1026)

	time.Sleep(1000 * time.Millisecond)

	proxy, err := NewLoadBalancer(
		&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234},
		&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1024},
		3,
	)
	if proxy != nil && proxy.listener != nil {
		proxy.Run()
	} else {
		log.Fatalf("Error opening socket - %s", err)
	}
}

// ==== MAIN - END


type LoadBalancer struct {
	listener     *net.TCPListener
	frontendAddr *net.TCPAddr
	backendAddr  *net.TCPAddr
	loadBalanceCount int
	stop             chan bool
	uuidGenerator func() string
}

func NewLoadBalancer(frontendAddr, backendAddr *net.TCPAddr, loadBalanceCount int) (*LoadBalancer, error) {
	listener, err := net.ListenTCP("tcp", frontendAddr)
	if err != nil {
		return nil, err
	}
	// If the port in frontendAddr was 0 then ListenTCP will have a picked
	// a port to listen on, hence the call to Addr to get that actual port:
	return &LoadBalancer{
		listener:     listener,
		frontendAddr: listener.Addr().(*net.TCPAddr),
		backendAddr:  backendAddr,
		loadBalanceCount: loadBalanceCount,
		stop: make(chan bool),
		uuidGenerator: func() string {
			return uuid.NewUUID().String()
		},
	}, nil
}

var (
	cookieHeaderRegex = regexp.MustCompile("Cookie: .*dynsoftup=([a-z0-9-]*);.*")
)

func (proxy *LoadBalancer) read(dst io.Writer, src io.Reader, isRequest bool) (read int64, written int64, err error) {
	data := make([]byte, 32*1024)
	for {
		readSize, readError := src.Read(data)

		//		fmt.Printf("\nbefore insert: \n%s\n", data)

		if read <= 0 {
			if isRequest {
				submatches := cookieHeaderRegex.FindSubmatch(data[0:readSize])
				if len(submatches) >= 2 {
					fmt.Printf("dynsoftup value is: %s\n", string(submatches[1]))
				}
			} else if readSize > 0 {
//				fmt.Printf("Before insert: \n%s\n", data)
				setCookieHeader := []byte(fmt.Sprintf("Set-Cookie: dynsoftup=%s;\n", proxy.uuidGenerator()))
				searchString := "\n"
				insertLocation := bytes.Index(data[0:readSize], []byte(searchString))
				if insertLocation > 0 {
					byteutil.Insert(data[0:readSize], insertLocation+len(searchString), setCookieHeader)
					readSize += len(setCookieHeader)
				}
//				fmt.Printf("\nAfter insert: \n%s\n", data)
			}
		}
		if readSize > 0 {
			read += int64(readSize)
		}
		writeSize, writeError := proxy.write(dst, data[0:readSize])
		if writeSize > 0 {
			written += int64(writeSize)
		}
		if writeError != nil {
			err = readError
			break
		}
		if readError == io.EOF {
			break
		}
		if readError != nil {
			err = readError
			break
		}
	}
	return read, written, err
}

func (proxy *LoadBalancer) write(dst io.Writer, data []byte) (written int64, err error) {
	nr := len(data)
	if nr > 0 {
		nw, ew := dst.Write(data)
		if nw > 0 {
			written += int64(nw)
		}
		if ew != nil {
			err = ew
		} else if nr != nw {
			err = io.ErrShortWrite
		}
	}
	return written, err
}

func (proxy *LoadBalancer) clientLoop() func(*net.TCPConn, chan bool) {
	var requestNumber = 0
	return func(client *net.TCPConn, quit chan bool) {
		requestNumber++
		backend, err := proxy.routeRequest(client, requestNumber); if err != nil {
			return
		}

		event := make(chan int64)
		var broker = func(from, to *net.TCPConn, isRequest bool) {
			_, written, err := proxy.read(to, from, isRequest)
			if err != nil {
				// If the socket we are writing to is shutdown with
				// SHUT_WR, forward it to the other end of the pipe:
				if err, ok := err.(*net.OpError); ok && err.Err == syscall.EPIPE {
					from.CloseWrite()
				}
			}
			to.CloseRead()
			event <- written
		}

		go broker(client, backend, true)
		go broker(backend, client, false)

		var transferred int64 = 0
		for i := 0; i < 2; i++ {
			select {
			case written := <-event:
				transferred += written
			case <-quit:
				// Interrupt the two brokers and "join" them.
				client.Close()
				backend.Close()
				for ; i < 2; i++ {
					transferred += <-event
				}
				return
			}
		}
		client.Close()
		backend.Close()
	}
}

func (proxy *LoadBalancer) routeRequest(client *net.TCPConn, requestNumber int) (*net.TCPConn, error) {
	backendAddr := &net.TCPAddr{IP: proxy.backendAddr.IP, Port: proxy.backendAddr.Port + (requestNumber % proxy.loadBalanceCount)}
	backend, err := net.DialTCP("tcp", nil, backendAddr)
	if err != nil {
		log.Printf("Can't forward traffic to backend tcp/%v: %s\n", backendAddr, err)
		client.Close()
	}
	return backend, err
}

func (proxy *LoadBalancer) Run() {
	quitRequest := make(chan bool)
	defer close(quitRequest)

	running := true
	go func() {
		<-proxy.stop
		running = false
		proxy.listener.Close()
	}()

	clientLoop := proxy.clientLoop()
	for running {
		client, err := proxy.listener.Accept()
		if err != nil {
			if !isClosedError(err) {
				log.Printf("Stopping proxy on tcp/%v for tcp/%v (%v)", proxy.frontendAddr, proxy.backendAddr, err)
			}
			return
		}
		go clientLoop(client.(*net.TCPConn), quitRequest)
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

func (proxy *LoadBalancer) Stop() {
	close(proxy.stop)
}
