package proxy_c

import (
	"net"
	"testing"
	"code.google.com/p/go-uuid/uuid"
	assertion "util/test/assertion"
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
func Test_Route_For_Request_With_First_Chunk(testCtx *testing.T) {
	// given
	listener, err := net.Listen("tcp", ":1024")
	if err == nil {
		defer listener.Close()
	}
	var (
		expectedUuid    = uuid.NewUUID()
		mockContext     = NewTestRouteChunkContext("Cookie: dynsoftup="+expectedUuid.String()+";", uuid.NIL, true)
		mockWrite       = NewMockStage("mockWrite")
		mockCreatePipe  = NewMockStage("mockCreatePipe")
		routingContext  = &RoutingContext{backendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, requestCounter: -1, uuid: uuid.NewUUID()}
		routingContexts = &RoutingContexts{}
	)
	routingContexts.Add(routingContext)
	mockCreatePipe.close(1)

	// when
	route(mockWrite.mockStage, routingContexts, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct UUID From Request Cookie", testCtx, expectedUuid, mockContext.requestUUID)
	<-mockCreatePipe.mockStageCallChannel
	assertion.AssertDeepEqual("Correct New Pipe Created", testCtx, 1, mockCreatePipe.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.mockStageCallCounter)

}

// test firstChunk and not clientToServer and no requestUUID
// 	- should
// 		1. add cookie with new UUID value
// 		2. call next
func Test_Route_For_Response_With_No_RequestUUID(testCtx *testing.T) {
	// given
	var (
		uuidValue            = uuid.NewUUID()
		mockContext          = NewTestRouteChunkContext("this is a request with no cookie \n added", uuid.NewUUID(), false)
		mockWrite            = NewMockStage("mockWrite")
		mockCreatePipe       = NewMockStage("mockCreatePipe")
		initialTotalReadSize = int64(10)
		expectedCookieHeader = "Set-Cookie: dynsoftup=" + uuidValue.String() + ";\n"
		routingContext       = &RoutingContext{backendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, requestCounter: -1, uuid: uuid.NewUUID()}
		routingContexts      = &RoutingContexts{}
	)
	routingContexts.Add(routingContext)

	mockContext.totalReadSize = initialTotalReadSize

	// when
	route(mockWrite.mockStage, routingContexts, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk With Cookie", testCtx, []byte("this is a request with no cookie \n"+expectedCookieHeader+" added"), mockContext.data)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, int64(len(expectedCookieHeader))+initialTotalReadSize, mockContext.totalReadSize)
}

// test firstChunk and not clientToServer and context.requestUUID
// 	- should
// 		1. add cookie with context.requestUUID
// 		2. call next
func Test_Route_For_Response_With_RequestUUID(testCtx *testing.T) {
	// given
	var (
		uuidValue            = uuid.NewUUID()
		mockContext          = NewTestRouteChunkContext("this is a request with no cookie \n added\n", uuidValue, false)
		mockWrite            = NewMockStage("mockWrite")
		mockCreatePipe       = NewMockStage("mockCreatePipe")
		initialTotalReadSize = int64(10)
		expectedCookieHeader = "Set-Cookie: dynsoftup=" + uuidValue.String() + ";\n"
		routingContext       = &RoutingContext{backendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, requestCounter: -1, uuid: uuid.NewUUID()}
		routingContexts      = &RoutingContexts{}
	)
	routingContexts.Add(routingContext)
	mockContext.totalReadSize = initialTotalReadSize

	// when
	route(mockWrite.mockStage, routingContexts, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk With Cookie", testCtx, []byte("this is a request with no cookie \n"+expectedCookieHeader+" added\n"), mockContext.data)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, int64(len(expectedCookieHeader))+initialTotalReadSize, mockContext.totalReadSize)

}

// test not firstChunk and is clientToServer
// 	- should
// 		1. do not create backpipe
// 		3. call next
func Test_Route_For_Request_With_Not_First_Chunk(testCtx *testing.T) {
	// given
	var (
		mockContext     = NewTestRouteChunkContext("this is a request with no cookie \n added", nil, true)
		mockWrite       = NewMockStage("mockWrite")
		mockCreatePipe  = NewMockStage("mockCreatePipe")
		routingContext  = &RoutingContext{backendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, requestCounter: -1, uuid: uuid.NewUUID()}
		routingContexts = &RoutingContexts{}
	)
	routingContexts.Add(routingContext)
	mockContext.firstChunk = false
	mockCreatePipe.close(1)

	// when
	route(mockWrite.mockStage, routingContexts, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk Without Cookie", testCtx, []byte("this is a request with no cookie \n added"), mockContext.data)
	<-mockCreatePipe.mockStageCallChannel
	assertion.AssertDeepEqual("Correct New Pipe Created", testCtx, 0, mockCreatePipe.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.mockStageCallCounter)
}

// test not firstChunk and not clientToServer
// 	- should
// 		1. do not add cookie
// 		2. call next
func Test_Route_For_Response_With_Not_First_Chunk(testCtx *testing.T) {
	// given
	var (
		mockContext     = NewTestRouteChunkContext("this is a response with no cookie \n added", nil, false)
		mockWrite       = NewMockStage("mockWrite")
		mockCreatePipe  = NewMockStage("mockCreatePipe")
		routingContext  = &RoutingContext{backendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, requestCounter: -1, uuid: uuid.NewUUID()}
		routingContexts = &RoutingContexts{}
	)
	routingContexts.Add(routingContext)
	mockContext.firstChunk = false

	// when
	route(mockWrite.mockStage, routingContexts, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk Without Cookie", testCtx, []byte("this is a response with no cookie \n added"), mockContext.data)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.mockStageCallCounter)
}


