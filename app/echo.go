package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"strings"
)

func handleEchoRoute(conn net.Conn, req *HTTPRequest) {
	content := strings.TrimPrefix(req.Path, "/echo/")
	headers := []string{}
	
	encodings := req.Headers["Accept-Encoding"]
	for _, encoding := range strings.Split(encodings, ",") {
		if strings.TrimSpace(encoding) == "gzip" {
			var buffer bytes.Buffer
			w := gzip.NewWriter(&buffer)
			w.Write([]byte(content))
			w.Close()
			content = buffer.String()
			headers = append(headers, "Content-Encoding: gzip")
			break
		}
	}
	headers = append(headers, "Content-Type: text/plain", "Content-Length: " + fmt.Sprint(len(content)))
	response := createResponse("HTTP/1.1 200 OK", headers, content)
	writeToConnection(conn, []byte(response))
}