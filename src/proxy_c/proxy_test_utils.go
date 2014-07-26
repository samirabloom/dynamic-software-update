package proxy_c

import (
	"net"
	"code.google.com/p/go-uuid/uuid"
	"time"
	"container/list"
	"fmt"
)

var uuidGenerator = func(uuidValue uuid.UUID) func() uuid.UUID {
	return func() uuid.UUID {
		return uuidValue
	}
}(uuid.NewUUID())

func printOrder(orderedList list.List, printPrefix string) {
	for element := orderedList.Front(); element != nil; element = element.Next() {
		fmt.Printf("%s - version:\t%f\tuuid:\t%s\n", printPrefix, element.Value.(*RoutingContext).version, element.Value.(*RoutingContext).uuid)
	}
}

func NewTestChunkContext() *chunkContext {
	mockContext := &chunkContext{
		description: "",
		data: make([]byte, 64*1024),
		from: &net.TCPConn{},
		to: &net.TCPConn{},
		err: nil,
		totalReadSize: 0,
		totalWriteSize: 0,
		pipeComplete: make(chan int64, 100),
		firstChunk: true,
		performance: *&performance{
			read: new(int64),
			route: new(int64),
			write: new(int64),
			complete: new(int64),
		},
		routingContext: nil,
	}
	return mockContext
}

func CopyChunkContext(contextToCopy *chunkContext) *chunkContext {
	copiedChunkContext := &chunkContext{
		description: contextToCopy.description,
		data: make([]byte, 64*1024),
		from: contextToCopy.from, // todo warning not copied correctly
		to: contextToCopy.to, // todo warning not copied correctly
		err: contextToCopy.err,
		totalReadSize: contextToCopy.totalReadSize,
		totalWriteSize: contextToCopy.totalWriteSize,
		pipeComplete: contextToCopy.pipeComplete, // todo warning not copied correctly
		firstChunk: contextToCopy.firstChunk,
		performance: *&performance{
			read: new(int64),
			route: new(int64),
			write: new(int64),
			complete: new(int64),
		},
	}
	*copiedChunkContext.performance.read = *contextToCopy.performance.read
	*copiedChunkContext.performance.route = *contextToCopy.performance.route
	*copiedChunkContext.performance.write = *contextToCopy.performance.write
	*copiedChunkContext.performance.complete = *contextToCopy.performance.complete
	amountCopied := copy(copiedChunkContext.data, contextToCopy.data)
	copiedChunkContext.data = copiedChunkContext.data[0:amountCopied]
	return copiedChunkContext
}

type mockStage struct {
	mockStageCallCounter      int
	mockStageCallChannel      chan int
	mockStageChunkContexts    []*chunkContext
	description               string
}

func NewMockStage(description string) *mockStage {
	return &mockStage{
		mockStageCallCounter: 0,
		mockStageCallChannel: make(chan int, 5),
		mockStageChunkContexts: make([]*chunkContext, 5),
		description: description,
	}
}

func (mockStage *mockStage) mockStage(context *chunkContext) {
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