package stages

import (
	"io"
	"testing"
	mock "util/test/mock"
	assertion "util/test/assertion"
	"net"
	"proxy/contexts"
)

func NewTestWriteChunkContext(data string, err error) (*contexts.ChunkContext, *mock.MockConn) {
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
	mockDestination := mock.NewMockConn(err, 5)
	mockContext.Data = []byte(data)
	mockContext.To = mockDestination
	return mockContext, mockDestination
}

// test chunk no errors
// 	- should
// 		1. call write with context.Data once
// 		2. update context.TotalWriteSize
func Test_Write_With_Chunk_No_Error(testCtx *testing.T) {
	// given
	initialTotalWriteSize := int64(10)
	mockContext, mockDestination := NewTestWriteChunkContext("this is the data that is going to be written", nil)
	mockContext.TotalWriteSize = initialTotalWriteSize
	var expectedError error = nil

	// when
	write(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, int64(len(mockContext.Data))+initialTotalWriteSize, mockContext.TotalWriteSize)
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, expectedError, mockContext.Err)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 1, mockDestination.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, mockContext.Data, mockDestination.Data[0])
}

// test chunk none nil error
// 	- should
// 		1. call write with context.Data once
// 		2. update context.TotalWriteSize
// 		2. set context.Err
func Test_Write_With_Chunk_With_Error(testCtx *testing.T) {
	// given
	mockContext, mockDestination := NewTestWriteChunkContext("this is the data that is going to be written", io.EOF)

	// when
	write(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, int64(len(mockContext.Data)), mockContext.TotalWriteSize)
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, io.EOF, mockContext.Err)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 1, mockDestination.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, mockContext.Data, mockDestination.Data[0])
}

// test zero sized chunk
// 	- should
// 		1. not call write
func Test_Write_With_Zero_Chunk(testCtx *testing.T) {
	// given
	mockContext, mockDestination := NewTestWriteChunkContext("", nil)
	var expectedError error = nil

	// when
	write(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, int64(len(mockContext.Data)), mockContext.TotalWriteSize)
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, expectedError, mockContext.Err)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 0, mockDestination.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, string(mockContext.Data), string(mockDestination.Data[0]))
}

// test chunk amount written less than amountToWrite
// 	- should
// 		1. call write with context.Data once
// 		2. update context.TotalWriteSize
// 		2. set context.Err as io.ErrShortWrite
func Test_Write_With_Short_Write_Error(testCtx *testing.T) {
	// given
	mockContext, mockDestination := NewTestWriteChunkContext("this is the data that is going to be written", io.EOF)
	mockDestination.ShortWrite = true

	expectedData := make([]byte, len(mockContext.Data)/2)
	copy(expectedData, mockContext.Data)

	// when
	write(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, int64(len(mockContext.Data) / 2), mockContext.TotalWriteSize)
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, io.ErrShortWrite, mockContext.Err)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 1, mockDestination.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, expectedData, mockDestination.Data[0])
}

