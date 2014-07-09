package mock

import (
	"time"
	"net"
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
	Data           [][]byte
	Error          error
	LocalAddress   net.Addr
	RemoteAddress  net.Addr
	NumberOfReads  int
	NumberOfWrites int
	shortWrite     bool
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
	if mockConn.shortWrite {
		writeSize := copy(mockConn.Data[mockConn.NumberOfWrites], writeBuffer)
		mockConn.NumberOfWrites++
		return writeSize/2, nil
	} else {
		writeSize := copy(mockConn.Data[mockConn.NumberOfWrites], writeBuffer)
		mockConn.NumberOfWrites++
		return writeSize, mockConn.Error
	}
}

func (mockConn *MockConn) Close() error {
	return mockConn.Error
}

func (mockConn *MockConn) LocalAddr() net.Addr {
	return mockConn.LocalAddress
}

func (mockConn *MockConn) RemoteAddr() net.Addr {
	return mockConn.RemoteAddress
}

func (mockConn *MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (mockConn *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (mockConn *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
