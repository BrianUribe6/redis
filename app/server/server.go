package server

import (
	"fmt"
	"log"
	"net"

	command "github.com/codecrafters-io/redis-starter-go/app/commands"
	parser "github.com/codecrafters-io/redis-starter-go/app/resp/parser"
)

type Server struct {
	hostname string
	port     int
}

func New(hostname string, port int) *Server {
	return &Server{
		hostname: hostname,
		port:     port,
	}
}

func (s *Server) Listen() {
	address := fmt.Sprintf("%s:%d", s.hostname, s.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	defer listener.Close()

	log.Printf("Server running at %s\n", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
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

func (s *Server) AsReplica(masterHostname string, masterPort int) {
	masterAddr := fmt.Sprintf("%s:%d", masterHostname, masterPort)
	configureReplica(s.port, masterAddr)
}
