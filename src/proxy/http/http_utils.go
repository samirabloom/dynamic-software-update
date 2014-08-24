package http

import (
	"net"
	"strings"
	"strconv"
	"regexp"
)

var (
	hostHeaderRegex = regexp.MustCompile("Host: ([A-z0-9-_.:]*)")
)

func UpdateHostHeader(data []byte, to net.Addr) []byte {
	ipAddress := string(to.(*net.TCPAddr).IP)
	if len(strings.Trim(ipAddress, "\u007f\x00\x01")) == 0 {
		ipAddress = "127.0.0.1"
	}
	return hostHeaderRegex.ReplaceAllLiteral(data, []byte("Host: "+ipAddress+":"+strconv.Itoa(to.(*net.TCPAddr).Port)))
}
