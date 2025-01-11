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
	body        string
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
	fmt.Printf("Listening on port 4221\n")

	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConn(connection)
	}
}

func handleConn(connection net.Conn) {
	defer connection.Close()

	req, err := getRequest(connection)
	if err != nil {
		fmt.Println("error creating request object: ", err.Error())
		os.Exit(1)
	}

	switch req.method {
	case "GET":
		targetParts := strings.Split(req.target, "/")
		if targetParts[1] == "echo" {
			length := len(targetParts[2])
			connection.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", length, targetParts[2])))
		} else if req.target == "/user-agent" {
			length := len(req.header["user-agent"])
			connection.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", length, req.header["user-agent"])))
		} else if req.target == "/" {
			connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
		connection.Close()
	default:
		connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		connection.Close()
	}
}

func getRequest(conn net.Conn) (Request, error) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)

	if err != nil {
		return Request{}, fmt.Errorf("couldn`t read from connection: %w", err)
	}

	request := strings.Split(string(buffer), "\r\n\r\n")
	reqLineString := strings.Split(request[0], "\r\n")[0]
	reqLine := strings.Split(reqLineString, " ")
	header := strings.Split(request[0], "\r\n")[1:]

	fmt.Printf("reqLine: %v \n", header)
	body := request[1]

	fmt.Printf("body: %v \n", body)

	headerValues := make(map[string]string)
	for i := 0; i < len(header); i++ {
		v := strings.Split(header[i], ": ")
		headerValues[strings.ToLower(v[0])] = v[1]
	}

	return Request{
		method:      reqLine[0],
		target:      strings.ToLower(reqLine[1]),
		httpVersion: reqLine[2],
		header:      headerValues,
		body:        body,
	}, nil
}
