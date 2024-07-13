package main

import (
	"fmt"
	"net"
	"strings"
)

func writeToConnection(conn net.Conn, message []byte) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Can't write to connection: ", err.Error())
		return
	}
}

func createResponse(status string, headers []string, body string) string {
	response := string(status) + "\r\n"
	for _, header := range headers {
		response += header + "\r\n"
	}
	response += "\r\n"
	if body != "" {
		response += body
	}
	return response
}


func parseRequest(buf []byte) (*HTTPRequest, error) {
	var req HTTPRequest = HTTPRequest{}
	req.Headers = make(map[string]string)
	lines := strings.Split(string(buf), "\r\n")
	fmt.Println(lines)
	for i, line := range lines {
		if i == 0 {
			req.Method = strings.Split(line, " ")[0]
			req.Path = strings.Split(line, " ")[1]
			continue
		}
		if line == "" {
			req.Body = strings.Join(lines[i+1:], "\r\n")
			break
		}
		headers := strings.Split(line, ": ")
		if len(headers) < 2 {
			req.Body = headers[0]
			break
		}
		req.Headers[headers[0]] = headers[1]
	}
	return &req, nil
}