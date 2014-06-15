package main

import (
	"fmt"
	"net"
	"time"
	"strconv"
)

func main() {
	//create listener
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error Listening", err.Error())
		return // terminate program
	}

	// listen and accept connections from clients:
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // terminate program
		}
		go doServerStuff(conn)
	}
}

func doServerStuff(conn net.Conn) {
	body := "<html><head></head><body><div>This is a simple go server running in a docker containers</div></body></html>\n"
	conn.Write([]byte("HTTP/1.1 200 Ok\n" +
			"Content-Length: " + strconv.Itoa(len(body)) + "\n" +
			"Set-Cookie: JSESSIONID=F9ADC2A30C8255A49061992F48C8AAAB; Path=/; Secure; HttpOnly\n" +
			"\r\n" +
			body))

	time.Sleep(10 * time.Millisecond)
	conn.Close()
}
