package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// go run client.go [METHOD] [HOST:PORT] [PATH]
// e.g. go run client.go GET localhost:8080 /helloworld.html
func main() {
	if len(os.Args) < 3 {
		fmt.Println(("Usage: go run client.go [HOST:PORT] [PATH]"))
		return
	}
	method := os.Args[1]
	hostPort := os.Args[2]
	path := os.Args[3]

	// Connect to the server
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		fmt.Println("Error connecting", err)
		return
	}
	defer conn.Close()	// Close the connection when the main function ends

	// Send HTTP GET request
	request := fmt.Sprintf("%s %s HTTP/1.1\r\nHOST: %s\r\n\r\n", method, path, hostPort)
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Error writing request", err)
		return
	}

	// Read the response and print it to the console
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(line)
	}
}
