package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Listen on a port 8080
	listener , err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening", err)
		return
	}
	defer listener.Close()	// Close the listener when the main function ends
	fmt.Println("Listening on port 8080...")

	// Wait for connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err)
			continue
		}
		// Handle each new connection in a separate goroutine 
		go handleConnection(conn)
	}
}


func handleConnection(conn net.Conn) {
	defer conn.Close()	// Close the connection when the function ends

	// Create a new reader for the connection
	reader := bufio.NewReader(conn)

	// Read the request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request line", err)
		return
	}
	requestLine = strings.TrimSpace(requestLine)	// Remove leading and trailing whitespaces

	// Parse the request line
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		// Invalid format of the request line
		fmt.Println("Invalid request line", requestLine)
		return
	}
	method := parts[0]	// "GET" etc.
	path := parts[1]	// "/helloworld.html" etc.
	// httpVersion := parts[2]

	// ========== Read the headers ==========
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading header:", err)
			return
		}
		// Finished reading the headers when we encounter an empty line
		if strings.TrimSpace(line) == "" {
			break
		}
	}

	// Check if the method is GET (we only support GET requests now)
	if method != "GET" {
		fmt.Fprintf(conn, "HTTP/1.1 405 Method Not Allowed\r\n\r\n")
		return
	}

	// Default to "helloworld.html" if the path is "/"
	if path == "/" {
		path = "/helloworld.html"
	}

	// Remove the leading "/" from the path to get the file name
	if strings.HasPrefix(path, "/") {
		// Remove the leading "/"
		path = path[1:]
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		sendNotFound(conn)
		return
	}
	defer file.Close()

	// Get the file size to set the Content-Length header
	stat, _ := file.Stat()
	fileSize := stat.Size()

	// ========== Send the response ==========
	// Status line
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	// Headers
	fmt.Fprintf(conn, "Content-Type: text/html\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", fileSize)
	// End of headers
	fmt.Fprint(conn, "\r\n")

	// Send the file
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if n == 0 || err != nil {
			break
		}
		conn.Write(buffer[:n])
	}
}


// Send a 404 Not Found response
func sendNotFound(conn net.Conn) {
	notfoundFile := "notfound.html"
	file, err := os.Open(notfoundFile)
	if err != nil {
		fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
		return
	}
	defer file.Close()

	// Get the file size to set the Content-Length header
	stat, _ := file.Stat()
	fileSize := stat.Size()

	fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n")
	fmt.Fprintf(conn, "Content-Type: text/html\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", fileSize)
	fmt.Fprintf(conn, "\r\n")

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if n == 0 || err != nil {
			break
		}
		conn.Write(buffer[:n])
	}
}
