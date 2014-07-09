package mock

import (
	"net/http"
)

type MockResponseWriter struct {
	WritenBodyBytes map[int][]byte
}

func (rw *MockResponseWriter) Header() http.Header {
	return nil;
}

func (rw *MockResponseWriter) Write(data []byte) (int, error) {
	rw.WritenBodyBytes[len(rw.WritenBodyBytes)] = data
	return len(data), nil;
}

func (rw *MockResponseWriter) WriteHeader(int) {
	return;
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
