package http

import "regexp"

var (
	hostHeaderRegex = regexp.MustCompile("Host: ([A-z0-9-_.:]*)")
)

func UpdateHostHeader(data []byte, host string, port string, rewriteHeader bool) []byte {
	if rewriteHeader {
		if port == "80" {
			return hostHeaderRegex.ReplaceAllLiteral(data, []byte("Host: "+host))

		} else {
			return hostHeaderRegex.ReplaceAllLiteral(data, []byte("Host: "+host+":"+port))
		}
	} else {
		return data
	}
}
