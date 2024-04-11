package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

var portNumFlag = flag.Int("port", 6379, "the port at which the server will be listening to")
var replicaOfFlag = flag.String("replicaof", "localhost", "the address of this server master")

func main() {
	flag.Parse()

	if isReplica() {
		replicaPort := flag.Arg(0)
		if len(replicaPort) == 0 {
			fmt.Println("you must provide the master's port number")
			return
		}
		_, err := strconv.Atoi(replicaPort)
		if err != nil {
			fmt.Println("invalid port number")
			return
		}
		store.Info.SetRole(store.SLAVE_ROLE)
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

	fmt.Printf("Server running on %s", address)
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

func isReplica() bool {
	isSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "replicaof" {
			isSet = true
		}
	})

	return isSet
}
