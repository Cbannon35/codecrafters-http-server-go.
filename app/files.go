package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleFileRoute(conn net.Conn, req *HTTPRequest) {
	fileName := strings.TrimPrefix(req.Path, "/files/")
	dir := os.Args[2]
	switch method := req.Method; {
	case method == "GET":
		handleFileRouteGET(conn, req, dir, fileName)

	case method == "POST":
		handleFileRoutePOST(conn, req, dir, fileName)
	}
}

func handleFileRouteGET(conn net.Conn, req *HTTPRequest, dir string, fileName string) {
	data, err := os.ReadFile(dir + fileName)
	if err != nil {
		writeToConnection(conn, []byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}
	response := createResponse("HTTP/1.1 200 OK", []string{"Content-Type: application/octet-stream", "Content-Length: " + fmt.Sprint(len(data))}, string(data))
	writeToConnection(conn, []byte(response))
}

func handleFileRoutePOST(conn net.Conn, req *HTTPRequest, dir string, fileName string) {
	content := strings.Trim(req.Body, "\x00")
	err := os.WriteFile(dir+fileName, []byte(content), 0644)
	if err != nil {
		writeToConnection(conn, []byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}
	writeToConnection(conn, []byte("HTTP/1.1 201 Created\r\n\r\n"))
}