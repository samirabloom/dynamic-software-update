package proxy_c

import (
	"testing"
	"util/test/mock"
	"fmt"
	"net"
	"bytes"
	"io"
)

func TestWrite(testCtx *testing.T) {
	// given
	var (
		mockDestination = &mock.MockConn{
		Data: make([][]byte, 64*1024),
		Error: nil,
		LocalAddress:nil,
		RemoteAddress:nil,
	}
		source          = &net.TCPConn{}
		destination     = &net.TCPConn{}
		mockContext     = &chunkContext{
			description: "",
			data: make([]byte, 64*1024),
			from: source,
			to: destination,
			err: nil,
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
			requestNumber: 0,
	}
	)
	mockContext.data = []byte("this is the data that is going to be written")

	// when
	write(mockDestination)(mockContext)
	value := make([]byte, len(mockContext.data))
	copy(value, mockDestination.Data[0])

	// then
	writtenSize, _ := mockDestination.Write(mockContext.data)
	fmt.Printf("The write data size is %d and the actual written data is: [%s]\n", mockContext.totalWriteSize, value)
	if mockContext.totalWriteSize != 44 {
		testCtx.Fatalf("expected: [44] actual: [%d]", writtenSize)
	}
	if !bytes.Equal(mockContext.data, value) {
		testCtx.Fatalf("expected: [%s] actual: [%s]", mockContext.data, value)
	}
}


func TestWriteWithError(testCtx *testing.T) {
	// given
	var (
		mockDestination = &mock.MockConn{
		Data: make([][]byte, 64*1024),
		Error: io.ErrShortWrite,
		LocalAddress:nil,
		RemoteAddress:nil,
	}
		source          = &net.TCPConn{}
		destination     = &net.TCPConn{}
		mockContext     = &chunkContext{
		description: "",
		data: make([]byte, 64*1024),
		from: source,
		to: destination,
		err: nil,
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
		requestNumber: 0,
	}
	)
	mockContext.data = []byte("this is the data that is going to be written")

	// when
	write(mockDestination)(mockContext)
	value := make([]byte, len(mockContext.data))
	copy(value, mockDestination.Data[0])


	// then
	fmt.Printf("The expected error is: [%v] actual error is: [%v]\n", mockDestination.Error, mockContext.err)
	if mockContext.err != io.ErrShortWrite {
		testCtx.Fatalf("expected: [%v] actual: [%v]", io.ErrShortWrite, mockContext.err)
	}
}
