package proxy_c

import (
	"testing"
	"io"
	mock "util/test/mock"
	assertion "util/test/assertion"
)

func NewTestReadChunkContext(chunks []string, err error) (*chunkContext, *mock.MockConn) {
	mockContext := NewTestChunkContext()
	mockSource := mock.NewMockConn(err, len(chunks))
	mockContext.from = mockSource
	for index, chunk := range chunks {
		mockSource.Data[index] = []byte(chunk)
	}
	return mockContext, mockSource
}

// test readSize == 0 then EOF
// 	- should
// 		1. call src.Read once
// 		2. should not call next(context)
// 		3. should not set context.firstChunk to false
func Test_Read_With_No_Chunk_And_EOF_Error(testCtx *testing.T) {
	// given
	var (
		mockContext, mockSource = NewTestReadChunkContext([]string{}, io.EOF)
		mockRoute               = NewMockStage("mockRoute")
		mockComplete            = NewMockStage("mockComplete")
	)

	// when
	read(mockRoute.mockStage, mockComplete.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, mockSource.NumberOfReads, 1)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockSource.NumberOfWrites, 0)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, mockRoute.mockStageCallCounter, 0)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, mockComplete.mockStageCallCounter, 1)
	assertion.AssertDeepEqual("Correct Error", testCtx, mockSource.Error, io.EOF)
}

// test readSize == 0 then other error
// 	- should
// 		1. call src.Read once
// 		2. should not call next(context)
// 		3. should not set context.firstChunk to false
func Test_Read_With_Empty_Chunk_And_Non_EOF_Error(testCtx *testing.T) {
	// given
	var (
		mockContext, mockSource = NewTestReadChunkContext([]string{""}, io.ErrClosedPipe)
		mockRoute               = NewMockStage("mockRoute")
		mockComplete            = NewMockStage("mockComplete")
	)

	// when
	read(mockRoute.mockStage, mockComplete.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, mockSource.NumberOfReads, 2)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockSource.NumberOfWrites, 0)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, mockRoute.mockStageCallCounter, 0)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, mockComplete.mockStageCallCounter, 1)
	assertion.AssertDeepEqual("Correct Error", testCtx, mockSource.Error, io.ErrClosedPipe)
}

// test two chunks then EOF
// 	- should
// 		1. call src.Read twice
// 		2. set context.data to result from src.Read twice
// 		3. should update context.totalReadSize with amount read in both
// 		4. should call next(context) twice
// 		5. should set context.firstChunk to false
func Test_Read_With_Two_Chunks(testCtx *testing.T) {
	// given
	var (
		mockContext, mockSource = NewTestReadChunkContext([]string{"this is the first chunk", "this is the second chunk"}, io.EOF)
		mockRoute               = NewMockStage("mockRoute")
		mockComplete            = NewMockStage("mockComplete")
	)

	// when
	read(mockRoute.mockStage, mockComplete.mockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, mockSource.NumberOfReads, 3)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, mockSource.NumberOfWrites, 0)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, mockRoute.mockStageCallCounter, 2)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, mockComplete.mockStageCallCounter, 1)
	assertion.AssertDeepEqual("Correct Error", testCtx, mockSource.Error, io.EOF)
	assertion.AssertDeepEqual("Correct First Chunk", testCtx, mockRoute.mockStageChunkContexts[0].data, mockSource.Data[0])
	assertion.AssertDeepEqual("Correct First Chunk - firstChunk Indicator", testCtx, mockRoute.mockStageChunkContexts[0].firstChunk, true)
	assertion.AssertDeepEqual("Correct Second Chunk", testCtx, mockRoute.mockStageChunkContexts[1].data, mockSource.Data[1])
	assertion.AssertDeepEqual("Correct Second Chunk - firstChunk Indicator", testCtx, mockRoute.mockStageChunkContexts[1].firstChunk, false)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, mockContext.totalReadSize, int64(len(mockSource.Data[0]) + len(mockSource.Data[1])))
}
