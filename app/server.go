package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// @Citation: https://app.codecrafters.io/users/Powerisinschool
type HTTPRequest struct {
	Method    string
	Path      string
	Headers   map[string]string
	Body      string
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Connection accepted")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error handling request: ", err.Error())
		os.Exit(1)
	}
	req, err := parseRequest(buf)

	if err != nil {
		fmt.Fprintln(conn, "reading standard input:", err)
		os.Exit(1)
	}

	switch path := req.Path; {
	case strings.HasPrefix(path, "/files"):
		handleFileRoute(conn, req)
	case strings.HasPrefix(path, "/echo"):
		handleEchoRoute(conn, req)
	case path == "/user-agent":
		handleUserAgentRoute(conn, req)
	case path == "/":
		writeToConnection(conn, []byte("HTTP/1.1 200 OK\r\n\r\n"))
	default:
		writeToConnection(conn, []byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}