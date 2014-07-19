package proxy_c

import (
	"net"
	"time"
	"testing"
	networkutil "util/test/network"
	assertion "util/test/assertion"
	"code.google.com/p/go-uuid/uuid"
)

func NewLoadBalancer(frontendAddr *net.TCPAddr, backendAddresses []*net.TCPAddr) *LoadBalancer {
	return &LoadBalancer{
		frontendAddr: frontendAddr,
		router: &RoutingContext{
			backendAddresses:  backendAddresses,
			requestCounter: -1,
			uuid: uuid.NewUUID(),
		},
		stop: make(chan bool),
	}
}

func sendRequest(testCtx *testing.T, address *net.TCPAddr, request []byte) []byte {
	//   - a socket connected to proxy
	proxyConnection, err := net.DialTCP("tcp", nil, address); if err != nil {
		testCtx.Fatalf("Can't connect to the proxy: %v", err)
	} else {
		defer proxyConnection.Close()
		proxyConnection.SetDeadline(time.Now().Add(10 * time.Second))
	}

	if _, err = proxyConnection.Write(request); err != nil {
		testCtx.Fatal(err)
	}
	var size = 0
	recvBuf := make([]byte, len(request)*3)
	if size, err = proxyConnection.Read(recvBuf); err != nil {
		testCtx.Fatal(err)
	}

	return recvBuf[0:size]
}

func Test_Proxy_Basic_Request_And_Response(testCtx *testing.T) {
	// given
	//   - echo server running
	var (
		echoServerAddress = networkutil.FindFreeLocalSocket(testCtx)
		proxyAddress      = networkutil.FindFreeLocalSocket(testCtx)
	)
	networkutil.Run(testCtx, echoServerAddress)

	//   - proxy running
	var (
		proxy = NewLoadBalancer(proxyAddress, []*net.TCPAddr{echoServerAddress})
	)
	proxy.Start()

	//   - a example request
	var (
		testRequest      = []byte("some random request with no new lines")
		expectedResponse = []byte("some random request with no new lines")
	)

	// when
	response := sendRequest(testCtx, proxyAddress, testRequest)

	// then
	assertion.AssertDeepEqual("Correct Response", testCtx, expectedResponse, response)

	// clean-up
	proxy.Stop()
}

func Test_Proxy_Request_With_UUID(testCtx *testing.T) {
	// given
	//   - echo server running
	var (
		echoServerAddress = networkutil.FindFreeLocalSocket(testCtx)
		proxyAddress      = networkutil.FindFreeLocalSocket(testCtx)
	)
	networkutil.Run(testCtx, echoServerAddress)

	//   - proxy running
	var (
		proxy = NewLoadBalancer(proxyAddress, []*net.TCPAddr{echoServerAddress})
	)
	proxy.Start()

	//   - a example request
	var (
		testRequest      = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: www.test.co.uk\nAccept: */*\nAccept-Encoding: deflate, gzip\nCookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245a;")
		expectedResponse = []byte("GET / HTTP/1.1\nSet-Cookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245a;\nX-EchoServer: " + echoServerAddress.String() + "\nUser-Agent: curl/7.30.0\nHost: www.test.co.uk\nAccept: */*\nAccept-Encoding: deflate, gzip\nCookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245a;")
	)

	// when
	response := sendRequest(testCtx, proxyAddress, testRequest)

	// then
	assertion.AssertDeepEqual("Correct Response", testCtx, expectedResponse, response)

	// clean-up
	proxy.Stop()
}

func Test_Proxy_Load_Balances_Request(testCtx *testing.T) {
	// given
	//   - echo server running
	var (
		proxyAddress        = networkutil.FindFreeLocalSocket(testCtx)
		echoServerAddresses = networkutil.FindFreeLocalSocketRange(testCtx, 2, 10)
	)
	networkutil.Run(testCtx, echoServerAddresses[0])
	networkutil.Run(testCtx, echoServerAddresses[1])

	//   - proxy running
	var (
		proxy = NewLoadBalancer(proxyAddress, echoServerAddresses)
	)
	proxy.Start()

	//   - a example request
	var (
		testRequest         = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: www.test.co.uk\nAccept: */*\nAccept-Encoding: deflate, gzip\nCookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245a;")
		expectedResponseOne = []byte("GET / HTTP/1.1\nSet-Cookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245a;\nX-EchoServer: " + echoServerAddresses[0].String() + "\nUser-Agent: curl/7.30.0\nHost: www.test.co.uk\nAccept: */*\nAccept-Encoding: deflate, gzip\nCookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245a;")
		expectedResponseTwo = []byte("GET / HTTP/1.1\nSet-Cookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245a;\nX-EchoServer: " + echoServerAddresses[1].String() + "\nUser-Agent: curl/7.30.0\nHost: www.test.co.uk\nAccept: */*\nAccept-Encoding: deflate, gzip\nCookie: dynsoftup=0e5e6c61-0731-11e4-aaec-600308a8245a;")
	)

	// when
	// - first request
	responseOne := sendRequest(testCtx, proxyAddress, testRequest)
	responseTwo := sendRequest(testCtx, proxyAddress, testRequest)
	responseThree := sendRequest(testCtx, proxyAddress, testRequest)


	// then
	assertion.AssertDeepEqual("Correct Response", testCtx, expectedResponseOne, responseOne)
	assertion.AssertDeepEqual("Correct Response", testCtx, expectedResponseTwo, responseTwo)
	assertion.AssertDeepEqual("Correct Response", testCtx, expectedResponseOne, responseThree)

	// clean-up
	proxy.Stop()
}

// TODO:
// routing - with cookie send to specific server
//         - add cookies to match server that provided response
// handle incorrect response (by re-routing request)
// handle incorrect response (by only returning correct response from multiple servers)

