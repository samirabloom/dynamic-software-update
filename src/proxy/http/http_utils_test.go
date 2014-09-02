package http

import (
	"testing"
	assertion "util/test/assertion"
)

func Test_Update_Header_For_Requests_To_Local_Port_80_Address(testCtx *testing.T) {
	// given
	var (
		data            = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: www.google.co.uk\nAccept: */*")
		expectedData    = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: 127.0.0.1\nAccept: */*")
	)

	// when
	actualData := UpdateHostHeader(data, "127.0.0.1", "80", true)

	// then
	assertion.AssertDeepEqual("Correctly added Host header", testCtx, expectedData, actualData)
}

func Test_Update_Header_For_Requests_To_Remote_None_Port_80_Address(testCtx *testing.T) {
	// given
	var (
		data            = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: www.google.co.uk\nAccept: */*")
		expectedData    = []byte("GET / HTTP/1.1\nUser-Agent: curl/7.30.0\nHost: www.google.co.uk:443\nAccept: */*")
	)

	// when
	actualData := UpdateHostHeader(data, "www.google.co.uk", "443", true)

	// then
	assertion.AssertDeepEqual("Correctly added Host header", testCtx, expectedData, actualData)
}
