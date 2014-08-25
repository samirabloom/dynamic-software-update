package mock

import (
	"time"
	"net"
	"io"
)

type MockReader struct {
	Data  []byte
	Error error
}

func (mockReader MockReader) Read(readBuffer []byte) (n int, err error) {
	return copy(readBuffer, mockReader.Data), mockReader.Error
}

type MockWriter struct {
	Data  []byte
	Error error
}

func (mockWriter MockWriter) Write(writeBuffer []byte) (n int, err error) {
	return copy(mockWriter.Data, writeBuffer), mockWriter.Error
}

type MockConn struct {
	*net.TCPConn
	Data            [][]byte
	Error           error
	LocalAddress    net.Addr
	RemoteAddress   net.Addr
	NumberOfReads   int
	NumberOfWrites  int
	ShortWrite      bool
	ReadClosed      bool
	WriteClosed     bool
	Closed          bool
}

func NewMockConn(err error, size int) *MockConn {
	return &MockConn{
		Data: make([][]byte, size),
		Error: err,
	}
}

func (mockConn *MockConn) Read(readBuffer []byte) (n int, err error) {
	mockConn.NumberOfReads++
	if mockConn.NumberOfReads <= len(mockConn.Data) {
		return copy(readBuffer, mockConn.Data[mockConn.NumberOfReads-1]), nil
	} else {
		return 0, mockConn.Error
	}
}

func (mockConn *MockConn) Write(writeBuffer []byte) (n int, err error) {
	if mockConn.ShortWrite {
		mockConn.Data[mockConn.NumberOfWrites] = make([]byte, len(writeBuffer)/2)
		writeSize := copy(mockConn.Data[mockConn.NumberOfWrites], writeBuffer)
		mockConn.NumberOfWrites++
		return writeSize, nil
	} else {
		mockConn.Data[mockConn.NumberOfWrites] = make([]byte, len(writeBuffer))
		writeSize := copy(mockConn.Data[mockConn.NumberOfWrites], writeBuffer)
		mockConn.NumberOfWrites++
		return writeSize, mockConn.Error
	}
}

func (mockConn *MockConn) Close() error {
	mockConn.Closed = true
	return mockConn.Error
}

func (mockConn *MockConn) LocalAddr() net.Addr {
	return mockConn.LocalAddress
}

func (mockConn *MockConn) RemoteAddr() net.Addr {
	return mockConn.RemoteAddress
}

func (mockConn *MockConn) SetDeadline(t time.Time) error {
	return mockConn.Error
}

func (mockConn *MockConn) SetReadDeadline(t time.Time) error {
	return mockConn.Error
}

func (mockConn *MockConn) SetWriteDeadline(t time.Time) error {
	return mockConn.Error
}

func (mockConn *MockConn) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}

func (mockConn *MockConn) CloseRead() error {
	mockConn.ReadClosed = true
	return mockConn.Error
}

func (mockConn *MockConn) CloseWrite() error {
	mockConn.WriteClosed = true
	return mockConn.Error
}

func (mockConn *MockConn) SetLinger(sec int) error {
	return mockConn.Error
}

func (mockConn *MockConn) SetKeepAlive(keepalive bool) error {
	return mockConn.Error
}

func (mockConn *MockConn) SetKeepAlivePeriod(d time.Duration) error {
	return mockConn.Error
}

func (mockConn *MockConn) SetNoDelay(noDelay bool) error {
	return mockConn.Error
}
