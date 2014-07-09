package proxy_b

import (
	"net"
	"testing"
	"time"
	"bytes"
	"fmt"
	networkutil "util/test/network"
)

func TestShouldProxyRequestAndResponse(testCtx *testing.T) {
	// given
	//   - echo server running
	echoServerAddress, echoServerTCPAddress := networkutil.FindFreeLocalSocket(testCtx)
	_, proxyTCPAddress := networkutil.FindFreeLocalSocket(testCtx)
	networkutil.Run(testCtx, echoServerAddress)

	//   - proxy running
	proxy, _ := NewLoadBalancer(
		proxyTCPAddress,
		echoServerTCPAddress,
		1,
	)
	proxy.uuidGenerator = func() string {
		return "uuid"
	}
	go proxy.Run()

	//   - a socket connected to proxy
	proxyConnection, err := net.DialTCP("tcp", nil, proxyTCPAddress)
	if err != nil {
		testCtx.Fatalf("Can't connect to the proxy: %v", err)
	}
	defer proxyConnection.Close()
	proxyConnection.SetDeadline(time.Now().Add(10 * time.Second))

	//   - a example request
	testRequest := []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: www.test.co.uk\nAccept: */*\nAccept-Encoding: deflate, gzip")
	expectedResponse := []byte("GET / HTTP/1.1\nSet-Cookie: dynsoftup=uuid;\nX-EchoServer: " + echoServerAddress + "\nUser-Agent: curl/7.30.0\nHost: www.test.co.uk\nAccept: */*\nAccept-Encoding: deflate, gzip")

	// when
	if _, err = proxyConnection.Write(testRequest); err != nil {
		testCtx.Fatal(err)
	}
	var size = 0
	recvBuf := make([]byte, len(testRequest)*3)
	if size, err = proxyConnection.Read(recvBuf); err != nil {
		testCtx.Fatal(err)
	}


	// then
//	networkutil.AssertArrayEquals(testCtx, testRequest, recvBuf[0:size])
	if !bytes.Equal(expectedResponse, recvBuf[0:size]) {
		testCtx.Fatal(fmt.Errorf("\nExpected:\n[%s]\nActual:\n[%s]", expectedResponse, recvBuf[0:size]))
	}

	// clean-up
	go proxy.Stop()
}

// load balancing
// routing - with cookie send to specific server
//         - add cookies to match server that provided response
// handle incorrect response (by re-routing request)
// handle incorrect response (by only returning correct response from multiple servers)

