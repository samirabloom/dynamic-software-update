package mock

import (
	"net/http"
)

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{WrittenBodyBytes: make(map[int][]byte), Headers: make(http.Header), ResponseCodes: make(map[int]int)}
}

type MockResponseWriter struct {
	WrittenBodyBytes map[int][]byte
	Headers         http.Header
	ResponseCodes   map[int]int
}

func (rw *MockResponseWriter) Header() http.Header {
	return rw.Headers;
}

func (rw *MockResponseWriter) Write(data []byte) (int, error) {
	currentSize := len(rw.WrittenBodyBytes)
	rw.WrittenBodyBytes[currentSize] = make([]byte, len(data))
	copy(rw.WrittenBodyBytes[currentSize], data)
	return len(data), nil;
}

func (rw *MockResponseWriter) WriteHeader(code int) {
	rw.ResponseCodes[len(rw.ResponseCodes)] = code
}

type MockBody struct {
	BodyBytes []byte
}

func (body *MockBody) Read(data []byte) (n int, err error) {
	copy(data, body.BodyBytes)
	return len(body.BodyBytes), nil;
}

func (body *MockBody) Close() error {
	return nil;
}

type Writer interface {
	Write(p []byte) (n int, err error)
}
