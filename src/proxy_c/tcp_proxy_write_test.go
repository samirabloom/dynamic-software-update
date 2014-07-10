package proxy_c

import (
	"testing"
	"util/test/mock"
	"util/test/assertion"
	"fmt"
	"net"
	"io"
)

// test chunk no errors
// 	- should
// 		1. call write with context.data once
// 		2. update context.totalWriteSize

func Test_On_Write_With_Call_Write_And_No_Error(testCtx *testing.T) {
	// given
	var (
		mockDestination = &mock.MockConn{
		Data: make([][]byte, 64*1024),
		Error: nil,
		NumberOfWrites: 0,
		LocalAddress:nil,
		RemoteAddress:nil,
	}
		mockContext     = &chunkContext{
		description: "",
		data: make([]byte, 64*1024),
		from: &net.TCPConn{},
		to: &net.TCPConn{},
		err: nil,
		totalReadSize: 0,
		totalWriteSize: 10,
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

	// then
	assertion.AssertDeepEqual("Correct Number Of total WriteSize", testCtx, mockContext.totalWriteSize, int64(54))
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, mockContext.err, nil)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockDestination.NumberOfWrites, 1)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, mockContext.data, mockDestination.Data[0])
}

// test chunk none nil error
// 	- should
// 		1. call write with context.data once
// 		2. update context.totalWriteSize
// 		2. set context.err
func Test_On_Write_With_Call_Write_And_Set_Error(testCtx *testing.T) {
	// given
	var (
		mockDestination = &mock.MockConn{
		Data: make([][]byte, 64*1024),
		Error: io.EOF,
		NumberOfWrites: 0,
		LocalAddress:nil,
		RemoteAddress:nil,
	}
		mockContext     = &chunkContext{
		description: "",
		data: make([]byte, 64*1024),
		from: &net.TCPConn{},
		to: &net.TCPConn{},
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

	// then
	assertion.AssertDeepEqual("Correct Number Of total WriteSize", testCtx, mockContext.totalWriteSize, int64(44))
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, mockContext.err, io.EOF)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockDestination.NumberOfWrites, 1)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, mockContext.data, mockDestination.Data[0])
}

// test zero sized chunk
// 	- should
// 		1. not call write

func Test_On_Write_With_Zero_Chunk(testCtx *testing.T) {
	// given
	var (
		mockDestination = &mock.MockConn{
		Data: make([][]byte, 64*1024),
		Error: nil,
		NumberOfWrites: 0,
		LocalAddress:nil,
		RemoteAddress:nil,
	}
		mockContext     = &chunkContext{
		description: "",
		data: make([]byte, 0),
		from: &net.TCPConn{},
		to: &net.TCPConn{},
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
	fmt.Printf("data context length %d\n", len(mockContext.data))
	//	mockContext.data = []byte("this is the data that is going to be written")

	// when
	write(mockDestination)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of total WriteSize", testCtx, mockContext.totalWriteSize, int64(0))
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, mockContext.err, nil)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockDestination.NumberOfWrites, 0)
}

// test chunk amount written less than amountToWrite
// 	- should
// 		1. call write with context.data once
// 		2. update context.totalWriteSize
// 		2. set context.err as io.ErrShortWrite

func Test_On_Write_With_Call_Write_Short_Write_Error(testCtx *testing.T) {
	// given
	var (
		mockDestination = &mock.MockConn{
		Data: make([][]byte, 64*1024),
		Error: io.EOF,
		NumberOfWrites: 0,
		LocalAddress:nil,
		RemoteAddress:nil,
//		shortWrite :true,
	}
		mockContext     = &chunkContext{
		description: "",
		data: make([]byte, 64*1024),
		from: &net.TCPConn{},
		to: &net.TCPConn{},
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
	//	mockDestination.shortWrite = true
	mockContext.data = []byte("this is the data that is going to be written")

	// when
	write(mockDestination)(mockContext)

	// then
//	assertion.AssertDeepEqual("Correct Number Of total WriteSize", testCtx, mockContext.totalWriteSize, int64(22))
//	assertion.AssertDeepEqual("Correct Context Error ", testCtx, mockContext.err, io.ErrShortWrite)
//	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockDestination.NumberOfWrites, 1)
//	assertion.AssertDeepEqual("Correct Written Data", testCtx, mockContext.data, mockDestination.Data[0])
}

