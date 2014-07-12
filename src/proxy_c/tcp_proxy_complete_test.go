package proxy_c

import (
	"testing"
	"net"
	"fmt"
	"syscall"
)

func Test_On_Complete_With_(testCtx *testing.T) {
	// given
	var (
		source      = &net.TCPConn{}
		destination = &net.TCPConn{}
		mockContext = &chunkContext{
			description: "",
			data: make([]byte, 64*1024),
			from: source,
			to: destination,
			err: syscall.EPIPE,
			totalReadSize: 0,
			totalWriteSize: 0,
			event: make(chan int64, 100),
			firstChunk: true,
			performance: *&performance{
				read: new(int64),
				route: new(int64),
				write: new(int64),
				complete: new(int64),
		},
	}
		//		expectedResult error =
	)

	// when
	fmt.Printf("Before calling the complete function [%v]\n", mockContext.err)
	complete(mockContext)
	//	fmt.Printf("After calling the complete function [%v]\n", mockContext.err)

	// then
//	if mockContext.err != io.ErrShortWrite {
//		testCtx.Fatalf("expected: [%v] actual: [%v]", io.ErrShortWrite, mockContext.err)
//	}

}


// TEST - START

//func TestFoo() {
//	var testValue
//
//	var mockFunction = func(someParameter string) {
//		testValue = someParameter
//	}
//
//	Foo(mockFunction)
//
//	if(testValue != "hello") {
//		// test failed
//	}
//}
//
//func Foo(next func(someParameter string)) {
//	// do something
//	next("hello")
//	// do something else
//}

// - END
