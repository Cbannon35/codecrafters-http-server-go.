package main

import (
	"fmt"
	"net"
	"strings"
)

func handleEchoRoute(conn net.Conn, req *HTTPRequest) {
	content := strings.TrimPrefix(req.Path, "/echo/")
	headers := []string{"Content-Type: text/plain", "Content-Length: " + fmt.Sprint(len(content))}
	if req.Headers["Accept-Encoding"] == "gzip" {
		// content = gzipContent(content)
		headers = append(headers, "Content-Encoding: gzip")
	}
	response := createResponse("HTTP/1.1 200 OK", headers, content)
	writeToConnection(conn, []byte(response))
}