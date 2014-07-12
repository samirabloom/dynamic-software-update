package proxy_c

import (
	"testing"
	"code.google.com/p/go-uuid/uuid"
	"util/test/assertion"
	"net"
)

func NewTestRouteChunkContext(data string, requestUUID uuid.UUID, clientToServer bool) *chunkContext {
	mockContext := NewTestChunkContext()
	mockContext.requestUUID = requestUUID
	mockContext.data = []byte(data)
	mockContext.clientToServer = clientToServer
	return mockContext
}

func panicingUUIDGenerator() uuid.UUID {
	panic("THIS FUNCTION SHOULD NOT BE CALLED")
	return nil
}

// test firstChunk and clientToServer
// 	- should
// 		1. read dynsofyup cookie
// 		2. create backpipe
// 		3. call next
func Test_On_Route_For_Request_With_First_Chunk(testCtx *testing.T) {
	// given
	var (
		uuidValue      = uuid.NewUUID()
		mockContext    = NewTestRouteChunkContext("Cookie: dynsoftup="+uuidValue.String()+";", uuid.NIL, true)
		mockWrite      = NewMockStage("mockWrite")
		mockCreatePipe = NewMockStage("mockCreatePipe")
		routingContext = &RoutingContext{backendBaseAddr: &net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}, loadBalanceCount: 3}
	)
	mockCreatePipe.close(1)

	// when
	route(mockWrite.mockStage, panicingUUIDGenerator, routingContext, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct UUID From Request Cookie", testCtx, mockContext.requestUUID, uuidValue)
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
		uuidValue            = uuid.NewUUID()
		mockContext          = NewTestRouteChunkContext("this is a request with no cookie \n added", nil, false)
		mockWrite            = NewMockStage("mockWrite")
		mockCreatePipe       = NewMockStage("mockCreatePipe")
		initialTotalReadSize = int64(10)
		expectedCookieHeader = "Set-Cookie: dynsoftup=" + uuidValue.String() + ";\n"
		uuidGenerator        = func() uuid.UUID {return uuidValue}
	)
	mockContext.totalReadSize = initialTotalReadSize

	// when
	route(mockWrite.mockStage, uuidGenerator, &RoutingContext{}, mockCreatePipe.mockStage)(mockContext)

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
		uuidValue            = uuid.NewUUID()
		mockContext          = NewTestRouteChunkContext("this is a request with no cookie \n added\n", uuidValue, false)
		mockWrite            = NewMockStage("mockWrite")
		mockCreatePipe       = NewMockStage("mockCreatePipe")
		initialTotalReadSize = int64(10)
		expectedCookieHeader = "Set-Cookie: dynsoftup=" + uuidValue.String() + ";\n"
	)
	mockContext.totalReadSize = initialTotalReadSize

	// when
	route(mockWrite.mockStage, panicingUUIDGenerator, &RoutingContext{}, mockCreatePipe.mockStage)(mockContext)

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
	)
	mockContext.firstChunk = false
	mockCreatePipe.close(1)

	// when
	route(mockWrite.mockStage, panicingUUIDGenerator, &RoutingContext{}, mockCreatePipe.mockStage)(mockContext)

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
	)
	mockContext.firstChunk = false

	// when
	route(mockWrite.mockStage, panicingUUIDGenerator, &RoutingContext{}, mockCreatePipe.mockStage)(mockContext)

	// then
	println("mockContext.data", "["+string(mockContext.data)+"]")
	println("[]byte(\"this is a response with no cookie \n added\")", "["+"this is a response with no cookie \n added"+"]")
	assertion.AssertDeepEqual("Correct Chunk Without Cookie", testCtx, mockContext.data, []byte("this is a response with no cookie \n added"))
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, mockWrite.mockStageCallCounter, 1)
}


