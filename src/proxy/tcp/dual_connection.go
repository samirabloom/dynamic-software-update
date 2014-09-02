package tcp

import (
	log "proxy/log"
	"net"
	"strconv"
	"time"
	"fmt"
	"regexp"
	"io"
	"errors"
	"proxy/http"
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
	Connections         []TCPConnection
	Hosts               []string
	Ports               []string
	SuccessfulIndex     int
}

func NewDualTCPConnection(addresses []*net.TCPAddr, hosts, ports []string) *DualTCPConnection {
	connections := make([]TCPConnection, len(addresses))
	for index, address := range addresses {
		connection, err := net.DialTCP("tcp", nil, address)
		if err != nil {
			log.LoggerFactory().Error("error [%s] making connection for address [%s]\n", err, address)
		}
		connections[index] = connection
	}
	return &DualTCPConnection{
		Connections:        connections,
		Hosts:              hosts,
		Ports:              ports,
		SuccessfulIndex:    -1,
	}

}

type DualAddr struct {
	Addresses []net.Addr
}

func (da *DualAddr) Network() string {
	var network = ""
	for index, address := range da.Addresses {
		if index > 0 {
			network += ", "
		}
		if address != nil {
			network += address.Network()
		}
	}
	return network
}

func (da *DualAddr) String() string {
	var network = ""
	for index, address := range da.Addresses {
		if index > 0 {
			network += ", "
		}
		if address != nil {
			network += address.String()
		}
	}
	return network
}

func (dualTCPConnection *DualTCPConnection) Read(readBuffer []byte) (int, error) {
	if dualTCPConnection.SuccessfulIndex == -1 {
		return dualTCPConnection.readMultiple(readBuffer)
	} else {
		return dualTCPConnection.Connections[dualTCPConnection.SuccessfulIndex].Read(readBuffer)
	}
}

func (dualTCPConnection *DualTCPConnection) readMultiple(readBuffer []byte) (int, error) {
	var (
		groupCount int
		groupErr error
		groupReadBuffer []byte
	)
	for index, connection := range dualTCPConnection.Connections {
		connectionReadBuffer := make([]byte, len(readBuffer))
		count, err := connection.Read(connectionReadBuffer)
		if readSuccessful(connectionReadBuffer) {
			groupCount = count
			groupErr = err
			groupReadBuffer = connectionReadBuffer
			dualTCPConnection.SuccessfulIndex = index
		}
		// allow for all reads to fail
		if groupErr == nil {
			groupErr = err
		}
		if groupReadBuffer == nil {
			groupCount = count
			groupReadBuffer = connectionReadBuffer
		}
	}

	copy(readBuffer, groupReadBuffer)
	return groupCount, groupErr
}

func readSuccessful(readBuffer []byte) bool {
	statusCodeRegex := regexp.MustCompile("HTTP/[0-9].[0-9] ([a-z0-9-]*) .*")
	subMatches := statusCodeRegex.FindSubmatch(readBuffer)
	if len(subMatches) > 1 {
		statusCode, conversionErr := strconv.Atoi(string(subMatches[1]))
		return conversionErr == nil && int(statusCode / 100) != 5
	}
	return false
}

func (dualTCPConnection *DualTCPConnection) Write(writeBuffer []byte) (int, error) {
	var (
		maxCount int
		errorsList string
	)
	for index, connection := range dualTCPConnection.Connections {
		writeBuffer = http.UpdateHostHeader(writeBuffer, dualTCPConnection.Hosts[index], dualTCPConnection.Ports[index], false)

		count, err := connection.Write(writeBuffer)
		if err != nil {
			if len(errorsList) > 0 {
				errorsList += ", "
			}
			errorsList += fmt.Sprintf("connections[%d]: %s -> %s - error: %s", index, connection.LocalAddr(), connection.RemoteAddr(), err)
		} else if len(writeBuffer) != count {
			return count, io.ErrShortWrite
		}
		if count > maxCount {
			maxCount = count
		}
	}
	if len(errorsList) > 0 {
		return maxCount, errors.New(errorsList)
	} else {
		return maxCount, nil
	}
}

func (dualTCPConnection *DualTCPConnection) String() string {
	var output string = ""
	output += "\n{\n"
	for index, connection := range dualTCPConnection.Connections {
		output += fmt.Sprintf("\t connections[%d]: %s -> %s\n", index, connection.LocalAddr(), connection.RemoteAddr())
	}
	output += "}\n"
	return output
}

func (dualTCPConnection *DualTCPConnection) Close() error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.Close()
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dualTCPConnection *DualTCPConnection) LocalAddr() net.Addr {
	addresses := make([]net.Addr, len(dualTCPConnection.Connections))
	for index, connection := range dualTCPConnection.Connections {
		AllowForNilConnection(connection, func(connection TCPConnection) {
			addresses[index] = connection.LocalAddr()
		})
	}
	return &DualAddr{Addresses: addresses}
}

func (dualTCPConnection *DualTCPConnection) RemoteAddr() net.Addr {
	addresses := make([]net.Addr, len(dualTCPConnection.Connections))
	for index, connection := range dualTCPConnection.Connections {
		AllowForNilConnection(connection, func(connection TCPConnection) {
			addresses[index] = connection.RemoteAddr()
		})
	}
	return &DualAddr{Addresses: addresses}
}

func (dualTCPConnection *DualTCPConnection) SetDeadline(time time.Time) error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.SetDeadline(time)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dualTCPConnection *DualTCPConnection) SetReadDeadline(time time.Time) error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.SetReadDeadline(time)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dualTCPConnection *DualTCPConnection) SetWriteDeadline(time time.Time) error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.SetWriteDeadline(time)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dualTCPConnection *DualTCPConnection) ReadFrom(reader io.Reader) (int64, error) {
	panic("ReadFrom not supported for DualTCPConnection")
	return 0, nil
}

func (dualTCPConnection *DualTCPConnection) CloseRead() error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.CloseRead()
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dualTCPConnection *DualTCPConnection) CloseWrite() error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.CloseWrite()
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dualTCPConnection *DualTCPConnection) SetLinger(seconds int) error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.SetLinger(seconds)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dualTCPConnection *DualTCPConnection) SetKeepAlive(keepAlive bool) error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.SetKeepAlive(keepAlive)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dualTCPConnection *DualTCPConnection) SetKeepAlivePeriod(duration time.Duration) error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.SetKeepAlivePeriod(duration)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dualTCPConnection *DualTCPConnection) SetNoDelay(noDelay bool) error {
	for _, connection := range dualTCPConnection.Connections {
		var err error
		AllowForNilConnection(connection, func(connection TCPConnection) {
			err = connection.SetNoDelay(noDelay)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ==== DUAL_TCP_CONNECTION - END
