package network

import (
	"net"
	"testing"
)

func FindFreeLocalSocket(testCtx *testing.T) (string, *net.TCPAddr) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		testCtx.Fatal(err)
	}
	localAddress := listener.Addr().String()
	defer listener.Close()
	localTCPAddress, err := net.ResolveTCPAddr("tcp", localAddress)
	if err != nil {
		testCtx.Fatal(err)
	}
	return localAddress, localTCPAddress
}

