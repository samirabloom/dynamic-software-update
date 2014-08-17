package zmq_proxy

import (
	"os"
	zmq "github.com/pebbe/zmq4"
	"time"
	"regexp"
	"strconv"
	"fmt"
)

const (
	REQUEST_TIMEOUT = 10 * time.Millisecond //  msecs
)

// mkdir /tmp/feeds
// proxy 2 9090 1080 1081

// define ZMQ response handler
func CLI() {
	logger := log()

	logger(0, "\nForwarding: "+os.Args[2]+" -> 127.0.0.1:"+os.Args[3]+"/"+os.Args[4]+"\n")

	go router(logger)
	go upstream(logger, os.Args[2])
	go downstream(logger, "127.0.0.1:"+os.Args[3], "4")
	go downstream(logger, "127.0.0.1:"+os.Args[4], "3")
	go downstream(logger, "127.0.0.1:"+os.Args[3], "2")
	go downstream(logger, "127.0.0.1:"+os.Args[4], "1")

	//  Run for 20 minutes then quit
	time.Sleep(1200 * time.Second)
}

func log() func(level int, format string, a... interface{}) {
	// ERROR = 4, WARN = 3, INFO = 2, DEBUG = 1, TRACE = 0
	logLevel, _ := strconv.Atoi(os.Args[1])
	return func(level int, format string, a... interface{}) {
		if (level >= logLevel) {
			fmt.Printf(format, a...)
		}
	}
}

func upstream(log func(level int, format string, a... interface{}), upstreamPort string) {
	//  socket for incoming HTTP request/response
	upstream, error := zmq.NewSocket(zmq.STREAM)
	if error != nil {
		log(3, "error creating upstream socket: %s\n", error)
	} else {
		defer upstream.Close()
		log(0, "created upstream socket\n")
	}
	error = upstream.Bind("tcp://*:"+upstreamPort)
	if error != nil {
		log(3, "error binding [tcp://*:%d] upstream socket: %s\n", upstreamPort, error)
	} else {
		log(0, "bound upstream socket to [%s]\n", upstreamPort)
	}

	//  inproc socket for sending requests to downstream
	backend, error := zmq.NewSocket(zmq.REQ)
	if error != nil {
		log(3, "error creating backend socket: %s\n", error)
	} else {
		defer backend.Close()
		log(0, "created backend socket\n")
	}
	error = backend.Connect("ipc:///tmp/feeds/upstream")
	if error != nil {
		log(3, "error connecting to [ipc:///tmp/feeds/upstream] backend socket: %s\n", error)
	} else {
		log(0, "bound backend socket to [ipc:///tmp/feeds/upstream]\n")
	}

	upstreamPoller := zmq.NewPoller()
	upstreamPoller.Add(upstream, zmq.POLLIN)

	backendPoller := zmq.NewPoller()
	backendPoller.Add(backend, zmq.POLLIN)

	for {
		polled, error := upstreamPoller.Poll(REQUEST_TIMEOUT)

		if error == nil && len(polled) == 1 {
			id, error := upstream.Recv(0)
			if error != nil {
				log(3, "error receiveing request id from upstream: %s\n", error)
			} else {
				log(0, "\n\nreceived upstream request id message [% x]\n", id)

				// send request to backend
				backend.Send(id, zmq.SNDMORE)
			}
			message, error := upstream.Recv(0)
			if error != nil {
				log(3, "error receiveing request from upstream: %s\n", error)
			} else {
				log(0, "received request from upstream [%s]\n", message)

				// send request to backend
				backend.Send(message, 0)
			}
		}

		polled, error = backendPoller.Poll(REQUEST_TIMEOUT)

		if error == nil && len(polled) == 1 {
			id, error := backend.Recv(0)
			if error != nil {
				log(3, "error receiveing request id from upstream: %s\n", error)
			} else {
				log(0, "\n\nreceived upstream request id message [% x]\n", id)
			}
			message, error := backend.Recv(0)
			if error != nil {
				log(3, "error receiveing response from backend: %s\n", error)
			} else {
				log(1, "\n\nsending response id upstream [% x]\n", id)
				log(0, "sending response message upstream [\n"+message+"]")
				if (message != "CLOSE") {
					// send response upstream
					upstream.Send(id, zmq.SNDMORE)
					upstream.Send(message, zmq.SNDMORE)
				}
			}
		}
	}
}

