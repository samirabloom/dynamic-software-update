package proxy_c

import (
	"testing"
	"net"
	"syscall"
	mock "util/test/mock"
	assertion "util/test/assertion"
)

func NewTestCompleteChunkContext(chunks []string, err error) (*chunkContext, *mock.MockConn, *mock.MockConn) {
	mockContext := NewTestChunkContext()
	mockContext.err = err
	mockContext.from = mock.NewMockConn(err, len(chunks))
	mockContext.to = mock.NewMockConn(err, len(chunks))
	return mockContext, mockContext.to.(*mock.MockConn), mockContext.from.(*mock.MockConn)
}

// test no error
// 	- should
// 		1. read closed
func Test_Complete_With_No_Error(testCtx *testing.T) {
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
func Test_Complete_With_EPIPE_Error(testCtx *testing.T) {
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
