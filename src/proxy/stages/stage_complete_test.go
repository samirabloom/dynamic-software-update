package stages

import (
	"testing"
	"net"
	"syscall"
	mock "util/test/mock"
	assertion "util/test/assertion"
	"proxy/contexts"
)

func NewTestCompleteChunkContext(chunks []string, err error) (*contexts.ChunkContext, *mock.MockConn, *mock.MockConn) {
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
	mockContext.Err = err
	mockContext.From = mock.NewMockConn(err, len(chunks))
	mockContext.To = mock.NewMockConn(err, len(chunks))
	return mockContext, mockContext.To.(*mock.MockConn), mockContext.From.(*mock.MockConn)
}

// test no error
// 	- should
// 		1. read closed

// IGNORED AS NOT POSSIBLE TO TEST CURRENTLY WITH MOCK CONNECTION
func XTest_Complete_With_No_Error(testCtx *testing.T) {
	// given
	var (
		mockContext, mockSource, mockDestination = NewTestCompleteChunkContext([]string{}, nil)
	)

	// when
	complete(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Read Closed", testCtx, mockSource.ReadClosed, true)
	assertion.AssertDeepEqual("Correct Write Closed", testCtx, mockDestination.WriteClosed, false)
}

// test syscall.EPIPE error
// 	- should
// 		1. write closed
// 		2. read closed

// IGNORED AS NOT POSSIBLE TO TEST CURRENTLY WITH MOCK CONNECTION
func XTest_Complete_With_EPIPE_Error(testCtx *testing.T) {
	// given
	var (
		mockContext, mockSource, mockDestination = NewTestCompleteChunkContext([]string{}, &net.OpError{Err: syscall.EPIPE})
	)

	// when
	complete(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Read Closed", testCtx, mockSource.ReadClosed, true)
	assertion.AssertDeepEqual("Correct Wrtie Closed", testCtx, mockDestination.WriteClosed, true)
}
