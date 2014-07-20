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
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, 1, mockSource.NumberOfReads)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 0, mockSource.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, 0, mockRoute.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, 1, mockComplete.mockStageCallCounter)
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
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, 2, mockSource.NumberOfReads)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 0, mockSource.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, 0, mockRoute.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, 1, mockComplete.mockStageCallCounter)
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
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, 3, mockSource.NumberOfReads)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 0, mockSource.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, 2, mockRoute.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, 1, mockComplete.mockStageCallCounter)
	assertion.AssertDeepEqual("Correct Error", testCtx, io.EOF, mockSource.Error)
	assertion.AssertDeepEqual("Correct First Chunk", testCtx, mockSource.Data[0], mockRoute.mockStageChunkContexts[0].data)
	assertion.AssertDeepEqual("Correct First Chunk - firstChunk Indicator", testCtx, true, mockRoute.mockStageChunkContexts[0].firstChunk)
	assertion.AssertDeepEqual("Correct Second Chunk", testCtx, mockSource.Data[1], mockRoute.mockStageChunkContexts[1].data)
	assertion.AssertDeepEqual("Correct Second Chunk - firstChunk Indicator", testCtx, false, mockRoute.mockStageChunkContexts[1].firstChunk)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, int64(len(mockSource.Data[0]) + len(mockSource.Data[1])), mockContext.totalReadSize)
}