func router(log func(level int, format string, a... interface{})) {
	// inproc socket for receiving requests from frontend
	frontend, error := zmq.NewSocket(zmq.ROUTER)
	if error != nil {
		log(3, "error creating frontend socket: %s\n", error)
	} else {
		defer frontend.Close()
		log(0, "created frontend socket\n")
	}
	error = frontend.Bind("ipc:///tmp/feeds/upstream")
	if error != nil {
		log(3, "error binding [ipc:///tmp/feeds/upstream] frontend socket: %s\n", error)
	} else {
		log(0, "bound frontend socket to [ipc:///tmp/feeds/upstream]\n")
	}

	//  inproc socket for sending requests to downstream
	backend, error := zmq.NewSocket(zmq.DEALER)
	if error != nil {
		log(3, "error creating backend socket: %s\n", error)
	} else {
		defer backend.Close()
		log(0, "created backend socket\n")
	}
	error = backend.Bind("ipc:///tmp/feeds/downstream")
	if error != nil {
		log(3, "error binding [ipc:///tmp/feeds/downstream] backend socket: %s\n", error)
	} else {
		log(0, "bound backend socket to [ipc:///tmp/feeds/downstream]\n")
	}

	//  Connect backend to frontend via a proxy
	err := zmq.Proxy(frontend, backend, nil)
	log(0, "Proxy interrupted: %s\n", err)
}

