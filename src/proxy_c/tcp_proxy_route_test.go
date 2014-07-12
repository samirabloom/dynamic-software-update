package proxy_c

import (
	"testing"
	"code.google.com/p/go-uuid/uuid"
	"util/test/assertion"
)

func NewTestRouteChunkContext(data string, requestUUID uuid.UUID, clientToServer bool) *chunkContext {
	mockContext := NewTestChunkContext()
	mockContext.requestUUID = requestUUID
	mockContext.data = []byte(data)
	mockContext.clientToServer = clientToServer
	return mockContext
}

// test firstChunk and clientToServer
// 	- should
// 		1. read dynsofyup cookie
// 		2. create backpipe
// 		3. call next
func Test_On_Route_For_Request_With_First_Chunk(testCtx *testing.T) {
	// given
	var (
		uuidString     = "3f872698-0852-11e4-b87a-600308a8245e"
		mockContext    = NewTestRouteChunkContext("Cookie: dynsoftup="+uuidString+";", uuid.NIL, true)
		mockWrite      = NewMockStage("mockWrite")
		mockCreatePipe = NewMockStage("mockCreatePipe")
		uuidGenerator  = func() uuid.UUID {return nil}
	)
	mockCreatePipe.close(1)

	// when
	route(mockWrite.mockStage, uuidGenerator, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct UUID From Request Cookie", testCtx, mockContext.requestUUID, uuid.Parse(uuidString))
	<-mockCreatePipe.mockStageCallChannel
	assertion.AssertDeepEqual("Correct New Pipe Created", testCtx, mockCreatePipe.mockStageCallCounter, 1)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, mockWrite.mockStageCallCounter, 1)

}

// test firstChunk and not clientToServer and no requestUUID
// 	- should
// 		1. add cookie with new UUID value
// 		2. call next
func Test_On_Route_For_Response_With_No_RequestUUID(testCtx *testing.T) {
	// given
	var (
		mockContext          = NewTestRouteChunkContext("this is a request with no cookie \n added", nil, false)
		mockWrite            = NewMockStage("mockWrite")
		mockCreatePipe       = NewMockStage("mockCreatePipe")
		uuidString           = "3f872698-0852-11e4-b87a-600308a8245e"
		initialTotalReadSize = int64(10)
		expectedCookieHeader = "Set-Cookie: dynsoftup=" + uuidString + ";\n"
		uuidGenerator        = func() uuid.UUID {return uuid.Parse(uuidString)}
	)
	mockContext.totalReadSize = initialTotalReadSize

	// when
	route(mockWrite.mockStage, uuidGenerator, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk With Cookie", testCtx, mockContext.data, []byte("this is a request with no cookie \n"+expectedCookieHeader+" added"))
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, mockWrite.mockStageCallCounter, 1)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, mockContext.totalReadSize, int64(len(expectedCookieHeader))+initialTotalReadSize)
}

// test firstChunk and not clientToServer and context.requestUUID
// 	- should
// 		1. add cookie with context.requestUUID
// 		2. call next
func Test_On_Route_For_Response_With_RequestUUID(testCtx *testing.T) {
	// given
	var (
		uuidString           = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		mockContext          = NewTestRouteChunkContext("this is a request with no cookie \n added\n", uuid.Parse(uuidString), false)
		mockWrite            = NewMockStage("mockWrite")
		mockCreatePipe       = NewMockStage("mockCreatePipe")
		initialTotalReadSize = int64(10)
		expectedCookieHeader = "Set-Cookie: dynsoftup=" + uuidString + ";\n"
		uuidGenerator        = func() uuid.UUID {return nil}
	)
	mockContext.totalReadSize = initialTotalReadSize

	// when
	route(mockWrite.mockStage, uuidGenerator, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk With Cookie", testCtx, mockContext.data, []byte("this is a request with no cookie \n"+expectedCookieHeader+" added\n"))
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, mockWrite.mockStageCallCounter, 1)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, mockContext.totalReadSize, int64(len(expectedCookieHeader))+initialTotalReadSize)

}

// test not firstChunk and is clientToServer
// 	- should
// 		1. do not create backpipe
// 		3. call next
func Test_On_Route_For_Request_With_Not_First_Chunk(testCtx *testing.T) {
	// given
	var (
		mockContext    = NewTestRouteChunkContext("this is a request with no cookie \n added", nil, true)
		mockWrite      = NewMockStage("mockWrite")
		mockCreatePipe = NewMockStage("mockCreatePipe")
		uuidGenerator  = func() uuid.UUID {return nil}
	)
	mockContext.firstChunk = false
	mockCreatePipe.close(1)

	// when
	route(mockWrite.mockStage, uuidGenerator, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk Without Cookie", testCtx, mockContext.data, []byte("this is a request with no cookie \n added"))
	<-mockCreatePipe.mockStageCallChannel
	assertion.AssertDeepEqual("Correct New Pipe Created", testCtx, mockCreatePipe.mockStageCallCounter, 0)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, mockWrite.mockStageCallCounter, 1)
}

// test not firstChunk and not clientToServer
// 	- should
// 		1. do not add cookie
// 		2. call next
func Test_On_Route_For_Response_With_Not_First_Chunk(testCtx *testing.T) {
	// given
	var (
		mockContext    = NewTestRouteChunkContext("this is a response with no cookie \n added", nil, false)
		mockWrite      = NewMockStage("mockWrite")
		mockCreatePipe = NewMockStage("mockCreatePipe")
		uuidGenerator  = func() uuid.UUID {return nil}
	)
	mockContext.firstChunk = false

	// when
	route(mockWrite.mockStage, uuidGenerator, mockCreatePipe.mockStage)(mockContext)

	// then
	println("mockContext.data", "["+string(mockContext.data)+"]")
	println("[]byte(\"this is a response with no cookie \n added\")", "["+"this is a response with no cookie \n added"+"]")
	assertion.AssertDeepEqual("Correct Chunk Without Cookie", testCtx, mockContext.data, []byte("this is a response with no cookie \n added"))
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, mockWrite.mockStageCallCounter, 1)
}


