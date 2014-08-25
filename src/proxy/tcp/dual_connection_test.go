package tcp

import (
	networkutil "util/test/network"
	assertion "util/test/assertion"
	"net"
	"strconv"
	"net/http"
	"time"
	"testing"
	"io"
	"fmt"
)

type SimpleHandler struct {
	Port         int
	ResponseCode int
}

func (handler *SimpleHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Date", "Sat, 02 Aug 2014 20:11:40 GMT")
	response.WriteHeader(handler.ResponseCode)
	fmt.Fprintf(response, "Response from the server with port: %d and status code: %d", handler.Port, handler.ResponseCode)
}

func StartDualServers(addresses []*net.TCPAddr, statusCode []int) {
	for index, address := range addresses {
		go http.ListenAndServe(":"+strconv.Itoa(address.Port), &SimpleHandler{address.Port, statusCode[index]})
	}
	time.Sleep(150 * time.Millisecond)
}

func Test_Dual_TCP_Connection_When_First_Response_Invalid(testCtx *testing.T) {
	// given
	var (
		servers            = []*net.TCPAddr{networkutil.FindFreeLocalSocket(testCtx), networkutil.FindFreeLocalSocket(testCtx)}
		hosts              = []string{"127.0.0.1", "127.0.0.1"}
		ports              = []string{fmt.Sprintf("%v", servers[0].Port), fmt.Sprintf("%v", servers[1].Port)}
		actualResponse     = make([]byte, 0)
		actualTotalRead    = 0
		expectedResponse   = []byte("HTTP/1.1 200 OK\r\n" +
			"Date: Sat, 02 Aug 2014 20:11:40 GMT\r\n" +
			"Content-Length: " + strconv.Itoa(57 + len(strconv.Itoa(servers[1].Port))) + "\r\n" +
			"Content-Type: text/plain; charset=utf-8\r\n" + "\r\n" +
			"Response from the server with port: " + strconv.Itoa(servers[1].Port) + " and status code: 200")
		expectedTotalRead  = 179
	)

	StartDualServers(servers, []int{http.StatusInternalServerError, http.StatusOK})
	dualConnection := NewDualTCPConnection(servers, hosts, ports)

	// when
	dualConnection.Write([]byte("POST / HTTP/1.1\n" +
			"Content-Length: 10\n" +
			"\r\n" +
			"some_body\n"))
	dualConnection.CloseWrite()

	for {
		readBuffer := make([]byte, 4096)
		lenRead, err := dualConnection.Read(readBuffer)
		actualResponse = append(actualResponse, readBuffer...)
		actualTotalRead += lenRead

		if err == io.EOF {
			dualConnection.Close()
			break
		} else if err != nil {
			testCtx.Fatalf("Unexpected error while reading reponse %s", err)
		}
	}

	// then
	assertion.AssertDeepEqual("Correct response", testCtx, expectedResponse, actualResponse[0:actualTotalRead])
	assertion.AssertDeepEqual("Correct length read", testCtx, expectedTotalRead, actualTotalRead)
}

func Test_Dual_TCP_Connection_When_Second_Response_Invalid(testCtx *testing.T) {
	// given
	var (
		servers            = []*net.TCPAddr{networkutil.FindFreeLocalSocket(testCtx), networkutil.FindFreeLocalSocket(testCtx)}
		hosts              = []string{"127.0.0.1", "127.0.0.1"}
		ports              = []string{fmt.Sprintf("%v", servers[0].Port), fmt.Sprintf("%v", servers[1].Port)}
		actualResponse     = make([]byte, 0)
		actualTotalRead    = 0
		expectedResponse   = []byte("HTTP/1.1 200 OK\r\n" +
			"Date: Sat, 02 Aug 2014 20:11:40 GMT\r\n" +
			"Content-Length: " + strconv.Itoa(57 + len(strconv.Itoa(servers[0].Port))) + "\r\n" +
			"Content-Type: text/plain; charset=utf-8\r\n" + "\r\n" +
			"Response from the server with port: " + strconv.Itoa(servers[0].Port) + " and status code: 200")
		expectedTotalRead  = 179
	)

	StartDualServers(servers, []int{http.StatusOK, http.StatusInternalServerError})
	dualConnection := NewDualTCPConnection(servers, hosts, ports)

	// when
	dualConnection.Write([]byte("POST / HTTP/1.1\n" +
			"Content-Length: 10\n" +
			"\r\n" +
			"some_body\n"))
	dualConnection.CloseWrite()

	for {
		readBuffer := make([]byte, 4096)
		lenRead, err := dualConnection.Read(readBuffer)
		actualResponse = append(actualResponse, readBuffer...)
		actualTotalRead += lenRead

		if err == io.EOF {
			dualConnection.Close()
			break
		} else if err != nil {
			testCtx.Fatalf("Unexpected error while reading reponse %s", err)
		}
	}

	// then
	assertion.AssertDeepEqual("Correct response", testCtx, expectedResponse, actualResponse[0:actualTotalRead])
	assertion.AssertDeepEqual("Correct length read", testCtx, expectedTotalRead, actualTotalRead)
}