func downstream(log func(level int, format string, a... interface{}), downstreamHostAndPort string, routerId string) {
	// inproc socket for receiving requests from frontend
	frontend, error := zmq.NewSocket(zmq.REP)
	if error != nil {
		log(3, "\n%s - error creating frontend socket: %s\n", routerId, error)
	} else {
		defer frontend.Close()
		log(0, "\n%s - created frontend socket\n", routerId)
	}
	error = frontend.Connect("ipc:///tmp/feeds/downstream")
	if error != nil {
		log(3, "\n%s - error connecting to [ipc:///tmp/feeds/downstream] frontend socket: %s\n", routerId, error)
	} else {
		log(0, "\n%s - bound frontend socket to [ipc:///tmp/feeds/downstream]\n", routerId)
	}

	// socket for downstream HTTP request/response
	downstream, error := zmq.NewSocket(zmq.STREAM)
	if error != nil {
		log(3, "\n%s - error creating downstream socket: %s\n", routerId, error)
	} else {
		defer downstream.Close()
		log(0, "\n%s - created downstream socket\n", routerId)
	}

	for {
		// receive clientId
		clientId, error := frontend.Recv(0)
		if error != nil {
			log(3, "\n%s - error receiveing clientId from frontend: %s\n", routerId, error)
		} else {
			log(2, "\n%s - received frontend clientId message [% x]\n", routerId, clientId)
		}

		// receive message
		message, error := frontend.Recv(0)
		if error != nil {
			log(3, "\n%s - error receiveing message from frontend: %s\n", routerId, error)
		} else {
			log(0, "\n%s - received message from frontend [%s]\n", routerId, message)

			error = downstream.Connect("tcp://"+downstreamHostAndPort)
			if error != nil {
				log(3, "\n%s - error connecting to [tcp://%d] downstream socket: %s\n", routerId, downstreamHostAndPort, error)
			} else {
				log(0, "\n%s - bound downstream socket to [%s]\n", routerId, downstreamHostAndPort)
			}

			// disable compression to simplify proxy
			message = regexp.MustCompile("Accept-Encoding: [a-z\\,]*").ReplaceAllLiteralString(message, "")
			// disable persistent connections to simplify proxy
			// this avoids HTTP response chunks from arriving in the same ZeroMQ message for different requests
			message = regexp.MustCompile("Connection: keep-alive").ReplaceAllLiteralString(message, "Connection: close")

			// ensure that the request has the correct Host header
			message = regexp.MustCompile("Host: [a-z\\,\\.\\:\\d]*").ReplaceAllLiteralString(message, "Host: "+downstreamHostAndPort)


			serverId, error := downstream.GetIdentity()
			if error != nil {
				log(3, "\n%s - error getting serverId: %s\n", routerId, error)
			} else {
				log(2, "\n%s - received serverId [% x]\n", routerId, serverId)
			}

			log(1, "\n%s - sending serverId downstream [% x]\n", routerId, serverId)
			log(0, "\n%s - sending message downstream [%s]\n", routerId, message)

			// send response message - 1st chunk
			downstream.Send(serverId, zmq.SNDMORE)
			downstream.Send(message, 0)

			contentReceived := 0
			contentLength := 0
			contentLengthHeader := ""
			transferEncodingHeader := ""
			for {
				// receive serverId
				serverId, error := downstream.Recv(0)
				if error != nil {
					log(3, "\n%s - error receiveing downstream response serverId message: %s\n", routerId, error)
				} else {
					log(0, "\n%s - received downstream response serverId message [% x]\n", routerId, serverId)
				}

				// receive message
				message, error := downstream.Recv(0)
				if error != nil {
					log(3, "\n%s - error receiveing downstream response message: %s\n", routerId, error)
				} else {
					log(0, "\n%s - received downstream message [%s]\n", routerId, message)

					// send message to frontend
					frontend.Send(clientId, zmq.SNDMORE)
					frontend.Send(message, 0)
				}

				if len(contentLengthHeader) == 0 && len(transferEncodingHeader) == 0 {
					contentLengthHeader = regexp.MustCompile("Content-Length: \\d+").FindString(message)
					contentLength, _ = strconv.Atoi(regexp.MustCompile("\\d+").FindString(contentLengthHeader))
					transferEncodingHeader = regexp.MustCompile("Transfer-Encoding: chunked").FindString(message)
				}


				if len(contentLengthHeader) > 0 { // work out when response has been fully received (i.e. content received == content length header)

					if (contentReceived == 0) {
						// first chunk with headers
						for i := 0; i < len(message); i++ {
							if (i >= 4 &&
								string(message[i - 4]) == string("\u000D") &&
								string(message[i - 3]) == string("\u000A") &&
								string(message[i - 2]) == string("\u000D") &&
								string(message[i - 1]) == string("\u000A")) {
								// start of body (end of headers)
								contentReceived += (len(message)-i)
							}
						}

					} else {
						// not first chunk so no header
						contentReceived += len(message)
					}

					log(0, "\n%s - contentReceived: %d\n", routerId, contentReceived)
					log(0, "\n%s - contentLength: %d\n", routerId, contentLength)

					if contentReceived >= contentLength {
						// received all content
						break;
					}


				} else if len(transferEncodingHeader) > 0 { // work out when response has been fully received (i.e. final chunk has been received)

					// check if downstream has sent final empty chunk (as per HTTP specification)
					if len(message) == 0 {
						// received final chunk
						break;
					}

					// check if downstream sends trailer indicating final chunk, as follows:
					// -5 U+0030 - '0'
					// -4 U+000D - Carriage return
					// -3 U+000A - Line feed
					// -2 U+000D - Carriage return
					// -1 U+000A - Line feed
					// this allows for downstream which doesn't conform strictly to HTTP specification
					length := len(message)
					if (string(message[length - 5]) == string("\u0030") &&
						string(message[length - 4]) == string("\u000D") &&
						string(message[length - 3]) == string("\u000A") &&
						string(message[length - 2]) == string("\u000D") &&
						string(message[length - 1]) == string("\u000A")) {
						// received final chunk
						break;
					}
				} else {
					log(3, "\n%s - Error can't find either Content-Length or Transfer-Encoding headers\n", routerId)
				}
			}

			// close connection

			// send empty chunk
			downstream.Send(serverId, zmq.SNDMORE)
			downstream.Send("", 0)

			log(2, "\n%s - response completed closing HTTP connection\n", routerId)
		}
	}
}
