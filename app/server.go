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

func createResponse(status []byte, headers []string, body string) string {
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

	endpoint := requestLine[1]

	if endpoint == "/" {
		response := createResponse([]byte("HTTP/1.1 200 OK"), nil, "")
		writeToConnection(conn, []byte(response))
	} else if strings.HasPrefix(endpoint, "/echo") {
		text := strings.TrimPrefix(endpoint, "/echo/")
		fmt.Println("Echoing: ", text)
		headers := []string{"Content-Type: text/plain", "Content-Length: " + fmt.Sprint(len(text))}
		response := createResponse([]byte("HTTP/1.1 200 OK"), headers, text)
		writeToConnection(conn, []byte(response))
	} else {
		response := createResponse([]byte("HTTP/1.1 404 Not Found"), nil, "")
		writeToConnection(conn, []byte(response))
	}
	conn.Close()
}