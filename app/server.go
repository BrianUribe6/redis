package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
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
		handleClient(conn);
	}
}


func handleClient(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Failed to read message:", err);
		os.Exit(1)
	}
	fmt.Println(string(buff[:n]));
	msg := []byte("+PONG\r\n")
	conn.Write(msg)

}