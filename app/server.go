package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func main() {
	portNumFlag := flag.Int("port", 6379, "the port at which the server will be listening to")
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *portNumFlag)

	fmt.Printf("Server running at %s\n", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 1024)
	for {
		_, err := conn.Read(buff)
		if err != nil {
			break
		}
		commandParser := parser.New(buff)
		cmd, err := commandParser.Parse()
		if err != nil {
			break
		}
		(*cmd).Execute(conn)
	}
}
