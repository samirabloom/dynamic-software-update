package proxy_a

import (
	"time"
	"fmt"
	"net"
	"strconv"
	"regexp"
	"log"
	"os"
	logging "github.com/op/go-logging"
	"flag"
	"strings"
	"io"
	"net/http"
	"runtime"
)

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

func Proxy() {
	logLevel := flag.String("logLevel", "WARN", "Set the log level as \"CRITICAL\", \"ERROR\", \"WARNING\", \"NOTICE\", \"INFO\" or \"DEBUG\"")
	flag.Parse()

	logger = loggerFactory(logLevel)

	go Server(1024, 1025)

	time.Sleep(1000 * time.Millisecond)

	AcceptLoop(1234)
}

// ==== MAIN - END

// ==== CHUNK_CONTEXT - START

type chunkContext struct {
	data             []byte
	complete *bool // todo - to remove
	clientConnection *net.Conn
	serverConnection *net.Conn
	err              error
	totalReadSize    int
	totalWriteSize   int
}

type ChunkContext interface {
	String()
}

func (context *chunkContext) String() string {
	var output string = ""
	output += "\n\n{"
	if len(context.data) > 0 {
		output += "\n\tdata:\n\t\t"+strings.Replace(string(context.data), "\n", "\n\t\t", -1)
	}
	output += fmt.Sprintf("\n\tcomplete: %t", *context.complete)
	if *context.clientConnection != nil {
		output += "\n\tclientConntection: "+(*context.clientConnection).LocalAddr().String()+" -> "+(*context.clientConnection).RemoteAddr().String()
	}
	if *context.serverConnection != nil {
		output += "\n\tserverConnection: "+(*context.serverConnection).LocalAddr().String()+" <- "+(*context.serverConnection).RemoteAddr().String()
	}
	output += "\n}\n"
	return output
}

func NewChunkContext() *chunkContext {
	return &chunkContext{make([]byte, 4096), new(bool), new(net.Conn), new(net.Conn), nil, 0, 0}
}

// ==== CHUNK_CONTEXT - END

// ==== ACCEPT LOOP - START

func AcceptLoop(port int) {
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		logger.Error("Error Listening on \"localhost:%d\" -- %v", port, err)
	} else {
		for {
			chunkContext := NewChunkContext();
			if *chunkContext.clientConnection, err = listener.Accept(); err != nil {
				logger.Error("Error Accepting on \"localhost:%d\" -- %v", port, err)
				break;
			}
			defer func() {
				if err := recover(); err != nil {
					const size = 4096
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					log.Printf("dynsoftup: panic serving %v -> %v: %v\n%s", (*chunkContext.clientConnection).LocalAddr(), (*chunkContext.clientConnection).RemoteAddr(), err, buf)
				}
			}()

			var EIGHT = WriteToClient()
			var SEVEN = DetectEndOfData(EIGHT)
			var SIX = SetServerOrigin(SEVEN)
			var FIVE = ReadFromConnection(SIX, chunkContext.serverConnection)
			var FOUR = WriteToServer(FIVE)
			var THREE = RouteToServer(FOUR)
			var TWO = DetectEndOfData(THREE)
			var ONE = ReadFromConnection(TWO, chunkContext.clientConnection)

			ONE(*chunkContext)
		}
	}
}

// ==== ACCEPT LOOP - END

// ==== STAGE ONE / FIVE - START

func ReadFromConnection(next func(chunkContext), connection *net.Conn) func(chunkContext) {

	return func(context chunkContext) {
		logger.Debug("ReadFromConnection " + context.String())
		for !*context.complete {
			size, err := (*connection).Read(context.data);
			if err == io.EOF {
				break;
			}
			if err != nil {
				logger.Error("Error Reading from client %s %+v", context.String(), err)
				return
			}
			context.data = context.data[0:size]
			next(context)
		}
	}
}

// ==== STAGE ONE / FIVE - END

// ==== STAGE TWO / SEVEN - START

var headerRegex = regexp.MustCompile("Content-Length: \\d+")
var contentLengthRegex = regexp.MustCompile("\\d+")

