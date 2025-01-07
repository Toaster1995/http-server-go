package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

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

		buffer := make([]byte, 1024)
		_, err = connection.Read(buffer)

		if err != nil {
			fmt.Println("couldn`t read from connection")
			connection.Close()
			continue
		}

		req := strings.Split(string(buffer), "\r\n")[0]
		target := strings.Split(req, " ")[1]

		if target != "/" {
			connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		} else {
			connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		}
		connection.Close()
	}
}
