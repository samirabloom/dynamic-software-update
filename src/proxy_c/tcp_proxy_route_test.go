package proxy_c

import (
	"testing"
	"net"
	"code.google.com/p/go-uuid/uuid"
	"util/test/assertion"
	"fmt"
)

// test firstChunk and clientToServer
// 	- should
// 		1. read dynsofyup cookie
// 		2. create backpipe
// 		3. call next
func Test_On_Route_With_First_Chunk(testCtx *testing.T) {
	// given
	var (
		mockContext = &chunkContext{
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
			requestUUID: uuid.NIL,
	}
		mockWriteCallCounter   = 0
		mockWriteChunkContexts = make([]*chunkContext, 5)
		mockWrite  = func(mockDestination *chunkContext) {
			mockWriteChunkContexts[mockWriteCallCounter] = mockDestination
			mockWriteCallCounter++
		}
		clientToServer bool = true

		createMockPipeCounter  = 0
		createMockPipe = func(context *chunkContext, clientToServer bool){
			createMockPipeCounter++
		// **************
		// TO DO AND MOVE THE ASSERT FUNCTION
		// ****************
			assertion.AssertDeepEqual("Correct new pipe Counter", testCtx, createMockPipeCounter, 1)
		}
	)
	mockContext.data = []byte("Cookie: dynsoftup=3f872698-0852-11e4-b87a-600308a8245e;")

	// when
	route(mockWrite, clientToServer, func() string {return uuid.NewUUID().String()},createMockPipe)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct uuid picked from cookie", testCtx, mockContext.requestUUID, uuid.Parse("3f872698-0852-11e4-b87a-600308a8245e"))
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, mockWriteCallCounter, 1)

}

// test firstChunk and not clientToServer and no requestUUID
// 	- should
// 		1. add cookie with new UUID value
// 		2. call next
func Test_On_Route_With_Not_ClientToServer_No_RequestUUID(testCtx *testing.T) {
	// given
	var (
		mockContext = &chunkContext{
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
		requestUUID: nil,
	}
		mockWriteCallCounter   = 0
		mockWriteChunkContexts = make([]*chunkContext, 5)
		mockWrite  = func(mockDestination *chunkContext) {
		mockWriteChunkContexts[mockWriteCallCounter] = mockDestination
		mockWriteCallCounter++
	}
		clientToServer bool = false

		createMockPipeCounter  = 0
		createMockPipe = func(context *chunkContext, clientToServer bool){
		createMockPipeCounter++
		}
		uuid string = "3f872698-0852-11e4-b87a-600308a8245e"
	)
	mockContext.data = []byte("this is a request with no cookie \n added")
	// when
	route(mockWrite, clientToServer, func() string {return uuid}, createMockPipe)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct context.data after adding cookie", testCtx, mockContext.data, []byte("this is a request with no cookie \nSet-Cookie: dynsoftup="+uuid+";\n added"))
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, mockWriteCallCounter, 1)
}

// test firstChunk and not clientToServer and context.requestUUID
// 	- should
// 		1. add cookie with context.requestUUID
// 		2. call next
func Test_On_Route_With_RequestUUID_Not_ClientToServer(testCtx *testing.T) {
}


// test not firstChunk
// 	- should
// 		1. do not create backpipe
// 		2. do not add cookie
// 		3. call next
func Test_On_Route_With_Not_First_Chunk(testCtx *testing.T) {
}



