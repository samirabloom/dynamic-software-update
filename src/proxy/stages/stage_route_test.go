package stages

import (
	"net"
	"testing"
	"code.google.com/p/go-uuid/uuid"
	assertion "util/test/assertion"
	"time"
	"proxy/contexts"
)

func NewTestRouteChunkContext(data string, clientToServer contexts.Direction) *contexts.ChunkContext {
	mockContext := &contexts.ChunkContext{
		Data: make([]byte, 64*1024),
		From: &net.TCPConn{},
		To: &net.TCPConn{},
		Err: nil,
		TotalReadSize: 0,
		TotalWriteSize: 0,
		PipeComplete: make(chan int64, 100),
		FirstChunk: true,
	}
	mockContext.Data = []byte(data)
	mockContext.Direction = clientToServer
	return mockContext
}

// test firstChunk and clientToServer
// 	- should
// 		1. read dynsofyup cookie
// 		2. create backpipe
// 		3. call next
func XTest_Route_For_Request_With_First_Chunk(testCtx *testing.T) {
	// given
	listener, err := net.Listen("tcp", ":1024")
	if err == nil {
		defer listener.Close()
	}
	var (
		mockWrite      = contexts.NewMockStage("mockWrite")
		mockCreatePipe = contexts.NewMockStage("mockCreatePipe")
		cluster        = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: &net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}}, RequestCounter: -1, Uuid: uuid.NewUUID(), Mode: contexts.InstantMode}
		clusters       = &contexts.Clusters{}
		mockContext    = NewTestRouteChunkContext("Cookie: dynsoftup="+cluster.Uuid.String()+";", true)
	)
	clusters.Add(cluster)
	mockCreatePipe.Close(1)

	// when
	route(mockWrite.MockStage, clusters, mockCreatePipe.MockStage)(mockContext)

	// then
	<-mockCreatePipe.MockStageCallChannel
	assertion.AssertDeepEqual("Correct New Pipe Created", testCtx, 1, mockCreatePipe.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.MockStageCallCounter)

}

// test firstChunk and not clientToServer and no requestUUID
// 	- should
// 		1. add cookie with new UUID value
// 		2. call next
func XTest_Route_For_Response_With_No_RequestUUID(testCtx *testing.T) {
	// given
	var (
		mockWrite             = contexts.NewMockStage("mockWrite")
		mockCreatePipe        = contexts.NewMockStage("mockCreatePipe")
		initialTotalReadSize  = int64(10)
		cluster               = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: &net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}}, RequestCounter: -1, Uuid: uuid.NewUUID(), Mode: contexts.SessionMode}
		clusters              = &contexts.Clusters{}
		expectedContentLength = "Content-Length: 40\n"
		expectedCookieHeader  = "Set-Cookie: dynsoftup=" + cluster.Uuid.String() + "; Expires=" + time.Now().Add(time.Second * time.Duration(0)).Format(time.RFC1123) + ";\n"
		mockContext           = NewTestRouteChunkContext("this is a request with no cookie \n added", false)
	)
	clusters.Add(cluster)

	mockContext.TotalReadSize = initialTotalReadSize

	// when
	route(mockWrite.MockStage, clusters, mockCreatePipe.MockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk With Cookie", testCtx, []byte("this is a request with no cookie \n"+expectedContentLength+expectedCookieHeader+" added"), mockContext.Data)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, int64(len(expectedContentLength) + len(expectedCookieHeader))+initialTotalReadSize, mockContext.TotalReadSize)
}

// test firstChunk and not clientToServer and context.requestUUID
// 	- should
// 		1. add cookie with context.requestUUID
// 		2. call next
func XTest_Route_For_Response_With_RequestUUID(testCtx *testing.T) {
	// given
	var (
		mockWrite            = contexts.NewMockStage("mockWrite")
		mockCreatePipe       = contexts.NewMockStage("mockCreatePipe")
		initialTotalReadSize = int64(10)
		cluster              = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: &net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}}, RequestCounter: -1, Uuid: uuid.NewUUID(), Mode: contexts.SessionMode}
		clusters             = &contexts.Clusters{}
		expectedCookieHeader = "Set-Cookie: dynsoftup=" + cluster.Uuid.String() + "; Expires=" + time.Now().Add(time.Second * time.Duration(0)).Format(time.RFC1123) + ";\n"
		mockContext          = NewTestRouteChunkContext("this is a request with no cookie \n added\n", false)
	)
	clusters.Add(cluster)
	mockContext.TotalReadSize = initialTotalReadSize

	// when
	route(mockWrite.MockStage, clusters, mockCreatePipe.MockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk With Cookie", testCtx, []byte("this is a request with no cookie \n"+expectedCookieHeader+" added\n"), mockContext.Data)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, +int64(len(expectedCookieHeader))+initialTotalReadSize, mockContext.TotalReadSize)

}

// test not firstChunk and is clientToServer
// 	- should
// 		1. do not create backpipe
// 		3. call next
func XTest_Route_For_Request_With_Not_First_Chunk(testCtx *testing.T) {
	// given
	var (
		mockContext    = NewTestRouteChunkContext("this is a request with no cookie \n added", true)
		mockWrite      = contexts.NewMockStage("mockWrite")
		mockCreatePipe = contexts.NewMockStage("mockCreatePipe")
		cluster        = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: &net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}}, RequestCounter: -1, Uuid: uuid.NewUUID()}
		clusters       = &contexts.Clusters{}
	)
	clusters.Add(cluster)
	mockContext.FirstChunk = false
	mockCreatePipe.Close(1)

	// when
	route(mockWrite.MockStage, clusters, mockCreatePipe.MockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk Without Cookie", testCtx, []byte("this is a request with no cookie \n added"), mockContext.Data)
	<-mockCreatePipe.MockStageCallChannel
	assertion.AssertDeepEqual("Correct New Pipe Created", testCtx, 0, mockCreatePipe.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.MockStageCallCounter)
}

// test not firstChunk and not clientToServer
// 	- should
// 		1. do not add cookie
// 		2. call next
func XTest_Route_For_Response_With_Not_First_Chunk(testCtx *testing.T) {
	// given
	var (
		mockContext    = NewTestRouteChunkContext("this is a response with no cookie \n added", false)
		mockWrite      = contexts.NewMockStage("mockWrite")
		mockCreatePipe = contexts.NewMockStage("mockCreatePipe")
		cluster        = &contexts.Cluster{BackendAddresses: []*contexts.BackendAddress{&contexts.BackendAddress{Address: &net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}}, RequestCounter: -1, Uuid: uuid.NewUUID()}
		clusters       = &contexts.Clusters{}
	)
	clusters.Add(cluster)
	mockContext.FirstChunk = false

	// when
	route(mockWrite.MockStage, clusters, mockCreatePipe.MockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk Without Cookie", testCtx, []byte("this is a response with no cookie \n added"), mockContext.Data)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.MockStageCallCounter)
}


