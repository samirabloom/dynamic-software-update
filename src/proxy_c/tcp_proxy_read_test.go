package proxy_c

import (
	"testing"
	"util/test/mock"
	"net"
	"io"
	"util/test/assertion"
)

// test readSize == 0 then EOF
// 	- should
// 		1. call src.Read once
// 		2. should not call next(context)
// 		3. should not set context.firstChunk to false
func Test_On_Read_With_No_Chunk_And_EOF_Error(testCtx *testing.T) {
	// given
	var (
		mockSource                = &mock.MockConn{
			Data: make([][]byte, 0),
			Error: io.EOF,
			LocalAddress:nil,
			RemoteAddress:nil,
		}
		mockContext               = &chunkContext{
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
		mockRouteCallCounter      = 0
		mockRouteChunkContexts    = make([]*chunkContext, 5)
		mockRoute                 = func(context *chunkContext) {
			mockRouteChunkContexts[mockRouteCallCounter] = context
			mockRouteCallCounter++
		}
		mockCompleteCallCounter   = 0
		mockCompleteChunkContexts = make([]*chunkContext, 5)
		mockComplete              = func(context *chunkContext) {
			mockCompleteChunkContexts[mockCompleteCallCounter] = context
			mockCompleteCallCounter++
		}
	)

	// when
	read(mockRoute, mockSource, mockComplete)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, mockSource.NumberOfReads, 1)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockSource.NumberOfWrites, 0)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, mockRouteCallCounter, 0)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, mockCompleteCallCounter, 1)
	assertion.AssertDeepEqual("Correct Error", testCtx, mockSource.Error, io.EOF)
}

// test readSize == 0 then other error
// 	- should
// 		1. call src.Read once
// 		2. should not call next(context)
// 		3. should not set context.firstChunk to false
func Test_On_Read_With_Empty_Chunk_And_Non_EOF_Error(testCtx *testing.T) {
	// given
	var (
		mockSource                = &mock.MockConn{
			Data: make([][]byte, 1),
			Error: io.ErrClosedPipe,
			LocalAddress:nil,
			RemoteAddress:nil,
		}
		mockContext               = &chunkContext{
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
		mockRouteCallCounter      = 0
		mockRouteChunkContexts    = make([]*chunkContext, 5)
		mockRoute                 = func(context *chunkContext) {
			mockRouteChunkContexts[mockRouteCallCounter] = context
			mockRouteCallCounter++
		}
		mockCompleteCallCounter   = 0
		mockCompleteChunkContexts = make([]*chunkContext, 5)
		mockComplete              = func(context *chunkContext) {
			mockCompleteChunkContexts[mockCompleteCallCounter] = context
			mockCompleteCallCounter++
		}
	)
	mockSource.Data[0] = []byte("")

	// when
	read(mockRoute, mockSource, mockComplete)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, mockSource.NumberOfReads, 2)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockSource.NumberOfWrites, 0)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, mockRouteCallCounter, 0)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, mockCompleteCallCounter, 1)
	assertion.AssertDeepEqual("Correct Error", testCtx, mockSource.Error, io.ErrClosedPipe)
}

// test two chunks then EOF
// 	- should
// 		1. call src.Read twice
// 		2. set context.data to result from src.Read twice
// 		3. should update context.totalReadSize with amount read in both
// 		4. should call next(context) twice
// 		5. should set context.firstChunk to false
func Test_On_Read_With_Two_Chunks(testCtx *testing.T) {
	// given
	var (
		mockSource                = &mock.MockConn{
			Data: make([][]byte, 2),
			Error: io.EOF,
			LocalAddress:nil,
			RemoteAddress:nil,
		}
		mockContext               = &chunkContext{
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
		mockRouteCallCounter      = 0
		mockRouteChunkContexts    = make([]*chunkContext, 5)
		mockRoute                 = func(context *chunkContext) {
			mockRouteChunkContexts[mockRouteCallCounter] = CopyChunkContext(context)
			mockRouteCallCounter++
		}
		mockCompleteCallCounter   = 0
		mockCompleteChunkContexts = make([]*chunkContext, 5)
		mockComplete              = func(context *chunkContext) {
			mockCompleteChunkContexts[mockCompleteCallCounter] = context
			mockCompleteCallCounter++
		}
	)
	mockSource.Data[0] = []byte("this is the first chunk")
	mockSource.Data[1] = []byte("this is the second chunk")

	// when
	read(mockRoute, mockSource, mockComplete)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, mockSource.NumberOfReads, 3)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockSource.NumberOfWrites, 0)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, mockRouteCallCounter, 2)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, mockCompleteCallCounter, 1)
	assertion.AssertDeepEqual("Correct Error", testCtx, mockSource.Error, io.EOF)
	assertion.AssertDeepEqual("Correct First Chunk", testCtx, mockRouteChunkContexts[0].data, mockSource.Data[0])
	assertion.AssertDeepEqual("Correct First Chunk - firstChunk Indicator", testCtx, mockRouteChunkContexts[0].firstChunk, true)
	assertion.AssertDeepEqual("Correct Second Chunk", testCtx, mockRouteChunkContexts[1].data, mockSource.Data[1])
	assertion.AssertDeepEqual("Correct First Chunk - firstChunk Indicator", testCtx, mockRouteChunkContexts[0].firstChunk, false)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, mockContext.totalReadSize, int64(len(mockSource.Data[0]) + len(mockSource.Data[1])))
}
