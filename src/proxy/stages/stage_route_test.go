package stages

import (
	"net"
	"testing"
	"code.google.com/p/go-uuid/uuid"
	assertion "util/test/assertion"
	"time"
)

func NewTestRouteChunkContext(data string, clientToServer bool) *ChunkContext {
	mockContext := &ChunkContext{
		description: "",
		data: make([]byte, 64*1024),
		from: &net.TCPConn{},
		to: &net.TCPConn{},
		err: nil,
		totalReadSize: 0,
		totalWriteSize: 0,
		pipeComplete: make(chan int64, 100),
		firstChunk: true,
	}
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
		mockWrite      = NewMockStage("mockWrite")
		mockCreatePipe = NewMockStage("mockCreatePipe")
		cluster        = &Cluster{BackendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, RequestCounter: -1, Uuid: uuid.NewUUID()}
		clusters       = &Clusters{}
		mockContext    = NewTestRouteChunkContext("Cookie: dynsoftup="+cluster.Uuid.String()+";", true)
	)
	clusters.Add(cluster)
	mockCreatePipe.close(1)

	// when
	route(mockWrite.mockStage, clusters, mockCreatePipe.mockStage)(mockContext)

	// then
	//	assertion.AssertDeepEqual("Correct Cluster for UUID in Cookie", testCtx, cluster, mockContext.cluster)
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
		mockWrite             = NewMockStage("mockWrite")
		mockCreatePipe        = NewMockStage("mockCreatePipe")
		initialTotalReadSize  = int64(10)
		cluster               = &Cluster{BackendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, RequestCounter: -1, Uuid: uuid.NewUUID(), Mode: SessionMode}
		clusters              = &Clusters{}
		expectedContentLength = "Content-Length: 40\n"
		expectedCookieHeader  = "Set-Cookie: dynsoftup=" + cluster.Uuid.String() + "; Expires=" + time.Now().Add(time.Second * time.Duration(0)).Format(time.RFC1123) + ";\n"
		mockContext           = NewTestRouteChunkContext("this is a request with no cookie \n added", false)
	)
	clusters.Add(cluster)

	mockContext.totalReadSize = initialTotalReadSize

	// when
	route(mockWrite.mockStage, clusters, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk With Cookie", testCtx, []byte("this is a request with no cookie \n"+expectedContentLength+expectedCookieHeader+" added"), mockContext.data)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, int64(len(expectedContentLength) + len(expectedCookieHeader))+initialTotalReadSize, mockContext.totalReadSize)
}

// test firstChunk and not clientToServer and context.requestUUID
// 	- should
// 		1. add cookie with context.requestUUID
// 		2. call next
func Test_Route_For_Response_With_RequestUUID(testCtx *testing.T) {
	// given
	var (
		mockWrite            = NewMockStage("mockWrite")
		mockCreatePipe       = NewMockStage("mockCreatePipe")
		initialTotalReadSize = int64(10)
		cluster              = &Cluster{BackendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, RequestCounter: -1, Uuid: uuid.NewUUID(), Mode: SessionMode}
		clusters             = &Clusters{}
		expectedCookieHeader = "Set-Cookie: dynsoftup=" + cluster.Uuid.String() + "; Expires=" + time.Now().Add(time.Second * time.Duration(0)).Format(time.RFC1123) + ";\n"
		mockContext          = NewTestRouteChunkContext("this is a request with no cookie \n added\n", false)
	)
	clusters.Add(cluster)
	mockContext.totalReadSize = initialTotalReadSize

	// when
	route(mockWrite.mockStage, clusters, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk With Cookie", testCtx, []byte("this is a request with no cookie \n"+expectedCookieHeader+" added\n"), mockContext.data)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, +int64(len(expectedCookieHeader))+initialTotalReadSize, mockContext.totalReadSize)

}

// test not firstChunk and is clientToServer
// 	- should
// 		1. do not create backpipe
// 		3. call next
func Test_Route_For_Request_With_Not_First_Chunk(testCtx *testing.T) {
	// given
	var (
		mockContext    = NewTestRouteChunkContext("this is a request with no cookie \n added", true)
		mockWrite      = NewMockStage("mockWrite")
		mockCreatePipe = NewMockStage("mockCreatePipe")
		cluster        = &Cluster{BackendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, RequestCounter: -1, Uuid: uuid.NewUUID()}
		clusters       = &Clusters{}
	)
	clusters.Add(cluster)
	mockContext.firstChunk = false
	mockCreatePipe.close(1)

	// when
	route(mockWrite.mockStage, clusters, mockCreatePipe.mockStage)(mockContext)

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
		mockContext    = NewTestRouteChunkContext("this is a response with no cookie \n added", false)
		mockWrite      = NewMockStage("mockWrite")
		mockCreatePipe = NewMockStage("mockCreatePipe")
		cluster        = &Cluster{BackendAddresses: []*net.TCPAddr{&net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 1024}}, RequestCounter: -1, Uuid: uuid.NewUUID()}
		clusters       = &Clusters{}
	)
	clusters.Add(cluster)
	mockContext.firstChunk = false

	// when
	route(mockWrite.mockStage, clusters, mockCreatePipe.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Chunk Without Cookie", testCtx, []byte("this is a response with no cookie \n added"), mockContext.data)
	assertion.AssertDeepEqual("Correct Write Call Counter", testCtx, 1, mockWrite.mockStageCallCounter)
}

func Test_Parse_Header(testCtx *testing.T) {
	// given
	var (
		parsedHeader = &headerMetrics{}
		data         = []byte("HTTP/1.1 200 OK\n" +
			"Content-Length: 143\n" +
			"Connection: keep-alive\n" +
			"Expires: 5\n" +
			"Content-Type: text/plain; charset=utf-8\n" +
			"Transfer-Encoding: chunked\n")
	)
	parsedHeader.headers = make(map[string]string)

	// when
	parseMetrics(parsedHeader, data)

	// then expected/actual
	assertion.AssertDeepEqual("Correct Content-Length", testCtx, int64(143), parsedHeader.contentLength)
	assertion.AssertDeepEqual("Correct HTTP Status", testCtx, 200, parsedHeader.statusCode)
	assertion.AssertDeepEqual("Correct Expires", testCtx, "5", parsedHeader.headers["Expires"])
	assertion.AssertDeepEqual("Correct Transfer-Encoding", testCtx, "chunked", parsedHeader.headers["Transfer-Encoding"])
	assertion.AssertDeepEqual("Correct Connection", testCtx, "keep-alive", parsedHeader.headers["Connection"])
	assertion.AssertDeepEqual("Correct Content-Type", testCtx, "text/plain; charset=utf-8", parsedHeader.headers["Content-Type"])
}


