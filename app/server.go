package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func setupListener() net.Listener {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	return listener
}

func writeToConnection(conn net.Conn, message []byte) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Can't write to connection: ", err.Error())
		return
	}
}

func parseRequest(conn net.Conn) []string {
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		fmt.Println("Can't read from connection: ", err.Error())
		return nil
	}
	return strings.Split(string(req[:n]), "\r\n")
}

func main() {
	l := setupListener()
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
	fmt.Println("Handling connection")
	request := parseRequest(conn)

	if request == nil || len(request) < 3 {
		fmt.Println("Invalid request?")
		return
	}

	requestLine := strings.Split(request[0], " ")

	if requestLine[1] != "/" {
		writeToConnection(conn, []byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	} else {
		writeToConnection(conn, []byte("HTTP/1.1 200 OK\r\n\r\n"))
	}
	conn.Close()
}