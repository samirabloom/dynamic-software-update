package http

import "regexp"

var (
	hostHeaderRegex = regexp.MustCompile("Host: ([A-z0-9-_.:]*)")
)

func UpdateHostHeader(data []byte, host string, port string) []byte {
	return hostHeaderRegex.ReplaceAllLiteral(data, []byte("Host: "+host+":"+port))
}
