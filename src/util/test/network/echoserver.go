package network

import (
	"io"
	byteutil "util/byte"
	"bytes"
	"fmt"
	"net"
	"testing"
)

func Run(testCtx *testing.T, address *net.TCPAddr) {
	go func() {
		listener, err := net.ListenTCP("tcp", address)
		if err != nil {
			testCtx.Fatal(err)
		}
		for {
			client, err := listener.Accept()
			if err != nil {
				return
			}
			go func(client net.Conn) {
				data := make([]byte, 32*1024)
				for {
					readSize, readError := client.Read(data)
					if readSize > 0 {
//						fmt.Printf("Before insert - ES: \n%s\n", data)
						echoServerHeader := []byte(fmt.Sprintf("X-EchoServer: %s\n", address))
						searchString := "\n"
						insertLocation := bytes.Index(data, []byte(searchString))
						if insertLocation > 0 {
							byteutil.Insert(data[0:readSize], insertLocation+len(searchString), echoServerHeader)
							readSize += len(echoServerHeader)
						}
						writeSize, writeError := client.Write(data[0:readSize])
						if writeError != nil {
							testCtx.Logf("error in echo server: %v\n", writeError.Error())
							break
						}
						if readSize != writeSize {
							testCtx.Logf("error in echo server: %v\n", io.ErrShortWrite.Error())
							break
						}
					}
					if readError == io.EOF {
						break
					}
					if readError != nil {
						testCtx.Logf("error in echo server: %v\n", readError.Error())
						break
					}
				}
				client.Close()
			}(client)
		}
	}()
}
