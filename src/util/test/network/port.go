package network

import (
	"net"
	"testing"
)

func FindFreeLocalSocket(testCtx *testing.T) *net.TCPAddr {
	//	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(byte(127), byte(0), byte(0), byte(1)), Port: 0})
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		listener, err = net.Listen("tcp6", "[::1]:0")
	}
	if err != nil {
		testCtx.Fatal(err)
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr)
}

func FindFreeLocalSocketRange(testCtx *testing.T, size, maximumAttempts int) []*net.TCPAddr {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		listener, err = net.Listen("tcp6", "[::1]:0")
	}
	if err != nil {
		testCtx.Fatal(err)
	}
	firstPort := listener.Addr().(*net.TCPAddr)
	var freePorts []*net.TCPAddr = make([]*net.TCPAddr, size)
	for i := 0 ; i < size ; i++ {
		nextPort, err := net.ListenTCP("tcp", &net.TCPAddr{IP: firstPort.IP, Port: firstPort.Port+i+1})
		if err != nil {
			if maximumAttempts > 0 {
				return FindFreeLocalSocketRange(testCtx, size, maximumAttempts-1)
			}
		} else {
			freePorts[i] = nextPort.Addr().(*net.TCPAddr)
		}
		nextPort.Close()
	}

	defer listener.Close()
	return freePorts
}

