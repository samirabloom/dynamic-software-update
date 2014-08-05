package stages

import (
	"fmt"
	"proxy/tcp"
	"net"
	"strings"
	"time"
)

type RoutingContext struct {
	headers []string
}

type ChunkContext struct {
	description            string
	data                   []byte
	to                     tcp.TCPConnection
	from                   tcp.TCPConnection
	err                    error
	totalReadSize          int64
	totalWriteSize         int64
	pipeComplete           chan int64
	firstChunk             bool
	routingContext         *RoutingContext
	clientToServer         bool
}

func (context *ChunkContext) Close() {
	// close sockets
	context.from.Close()
	if context.to != nil && context.to.(*net.TCPConn) != nil {
		context.to.Close()
	}
}

func (context *ChunkContext) String() string {
	var output string = ""
	output += "\n{\n"
	output += fmt.Sprintf("\t description: %s\n", context.description)
	if context.clientToServer {
		output += "\t direction: client->server\n"
	} else {
		output += "\t direction: server->client\n"
	}
	if len(context.data) > 0 {
		output += "\t data:\n\t\t"+strings.Replace(string(context.data), "\n", "\n\t\t", -1)
	}
	output += "\n"
	if context.from != nil && context.from.(*net.TCPConn) != nil && context.from.LocalAddr() != nil && context.from.RemoteAddr() != nil {
		output += fmt.Sprintf("\t from: %s -> %s\n", context.from.LocalAddr(), context.from.RemoteAddr())
	}
	if context.to != nil && context.to.(*net.TCPConn) != nil && context.to.LocalAddr() != nil && context.to.RemoteAddr() != nil {
		output += fmt.Sprintf("\t to: %s -> %s\n", context.to.LocalAddr(), context.to.RemoteAddr())
	}
	output += fmt.Sprintf("\t totalReadSize: %d\n", context.totalReadSize)
	output += fmt.Sprintf("\t totalWriteSize: %d\n", context.totalWriteSize)
	output += "}\n"
	return output
}

func NewForwardPipeChunkContext(from *net.TCPConn, pipeComplete chan int64) *ChunkContext {
	return &ChunkContext{
		description:    "forwardpipe",
		data:           make([]byte, 64*1024),
		from:           from,
		pipeComplete:   pipeComplete,
		firstChunk:     true,
		routingContext: nil,
		clientToServer: true,
	}
}

func NewBackPipeChunkContext(forwardContext *ChunkContext) *ChunkContext {
	return &ChunkContext{
		description:    "backpipe",
		data:           make([]byte, 64*1024),
		from:           forwardContext.to,
		to:             forwardContext.from,
		pipeComplete:   forwardContext.pipeComplete,
		firstChunk:     true,
		routingContext: forwardContext.routingContext,
		clientToServer: false,
	}
}

// --- TEST UTILS - START

func CopyChunkContext(contextToCopy *ChunkContext) *ChunkContext {
	copiedChunkContext := &ChunkContext{
		description: contextToCopy.description,
		data: make([]byte, 64*1024),
		from: contextToCopy.from, // todo warning not copied correctly
		to: contextToCopy.to, // todo warning not copied correctly
		err: contextToCopy.err,
		totalReadSize: contextToCopy.totalReadSize,
		totalWriteSize: contextToCopy.totalWriteSize,
		pipeComplete: contextToCopy.pipeComplete, // todo warning not copied correctly
		firstChunk: contextToCopy.firstChunk,
	}
	amountCopied := copy(copiedChunkContext.data, contextToCopy.data)
	copiedChunkContext.data = copiedChunkContext.data[0:amountCopied]
	return copiedChunkContext
}

type mockStage struct {
	mockStageCallCounter      int
	mockStageCallChannel      chan int
	mockStageChunkContexts    []*ChunkContext
	description               string
}

func NewMockStage(description string) *mockStage {
	return &mockStage{
		mockStageCallCounter: 0,
		mockStageCallChannel: make(chan int, 5),
		mockStageChunkContexts: make([]*ChunkContext, 5),
		description: description,
	}
}

func (mockStage *mockStage) mockStage(context *ChunkContext) {
	mockStage.mockStageChunkContexts[mockStage.mockStageCallCounter] = CopyChunkContext(context)
	mockStage.mockStageCallCounter++
	mockStage.mockStageCallChannel <- mockStage.mockStageCallCounter
}

func (mockStage * mockStage) close(secondsDelay time.Duration) {
	go func() {
		time.Sleep(time.Second * secondsDelay)
		close(mockStage.mockStageCallChannel)
	}()
}

// --- TEST UTILS - END

// ==== CHUNK_CONTEXT - END
