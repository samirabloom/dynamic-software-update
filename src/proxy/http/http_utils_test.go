package http

import (
	"net"
	"strconv"
	"testing"
	networkutil "util/test/network"
	assertion "util/test/assertion"
	"fmt"
)

func Test_Update_Header_For_Requests_To_Local_Addresses(testCtx *testing.T) {
	// given
	var (
		data            = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: www.google.co.uk\nAccept: */*")
		to *net.TCPAddr = networkutil.FindFreeLocalSocket(testCtx)
		expectedData    = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: 127.0.0.1:" + strconv.Itoa(to.Port) + "\nAccept: */*")
	)

	// when
	actualData := UpdateHostHeader(data, to)

	// then
	assertion.AssertDeepEqual("Correctly added Host header", testCtx, expectedData, actualData)
}

func XTest_Update_Header_For_Requests_To_Remote_Addresses(testCtx *testing.T) {
	// given
	var (
		addrs, _        = net.LookupIP("www.google.co.uk")
		data            = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: www.google.co.uk\nAccept: */*")
		to *net.TCPAddr = &net.TCPAddr{IP: addrs[0], Port: 80}
		expectedData    = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: www.google.co.uk:80\nAccept: */*")
	)
	fmt.Printf("to: %s %#v\n", to, to)

	// when
	actualData := UpdateHostHeader(data, to)

	// then
	assertion.AssertDeepEqual("Correctly added Host header", testCtx, expectedData, actualData)
}