func DetectEndOfData(next func(chunkContext)) func(chunkContext) {
	var (
		contentLengthString = ""
		contentLength       = 0
		receivedContent     = 0
		firstRequestChunk   = true
	)

	return func(context chunkContext) {
		logger.Debug("DetectEndOfData " + context.String())
		var method string
		chunkData := context.data
		chunkSize := len(chunkData)

		if firstRequestChunk {
			switch {
			case chunkData[0] == 'G' && chunkData[1] == 'E':
				method = "GET"
			case chunkData[0] == 'P' && chunkData[1] == 'O':
				method = "POST"
			case chunkData[0] == 'P' && chunkData[1] == 'U':
				method = "PUT"
			}
		}

		*context.complete = false;
		if (method == "GET") {
			*context.complete = true;
		} else {
			if firstRequestChunk {
				contentLengthString = headerRegex.FindString(string(chunkData))
				contentLength, _ = strconv.Atoi(contentLengthRegex.FindString(contentLengthString))
				firstRequestChunk = false
			}


			if (len(contentLengthString) > 0) {
				if receivedContent == 0 {
					// To detect end of headers in first chunk for Content-Length: x
					for i := 0; i < chunkSize; i++ {

						// 0x000D = \r
						// 0x000A = \n
						if (i >= 2 &&
							chunkData[i-2] == 0x000A &&
							chunkData[i-1] == 0x000D &&
							chunkData[i] == 0x000A) {
							// start of body (end of headers)
							receivedContent += (chunkSize-(i+1))
						}
					}
				} else {
					receivedContent += chunkSize
				}
				if receivedContent >= contentLength {
					*context.complete = true;
				}
			} else if chunkSize == 0 ||
					(chunkData[chunkSize-5] == 0x0030 &&
						chunkData[chunkSize-4] == 0x000D &&
						chunkData[chunkSize-3] == 0x000A &&
						chunkData[chunkSize-2] == 0x000D &&
						chunkData[chunkSize-1] == 0x000A) {
				// Check if final chunk for Transfer-Encoding: chunked
				*context.complete = true;
			}
		}
		next(context)
	}
}

// ==== STAGE TWO / SEVEN - END

// ==== STAGE THREE - START

var requestNumber = -1

func RouteToServer(next func(chunkContext)) func(chunkContext) {
	var err error
	requestNumber++

	return func(context chunkContext) {
		logger.Debug("RouteToServer " + context.String())
		if *context.serverConnection, err = net.Dial("tcp", "localhost:"+strconv.Itoa(1024+(requestNumber%2))); err != nil {
			logger.Error("Error Dialing \"localhost:%d\" %s %+v", 1024+(requestNumber%2), context.String(), err)
			//			(*context.clientConnection).Close()
			//			(*context.serverConnection).Close()
			return
		}
		next(context)
	}
}

// ==== STAGE THREE - END

// ==== STAGE FOUR - START

func WriteToServer(next func(chunkContext)) func(chunkContext) {
	var serverReadThreadStarted = false

	return func(requestContext chunkContext) {
		logger.Debug("WriteToServer " + requestContext.String())
		(*requestContext.serverConnection).Write(requestContext.data)
		if !serverReadThreadStarted {
			responseContext := NewChunkContext()
			responseContext.serverConnection = requestContext.serverConnection
			responseContext.clientConnection = requestContext.clientConnection
			go next(*responseContext)
		}
	}
}

// ==== STAGE FOUR - END

// ==== STAGE SIX - START

func SetServerOrigin(next func(chunkContext)) func(chunkContext) {

	return func(context chunkContext) {
		logger.Debug("SetServerOrigin " + context.String())
		// TODO

		next(context)
	}
}

// ==== STAGE SIX - END

// ==== STAGE EIGHT - START

func WriteToClient() func(chunkContext) {
	return func(context chunkContext) {
		logger.Debug("WriteToClient " + context.String())

		(*context.clientConnection).Write(context.data)
		if *context.complete {
			(*context.clientConnection).Close()
			(*context.serverConnection).Close()
		}
	}
}

// ==== STAGE EIGHT - END
