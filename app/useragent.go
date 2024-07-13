package main

import (
	"fmt"
	"net"
)

func handleUserAgentRoute(conn net.Conn, req *HTTPRequest) {
	body := req.Headers["User-Agent"]
	response := createResponse("HTTP/1.1 200 OK", []string{"Content-Type: text/plain", "Content-Length: " + fmt.Sprint(len(body))}, body)
	writeToConnection(conn, []byte(response))
}