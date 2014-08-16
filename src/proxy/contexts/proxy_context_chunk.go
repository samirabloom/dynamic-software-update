package contexts

import (
	"fmt"
	"proxy/tcp"
	"strings"
	"time"
)

type Direction bool

const (
	ClientToServer Direction = true
	ServerToClient Direction = false
)

var DirectionToDescription = map[Direction]string {
	ClientToServer: "client -> server",
	ServerToClient: "client <- server",
}

type RoutingContext struct {
	Headers []string
}

type ChunkContext struct {
	Data                   []byte
	To                     tcp.TCPConnection
	From                   tcp.TCPConnection
	Err                    error
	TotalReadSize          int64
	TotalWriteSize         int64
	PipeComplete           chan int64
	FirstChunk             bool
	RoutingContext         *RoutingContext
	Direction              Direction
}

func (context *ChunkContext) Close() {
	// close sockets
	context.From.Close()
	if context.To != nil {
		context.To.Close()
	}
}

func (context *ChunkContext) String() string {
	var output string = ""
	output += "\n{\n"
	output += "\t direction: "+DirectionToDescription[context.Direction]+"\n"
	if len(context.Data) > 0 {
		output += "\t data:\n\t\t"+strings.Replace(string(context.Data), "\n", "\n\t\t", -1)
	}
	output += "\n"
	if context.From != nil && context.From.LocalAddr() != nil && context.From.RemoteAddr() != nil {
		output += fmt.Sprintf("\t from: %s -> %s\n", context.From.LocalAddr(), context.From.RemoteAddr())
	}
	if context.To != nil && context.To.LocalAddr() != nil && context.To.RemoteAddr() != nil {
		output += fmt.Sprintf("\t to: %s -> %s\n", context.To.LocalAddr(), context.To.RemoteAddr())
	}
	output += fmt.Sprintf("\t totalReadSize: %d\n", context.TotalReadSize)
	output += fmt.Sprintf("\t totalWriteSize: %d\n", context.TotalWriteSize)
	output += "}\n"
	return output
}

func NewForwardPipeChunkContext(from tcp.TCPConnection, pipeComplete chan int64) *ChunkContext {
	return &ChunkContext{
		Data:           make([]byte, 64*1024),
		From:           from,
		PipeComplete:   pipeComplete,
		FirstChunk:     true,
		RoutingContext: nil,
		Direction:        ClientToServer,
	}
}

func NewBackPipeChunkContext(forwardContext *ChunkContext) *ChunkContext {
	return &ChunkContext{
		Data:           make([]byte, 64*1024),
		From:           forwardContext.To,
		To:             forwardContext.From,
		PipeComplete:   forwardContext.PipeComplete,
		FirstChunk:     true,
		RoutingContext: forwardContext.RoutingContext,
		Direction:        ServerToClient,
	}
}

// --- TEST UTILS - START

func CopyChunkContext(contextToCopy *ChunkContext) *ChunkContext {
	copiedChunkContext := &ChunkContext{
		Data: make([]byte, 64*1024),
		From: contextToCopy.From, // todo warning not copied correctly
		To: contextToCopy.To, // todo warning not copied correctly
		Err: contextToCopy.Err,
		TotalReadSize: contextToCopy.TotalReadSize,
		TotalWriteSize: contextToCopy.TotalWriteSize,
		PipeComplete: contextToCopy.PipeComplete, // todo warning not copied correctly
		FirstChunk: contextToCopy.FirstChunk,
	}
	amountCopied := copy(copiedChunkContext.Data, contextToCopy.Data)
	copiedChunkContext.Data = copiedChunkContext.Data[0:amountCopied]
	return copiedChunkContext
}

type MockStage struct {
	MockStageCallCounter      int
	MockStageCallChannel      chan int
	MockStageChunkContexts    []*ChunkContext
	Description               string
}

func NewMockStage(description string) *MockStage {
	return &MockStage{
		MockStageCallCounter: 0,
		MockStageCallChannel: make(chan int, 5),
		MockStageChunkContexts: make([]*ChunkContext, 5),
		Description: description,
	}
}

func (mockStage *MockStage) MockStage(context *ChunkContext) {
	mockStage.MockStageChunkContexts[mockStage.MockStageCallCounter] = CopyChunkContext(context)
	mockStage.MockStageCallCounter++
	mockStage.MockStageCallChannel <- mockStage.MockStageCallCounter
}

func (mockStage *MockStage) Close(secondsDelay time.Duration) {
	go func() {
		time.Sleep(time.Second * secondsDelay)
		close(mockStage.MockStageCallChannel)
	}()
}

// --- TEST UTILS - END

// ==== CHUNK_CONTEXT - END
