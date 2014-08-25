package tcp

import "net"

type TCPConnAndName struct {
	*net.TCPConn
	Host string
	Port string
}