func Test_Dual_TCP_Connection_When_Both_Responses_Valid(testCtx *testing.T) {
	// given
	var (
		servers            = []*net.TCPAddr{networkutil.FindFreeLocalSocket(testCtx), networkutil.FindFreeLocalSocket(testCtx)}
		hosts              = []string{"127.0.0.1", "127.0.0.1"}
		ports              = []string{fmt.Sprintf("%v", servers[0].Port), fmt.Sprintf("%v", servers[1].Port)}
		actualResponse     = make([]byte, 0)
		actualTotalRead    = 0
		expectedResponse   = []byte("HTTP/1.1 200 OK\r\n" +
			"Date: Sat, 02 Aug 2014 20:11:40 GMT\r\n" +
			"Content-Length: " + strconv.Itoa(57 + len(strconv.Itoa(servers[0].Port))) + "\r\n" +
			"Content-Type: text/plain; charset=utf-8\r\n" + "\r\n" +
			"Response from the server with port: " + strconv.Itoa(servers[1].Port) + " and status code: 200")
		expectedTotalRead  = 179
	)

	StartDualServers(servers, []int{http.StatusOK, http.StatusOK})
	dualConnection := NewDualTCPConnection(servers, hosts, ports)

	// when
	dualConnection.Write([]byte("POST / HTTP/1.1\n" +
			"Content-Length: 10\n" +
			"\r\n" +
			"some_body\n"))
	dualConnection.CloseWrite()

	for {
		readBuffer := make([]byte, 4096)
		lenRead, err := dualConnection.Read(readBuffer)
		actualResponse = append(actualResponse, readBuffer...)
		actualTotalRead += lenRead

		if err == io.EOF {
			dualConnection.Close()
			break
		} else if err != nil {
			testCtx.Fatalf("Unexpected error while reading reponse %s", err)
		}
	}

	// then
	assertion.AssertDeepEqual("Correct response", testCtx, expectedResponse, actualResponse[0:actualTotalRead])
	assertion.AssertDeepEqual("Correct length read", testCtx, expectedTotalRead, actualTotalRead)
}


func Test_Dual_TCP_Connection_When_Both_Responses_Invalid(testCtx *testing.T) {
	// given
	var (
		servers            = []*net.TCPAddr{networkutil.FindFreeLocalSocket(testCtx), networkutil.FindFreeLocalSocket(testCtx)}
		hosts              = []string{"127.0.0.1", "127.0.0.1"}
		ports              = []string{fmt.Sprintf("%v", servers[0].Port), fmt.Sprintf("%v", servers[1].Port)}
		actualResponse     = make([]byte, 0)
		actualTotalRead    = 0
		expectedResponse   = []byte("HTTP/1.1 500 Internal Server Error\r\n" +
			"Date: Sat, 02 Aug 2014 20:11:40 GMT\r\n" +
			"Content-Length: " + strconv.Itoa(57 + len(strconv.Itoa(servers[0].Port))) + "\r\n" +
			"Content-Type: text/plain; charset=utf-8\r\n" + "\r\n" +
			"Response from the server with port: " + strconv.Itoa(servers[0].Port) + " and status code: 500")
		expectedTotalRead  = 198
	)

	StartDualServers(servers, []int{http.StatusInternalServerError, http.StatusInternalServerError})
	dualConnection := NewDualTCPConnection(servers, hosts, ports)

	// when
	dualConnection.Write([]byte("POST / HTTP/1.1\n" +
			"Content-Length: 10\n" +
			"\r\n" +
			"some_body\n"))
	dualConnection.CloseWrite()

	for {
		readBuffer := make([]byte, 4096)
		lenRead, err := dualConnection.Read(readBuffer)
		actualResponse = append(actualResponse, readBuffer...)
		actualTotalRead += lenRead

		if err == io.EOF {
			dualConnection.Close()
			break
		} else if err != nil {
			testCtx.Fatalf("Unexpected error while reading reponse %s", err)
		}
	}

	// then
	assertion.AssertDeepEqual("Correct response", testCtx, expectedResponse, actualResponse[0:actualTotalRead])
	assertion.AssertDeepEqual("Correct length read", testCtx, expectedTotalRead, actualTotalRead)
}
