package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	command "github.com/codecrafters-io/redis-starter-go/app/commands"
	parser "github.com/codecrafters-io/redis-starter-go/app/resp/parser"
)

var portNumFlag = flag.Int("port", 6379, "the port at which the server will be listening to")
var masterHostname = flag.String("replicaof", "localhost", "the address of this server master")

func main() {
	flag.Parse()

	if isReplica() {
		masterPort := flag.Arg(0)
		configureReplica(*masterHostname, masterPort)
	}

	address := fmt.Sprint("0.0.0.0:", *portNumFlag)

	startServer(address)
}

func startServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	log.Printf("Server running at %s\n", address)
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
		n, err := conn.Read(buff)
		if err != nil {
			break
		}
		commandParser := parser.NewCommandParser(buff[:n])
		c, err := commandParser.Parse()
		if err != nil {
			break
		}
		command.New(c.Label, c.Args).Execute(conn)
	}
}

func isReplica() bool {
	isSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "replicaof" {
			isSet = true
		}
	})

	return isSet
}
