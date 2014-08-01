package dual_TCP_connection

import (
	log "proxy/log"
	"net"
	"strconv"
	"time"
	"fmt"
	"regexp"
	"io"
)


// ==== DUAL_TCP_CONNECTION - START

type TCPConnection interface {
	Read(readBuffer []byte) (n int, err error)
	Write(writeBuffer []byte) (n int, err error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	ReadFrom(r io.Reader) (int64, error)
	CloseRead() error
	CloseWrite() error
	SetLinger(sec int) error
	SetKeepAlive(keepAlive bool) error
	SetKeepAlivePeriod(d time.Duration) error
	SetNoDelay(noDelay bool) error
}



type DualTCPConnection struct {
	expectedStatusCode                     int
	connection1                            TCPConnection
	connection2                            TCPConnection
}

func NewDualTCPConnection(expectedStatusCode int, firstConnection *net.TCPAddr, secondConnection *net.TCPAddr) *DualTCPConnection {
	server1, server1Err := net.DialTCP("tcp", nil, firstConnection)
	if server1Err != nil {
		log.LoggerFactory().Error("error making 1st connection of DualConnection /%v:\n", server1Err)
	}
	server2, server2Err := net.DialTCP("tcp", nil, secondConnection)
	if server2Err != nil {
		log.LoggerFactory().Error("error making 2nd connection of DualConnection /%v:\n", server2Err)
	}
	return &DualTCPConnection{
		expectedStatusCode: expectedStatusCode,
		connection1:        server1,
		connection2:        server2,
	}

}

func (dualTCPConnection *DualTCPConnection) Read(readBuffer []byte) (int, error) {
	read1Buffer := make([]byte, len(readBuffer))
	connection1Read, err1 := dualTCPConnection.connection1.Read(read1Buffer)
	read2Buffer := make([]byte, len(readBuffer))
	connection2Read, err2 := dualTCPConnection.connection2.Read(read2Buffer)

	statusCodeRegex := regexp.MustCompile("HTTP/[0-9].[0-9] ([a-z0-9-]*) .*")
	statusCode1Matches := statusCodeRegex.FindSubmatch(read1Buffer)
	statusCode2Matches := statusCodeRegex.FindSubmatch(read2Buffer)
	if len(statusCode1Matches) >= 2 {
		statusCode1Match := string(statusCode1Matches[1])
		statusCode1, _ := strconv.Atoi(statusCode1Match)
		log.LoggerFactory().Debug("statusCode1 found is: %d", statusCode1)
		if err1 == nil && statusCode1 == dualTCPConnection.expectedStatusCode {
			copy(readBuffer, read1Buffer)
			return connection1Read, err1
		}
	}
	if len(statusCode2Matches) >= 2 {
		statusCode2Match := string(statusCode2Matches[1])
		statusCode2, _ := strconv.Atoi(statusCode2Match)
		log.LoggerFactory().Debug("statusCode2 found is: %d", statusCode2)
		if err2 == nil && statusCode2 == dualTCPConnection.expectedStatusCode {
			copy(readBuffer, read2Buffer)
			return connection2Read, err2
		}
	}
	// TODO add the scenario when non of the servers have the right status code
	return connection1Read, err1
}

func (dualTCPConnection *DualTCPConnection) Write(writeBuffer []byte) (int, error) {
	connection1Write, err1 := dualTCPConnection.connection1.Write(writeBuffer)
	connection2Write, err2 := dualTCPConnection.connection2.Write(writeBuffer)

	if err2 == nil {
		return connection2Write, err2
	} else if err1 == nil {
		return connection1Write, err1
	} else {
		log.LoggerFactory().Error("DualConnection write error %v: %s\n", err1, err1)
		return connection1Write, err1
	}
}

func (dualTCPConnection *DualTCPConnection) String() string {
	var output string = ""
	output += "\n{\n"
	if dualTCPConnection.connection1 != nil && dualTCPConnection.connection1.(*net.TCPConn) != nil && dualTCPConnection.connection1.LocalAddr() != nil && dualTCPConnection.connection1.RemoteAddr() != nil {
		output += fmt.Sprintf("\t connection1: %s -> %s\n", dualTCPConnection.connection1.LocalAddr(), dualTCPConnection.connection1.RemoteAddr())
	}
	if dualTCPConnection.connection2 != nil && dualTCPConnection.connection2.(*net.TCPConn) != nil && dualTCPConnection.connection2.LocalAddr() != nil && dualTCPConnection.connection2.RemoteAddr() != nil {
		output += fmt.Sprintf("\t connection2: %s -> %s\n", dualTCPConnection.connection2.LocalAddr(), dualTCPConnection.connection2.RemoteAddr())
	}
	output += "}\n"
	return output
}

func (dualTCPConnection *DualTCPConnection) Close() {
	dualTCPConnection.connection1.Close()
	dualTCPConnection.connection2.Close()
}

func (dualTCPConnection *DualTCPConnection) LocalAddr() net.Addr {
	panic("LocalAddr not supported for DualTCPConnection")
	return nil
}

func (dualTCPConnection *DualTCPConnection) RemoteAddr() net.Addr {
	panic("RemoteAddr not supported for DualTCPConnection")
	return nil
}

func (dualTCPConnection *DualTCPConnection) SetDeadline(t time.Time) error {
	error1 := dualTCPConnection.connection1.SetDeadline(t)
	error2 := dualTCPConnection.connection2.SetDeadline(t)
	if error1 != nil {
		return error1
	}
	return error2
}

func (dualTCPConnection *DualTCPConnection) SetReadDeadline(t time.Time) error {
	error1 := dualTCPConnection.connection1.SetReadDeadline(t)
	error2 := dualTCPConnection.connection2.SetReadDeadline(t)
	if error1 != nil {
		return error1
	}
	return error2
}

func (dualTCPConnection *DualTCPConnection) SetWriteDeadline(t time.Time) error {
	error1 := dualTCPConnection.connection1.SetWriteDeadline(t)
	error2 := dualTCPConnection.connection2.SetWriteDeadline(t)
	if error1 != nil {
		return error1
	}
	return error2
}

func (dualTCPConnection *DualTCPConnection) ReadFrom(r io.Reader) (int64, error) {
	panic("ReadFrom not supported for DualTCPConnection")
	return 0, nil
}

func (dualTCPConnection *DualTCPConnection) CloseRead() error {
	error1 := dualTCPConnection.connection1.CloseRead()
	error2 := dualTCPConnection.connection2.CloseRead()
	if error1 != nil {
		return error1
	}
	return error2
}

func (dualTCPConnection *DualTCPConnection) CloseWrite() error {
	error1 := dualTCPConnection.connection1.CloseWrite()
	error2 := dualTCPConnection.connection2.CloseWrite()
	if error1 != nil {
		return error1
	}
	return error2
}

func (dualTCPConnection *DualTCPConnection) SetLinger(sec int) error {
	error1 := dualTCPConnection.connection1.SetLinger(sec)
	error2 := dualTCPConnection.connection2.SetLinger(sec)
	if error1 != nil {
		return error1
	}
	return error2
}

func (dualTCPConnection *DualTCPConnection) SetKeepAlive(keepalive bool) error {
	error1 := dualTCPConnection.connection1.SetKeepAlive(keepalive)
	error2 := dualTCPConnection.connection2.SetKeepAlive(keepalive)
	if error1 != nil {
		return error1
	}
	return error2
}

func (dualTCPConnection *DualTCPConnection) SetKeepAlivePeriod(d time.Duration) error {
	error1 := dualTCPConnection.connection1.SetKeepAlivePeriod(d)
	error2 := dualTCPConnection.connection2.SetKeepAlivePeriod(d)
	if error1 != nil {
		return error1
	}
	return error2
}

func (dualTCPConnection *DualTCPConnection) SetNoDelay(noDelay bool) error {
	error1 := dualTCPConnection.connection1.SetNoDelay(noDelay)
	error2 := dualTCPConnection.connection2.SetNoDelay(noDelay)
	if error1 != nil {
		return error1
	}
	return error2
}

// ==== DUAL_TCP_CONNECTION - END
