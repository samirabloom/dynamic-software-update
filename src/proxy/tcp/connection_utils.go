package tcp

import "net"

func AllowForNilConnection(connection TCPConnection, operation func (TCPConnection)) {
	_, assertion := connection.(*net.TCPConn)
	if assertion {
		if connection.(*net.TCPConn) != nil {
			operation(connection);
		}
	}

	_, assertion = connection.(*TCPConnAndName)
	if assertion {
		if connection.(*TCPConnAndName) != nil {
			if connection.(*TCPConnAndName).TCPConn != (*net.TCPConn)(nil) {
				operation(connection.(*TCPConnAndName).TCPConn);
			}
		}
	}

	_, assertion = connection.(*DualTCPConnection)
	if assertion {
		if connection.(*DualTCPConnection) != nil {
			operation(connection);
		}
	}
}

func IsDualConnection(connection TCPConnection) bool {
	_, isDualConnection := connection.(*DualTCPConnection)
	return isDualConnection
}
