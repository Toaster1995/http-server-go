package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Request struct {
	method      string
	target      string
	httpVersion string
	header      map[string]string
	body        []string
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	handleConn(l)
}

func handleConn(l net.Listener) {
	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		req, err := getRequest(connection)
		if err != nil {
			fmt.Println("error creating request object: ", err.Error())
			os.Exit(1)
		}

		targetParts := strings.Split(req.target, "/")
		if targetParts[1] == "echo" {
			length := len(targetParts[2])
			fmt.Printf("Target: %v", targetParts[2])
			connection.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", length, targetParts[2])))
		} else if req.target == "/" {
			connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
		connection.Close()
	}
}

func getRequest(conn net.Conn) (Request, error) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)

	if err != nil {
		return Request{}, fmt.Errorf("couldn`t read from connection: %w", err)
	}

	request := strings.Split(string(buffer), "\r\n")
	reqLine := request[0]
	headerBody := strings.Split(request[1], "\r\n\r\n")
	reqParts := strings.Split(reqLine, " ")

	headerVars := strings.Split(headerBody[0], "\r\n")
	headerValues := make(map[string]string)
	for _, item := range headerVars {
		v := strings.Split(item, ": ")
		headerValues[v[0]] = v[1]
	}

	return Request{
		method:      reqParts[0],
		target:      reqParts[1],
		httpVersion: reqParts[2],
		header:      headerValues,
	}, nil
}
