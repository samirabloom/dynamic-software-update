package proxy_c

import (
	"io"
	"testing"
	mock "util/test/mock"
	assertion "util/test/assertion"
)

func NewTestWriteChunkContext(data string, err error) (*chunkContext, *mock.MockConn) {
	mockContext := NewTestChunkContext()
	mockDestination := mock.NewMockConn(err, 5)
	mockContext.data = []byte(data)
	mockContext.to = mockDestination
	return mockContext, mockDestination
}

// test chunk no errors
// 	- should
// 		1. call write with context.data once
// 		2. update context.totalWriteSize
func Test_Write_With_Chunk_No_Error(testCtx *testing.T) {
	// given
	initialTotalWriteSize := int64(10)
	mockContext, mockDestination := NewTestWriteChunkContext("this is the data that is going to be written", nil)
	mockContext.totalWriteSize = initialTotalWriteSize
	var expectedError error = nil

	// when
	write(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, int64(len(mockContext.data))+initialTotalWriteSize, mockContext.totalWriteSize)
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, expectedError, mockContext.err)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 1, mockDestination.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, mockContext.data, mockDestination.Data[0])
}

// test chunk none nil error
// 	- should
// 		1. call write with context.data once
// 		2. update context.totalWriteSize
// 		2. set context.err
func Test_Write_With_Chunk_With_Error(testCtx *testing.T) {
	// given
	mockContext, mockDestination := NewTestWriteChunkContext("this is the data that is going to be written", io.EOF)

	// when
	write(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, int64(len(mockContext.data)), mockContext.totalWriteSize)
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, io.EOF, mockContext.err)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 1, mockDestination.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, mockContext.data, mockDestination.Data[0])
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
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, int64(len(mockContext.data)), mockContext.totalWriteSize)
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, expectedError, mockContext.err)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 0, mockDestination.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, string(mockContext.data), string(mockDestination.Data[0]))
}

// test chunk amount written less than amountToWrite
// 	- should
// 		1. call write with context.data once
// 		2. update context.totalWriteSize
// 		2. set context.err as io.ErrShortWrite
func Test_Write_With_Short_Write_Error(testCtx *testing.T) {
	// given
	mockContext, mockDestination := NewTestWriteChunkContext("this is the data that is going to be written", io.EOF)
	mockDestination.ShortWrite = true

	expectedData := make([]byte, len(mockContext.data)/2)
	copy(expectedData, mockContext.data)

	// when
	write(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, int64(len(mockContext.data) / 2), mockContext.totalWriteSize)
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, io.ErrShortWrite, mockContext.err)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 1, mockDestination.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, expectedData, mockDestination.Data[0])
}

