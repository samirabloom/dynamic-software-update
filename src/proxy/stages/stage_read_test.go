package stages

import (
	"testing"
	"io"
	mock "util/test/mock"
	assertion "util/test/assertion"
	"net"
	"proxy/contexts"
)

func NewTestReadChunkContext(chunks []string, err error) (*contexts.ChunkContext, *mock.MockConn) {
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
	mockSource := mock.NewMockConn(err, len(chunks))
	mockContext.From = mockSource
	for index, chunk := range chunks {
		mockSource.Data[index] = []byte(chunk)
	}
	return mockContext, mockSource
}

// test readSize == 0 then EOF
// 	- should
// 		1. call src.Read once
// 		2. should not call next(context)
// 		3. should not set context.FirstChunk to false
func Test_Read_With_No_Chunk_And_EOF_Error(testCtx *testing.T) {
	// given
	var (
		mockContext, mockSource = NewTestReadChunkContext([]string{}, io.EOF)
		mockRoute               = contexts.NewMockStage("mockRoute")
		mockComplete            = contexts.NewMockStage("mockComplete")
	)

	// when
	read(mockRoute.MockStage, mockComplete.MockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, 1, mockSource.NumberOfReads)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 0, mockSource.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, 0, mockRoute.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, 1, mockComplete.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Error", testCtx, mockSource.Error, io.EOF)
}

// test readSize == 0 then other error
// 	- should
// 		1. call src.Read once
// 		2. should not call next(context)
// 		3. should not set context.FirstChunk to false
func Test_Read_With_Empty_Chunk_And_Non_EOF_Error(testCtx *testing.T) {
	// given
	var (
		mockContext, mockSource = NewTestReadChunkContext([]string{""}, io.ErrClosedPipe)
		mockRoute               = contexts.NewMockStage("mockRoute")
		mockComplete            = contexts.NewMockStage("mockComplete")
	)

	// when
	read(mockRoute.MockStage, mockComplete.MockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, 2, mockSource.NumberOfReads)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 0, mockSource.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, 0, mockRoute.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, 1, mockComplete.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Error", testCtx, mockSource.Error, io.ErrClosedPipe)
}

// test two chunks then EOF
// 	- should
// 		1. call src.Read twice
// 		2. set context.Data to result from src.Read twice
// 		3. should update context.TotalReadSize with amount read in both
// 		4. should call next(context) twice
// 		5. should set context.FirstChunk to false
func Test_Read_With_Two_Chunks(testCtx *testing.T) {
	// given
	var (
		mockContext, mockSource = NewTestReadChunkContext([]string{"this is the first chunk", "this is the second chunk"}, io.EOF)
		mockRoute               = contexts.NewMockStage("mockRoute")
		mockComplete            = contexts.NewMockStage("mockComplete")
	)

	// when
	read(mockRoute.MockStage, mockComplete.MockStage)(mockContext)

	// then
	assertion.AssertDeepEqual("Correct Number Of Reads", testCtx, 3, mockSource.NumberOfReads)
	assertion.AssertDeepEqual("Correct Number Of Writes", testCtx, 0, mockSource.NumberOfWrites)
	assertion.AssertDeepEqual("Correct Route Call Counter", testCtx, 2, mockRoute.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Complete Call Counter", testCtx, 1, mockComplete.MockStageCallCounter)
	assertion.AssertDeepEqual("Correct Error", testCtx, io.EOF, mockSource.Error)
	assertion.AssertDeepEqual("Correct First Chunk", testCtx, mockSource.Data[0], mockRoute.MockStageChunkContexts[0].Data)
	assertion.AssertDeepEqual("Correct First Chunk - firstChunk Indicator", testCtx, true, mockRoute.MockStageChunkContexts[0].FirstChunk)
	assertion.AssertDeepEqual("Correct Second Chunk", testCtx, mockSource.Data[1], mockRoute.MockStageChunkContexts[1].Data)
	assertion.AssertDeepEqual("Correct Second Chunk - firstChunk Indicator", testCtx, false, mockRoute.MockStageChunkContexts[1].FirstChunk)
	assertion.AssertDeepEqual("Correct Total Read Size", testCtx, int64(len(mockSource.Data[0]) + len(mockSource.Data[1])), mockContext.TotalReadSize)
}
