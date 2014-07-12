package proxy_c

import (
	"testing"
	"util/test/mock"
	"util/test/assertion"
	"io"
)

func NewTestWriteChunkContext(data string) *chunkContext {
	mockContext := NewTestChunkContext()
	mockContext.data = []byte(data)
	return mockContext
}

// test chunk no errors
// 	- should
// 		1. call write with context.data once
// 		2. update context.totalWriteSize
func Test_On_Write_With_Chunk_No_Error(testCtx *testing.T) {
	// given
	var (
		mockDestination             = mock.NewMockConn(nil, 5)
		initialTotalWriteSize int64 = 10
		mockContext                 = NewTestWriteChunkContext("this is the data that is going to be written")
	)
	mockContext.totalWriteSize = initialTotalWriteSize

	// when
	write(mockDestination)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, mockContext.totalWriteSize, int64(len(mockContext.data))+initialTotalWriteSize)
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, mockContext.err, nil)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockDestination.NumberOfWrites, 1)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, mockContext.data, mockDestination.Data[0])
}

// test chunk none nil error
// 	- should
// 		1. call write with context.data once
// 		2. update context.totalWriteSize
// 		2. set context.err
func Test_On_Write_With_Chunk_With_Error(testCtx *testing.T) {
	// given
	var (
		mockDestination = mock.NewMockConn(io.EOF, 5)
		mockContext     = NewTestWriteChunkContext("this is the data that is going to be written")
	)

	// when
	write(mockDestination)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, mockContext.totalWriteSize, int64(len(mockContext.data)))
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, mockContext.err, io.EOF)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockDestination.NumberOfWrites, 1)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, mockContext.data, mockDestination.Data[0])
}

// test zero sized chunk
// 	- should
// 		1. not call write
func Test_On_Write_With_Zero_Chunk(testCtx *testing.T) {
	// given
	var (
		mockDestination = mock.NewMockConn(nil, 5)
		mockContext     = NewTestWriteChunkContext("")
	)

	// when
	write(mockDestination)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, mockContext.totalWriteSize, int64(len(mockContext.data)))
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, mockContext.err, nil)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockDestination.NumberOfWrites, 0)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, string(mockContext.data), string(mockDestination.Data[0]))
}

// test chunk amount written less than amountToWrite
// 	- should
// 		1. call write with context.data once
// 		2. update context.totalWriteSize
// 		2. set context.err as io.ErrShortWrite
func Test_On_Write_With_Short_Write_Error(testCtx *testing.T) {
	// given
	var (
		mockDestination = mock.NewMockConn(io.EOF, 5)
		mockContext     = NewTestWriteChunkContext("this is the data that is going to be written")
	)
	mockDestination.ShortWrite = true
	expectedData := make([]byte, len(mockContext.data)/2)
	copy(expectedData, mockContext.data)

	// when
	write(mockDestination)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Total WriteSize", testCtx, mockContext.totalWriteSize, int64(len(mockContext.data) / 2))
	assertion.AssertDeepEqual("Correct Context Error ", testCtx, mockContext.err, io.ErrShortWrite)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockDestination.NumberOfWrites, 1)
	assertion.AssertDeepEqual("Correct Written Data", testCtx, expectedData, mockDestination.Data[0])
}

