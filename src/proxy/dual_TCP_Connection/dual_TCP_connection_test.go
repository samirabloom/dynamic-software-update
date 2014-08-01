package dual_TCP_connection

import (
	networkutil "util/test/network"
	log "proxy/log"
	assertion "util/test/assertion"
	"net"
	"strconv"
	"net/http"
	"time"
	"testing"
	"io"
	"regexp"
	"bytes"
)

func Test_Dual_TCP_Connection_And_Response_From_StatusCode_200(testCtx *testing.T) {
	// given
	var (
		serverOne, _       = net.ResolveTCPAddr("tcp", "127.0.0.1:9098")
		serverTwo, _       = net.ResolveTCPAddr("tcp", "127.0.0.1:9099")
		requestBody = "<html><head></head><body><div>This is a response from two different server</div></body></html>\n"
		serverResponse = make([]byte, 4096)
		numberOfResponses  = 0
		serverResponses    = make([][]byte, 1024)
		expectedStatusCode = 200
		actualStatusCode   int
		expectedResponseBody = []byte("Server responded with status code: 200 and Port: 9099\n")
		actualResponseBody []byte
		statusCodeRegex = regexp.MustCompile("HTTP/[0-9].[0-9] ([a-z0-9-]*) .*")
	)

	go http.ListenAndServe(":"+strconv.Itoa(9098), &networkutil.Handle1{9098})
	time.Sleep(2000 * time.Millisecond)
	go http.ListenAndServe(":"+strconv.Itoa(9099), &networkutil.Handle2{9099})
	time.Sleep(2000 * time.Millisecond)
	dualConnection := NewDualTCPConnection(expectedStatusCode, serverOne, serverTwo)

	// when
	dualConnection.Write([]byte("POST / HTTP/1.1\n" +
			"Content-Length: " + strconv.Itoa(len(requestBody)) + "\n" +
			"\r\n" +
			requestBody))
	dualConnection.CloseWrite()

	for {
		numberOfResponses++
		lenRead, err := dualConnection.Read(serverResponse)
		serverResponses[numberOfResponses] = serverResponse[0:lenRead]
		if err == io.EOF {
			log.LoggerFactory().Debug("Read Loop EOF - %s", dualConnection)
			break
		} else if err != nil {
			log.LoggerFactory().Error("error reading the response from dual connection %v:\n", err)
			return // terminate program
		}
	}
	dualConnection.CloseRead()
	dualConnection.Close()

	// then

		// get response status code
	statusCodeMatches := statusCodeRegex.FindSubmatch(serverResponses[1])
	if len(statusCodeMatches) >= 2 {
		statusCodeMatch := string(statusCodeMatches[1])
		actualStatusCode, _ = strconv.Atoi(statusCodeMatch)
	}

		// get response body
	bodyStartPosition := bytes.Index(serverResponses[1], []byte("Server"))
	if bodyStartPosition > 0 {
		actualResponseBody = serverResponse[bodyStartPosition:len(serverResponses[1])]
	}
		// check the status code and response body
	assertion.AssertDeepEqual("Correct status code ", testCtx, expectedStatusCode, actualStatusCode)
	assertion.AssertDeepEqual("Correct response body ", testCtx, expectedResponseBody, actualResponseBody)
}


func Test_Dual_TCP_Connection_And_Response_From_StatusCode_500(testCtx *testing.T) {
	// given
	var (
		serverOne, _       = net.ResolveTCPAddr("tcp", "127.0.0.1:9098")
		serverTwo, _       = net.ResolveTCPAddr("tcp", "127.0.0.1:9099")
		requestBody = "<html><head></head><body><div>This is a response from two different server</div></body></html>\n"
		serverResponse = make([]byte, 4096)
		numberOfResponses  = 0
		serverResponses    = make([][]byte, 1024)
		expectedStatusCode = 500
		actualStatusCode   int
		expectedResponseBody = []byte("Response from the server with status code: 500 and Port: 9098\n")
		actualResponseBody []byte
		statusCodeRegex = regexp.MustCompile("HTTP/[0-9].[0-9] ([a-z0-9-]*) .*")
	)

	go http.ListenAndServe(":"+strconv.Itoa(9098), &networkutil.Handle1{9098})
	time.Sleep(2000 * time.Millisecond)
	go http.ListenAndServe(":"+strconv.Itoa(9099), &networkutil.Handle2{9099})
	time.Sleep(2000 * time.Millisecond)
	dualConnection := NewDualTCPConnection(expectedStatusCode, serverOne, serverTwo)

	// when
	dualConnection.Write([]byte("POST / HTTP/1.1\n" +
			"Content-Length: " + strconv.Itoa(len(requestBody)) + "\n" +
			"\r\n" +
			requestBody))
	dualConnection.CloseWrite()

	for {
		numberOfResponses++
		lenRead, err := dualConnection.Read(serverResponse)
		serverResponses[numberOfResponses] = serverResponse[0:lenRead]
		if err == io.EOF {
			log.LoggerFactory().Debug("Read Loop EOF - %s", dualConnection)
			break
		} else if err != nil {
			log.LoggerFactory().Error("error reading the response from dual connection %v:\n", err)
			return // terminate program
		}
	}
	dualConnection.CloseRead()
	dualConnection.Close()

	// then

	// get response status code
	statusCodeMatches := statusCodeRegex.FindSubmatch(serverResponses[1])
	if len(statusCodeMatches) >= 2 {
		statusCodeMatch := string(statusCodeMatches[1])
		actualStatusCode, _ = strconv.Atoi(statusCodeMatch)
	}

	// get response body
	bodyStartPosition := bytes.Index(serverResponses[1], []byte("Response"))
	if bodyStartPosition > 0 {
		actualResponseBody = serverResponse[bodyStartPosition:len(serverResponses[1])]
	}
	// check the status code and response body
	assertion.AssertDeepEqual("Correct status code ", testCtx, expectedStatusCode, actualStatusCode)
	assertion.AssertDeepEqual("Correct response body ", testCtx, expectedResponseBody, actualResponseBody)
}


