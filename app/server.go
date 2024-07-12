package main

import (
	"bufio"
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

// @Citation: https://app.codecrafters.io/users/Powerisinschool
func parseRequest(scanner *bufio.Scanner) (*HTTPRequest, error) {
	var req HTTPRequest = HTTPRequest{}
	req.Headers = make(map[string]string)
	for i := 0; scanner.Scan(); i++ {
		if i == 0 {
			req.Method = strings.Split(scanner.Text(), " ")[0]
			req.Path = strings.Split(scanner.Text(), " ")[1]
			continue
		}
		headers := strings.Split(scanner.Text(), ": ")
		if len(headers) < 2 {
			req.Body = headers[0]
			break
		}
		req.Headers[headers[0]] = headers[1]
	}
	return &req, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	req, err := parseRequest(scanner)

	if err != nil {
		fmt.Fprintln(conn, "reading standard input:", err)
		os.Exit(1)
	}

	fmt.Println("Request: ", req.Method, req.Path, req.Headers, req.Body)

	switch path := req.Path; {
	case strings.HasPrefix(path, "/files"):
		fileName := strings.TrimPrefix(path, "/files/")
		dir := os.Args[2]
		data, err := os.ReadFile(dir + fileName)
		if err != nil {
			writeToConnection(conn, []byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			return
		}
		response := createResponse("HTTP/1.1 200 OK", []string{"Content-Type: application/octet-stream", "Content-Length: " + fmt.Sprint(len(data))}, string(data))
		writeToConnection(conn, []byte(response))
	case strings.HasPrefix(path, "/echo"):
		content := strings.TrimPrefix(path, "/echo/")
		response := createResponse("HTTP/1.1 200 OK", []string{"Content-Type: text/plain", "Content-Length: " + fmt.Sprint(len(content))}, content)
		writeToConnection(conn, []byte(response))
	case path == "/user-agent":
		body := req.Headers["User-Agent"]
		response := createResponse("HTTP/1.1 200 OK", []string{"Content-Type: text/plain", "Content-Length: " + fmt.Sprint(len(body))}, body)
		writeToConnection(conn, []byte(response))
	case path == "/":
		fmt.Print("Root path requested")
		writeToConnection(conn, []byte("HTTP/1.1 200 OK\r\n\r\n"))
	default:
		writeToConnection(conn, []byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}