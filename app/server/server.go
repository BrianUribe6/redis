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
	replicas []net.Conn
}

func New(hostname string, port int) *Server {
	return &Server{
		hostname: hostname,
		port:     port,
		replicas: []net.Conn{},
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
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	buff := make([]byte, 1024)
	isReplica := false
	for {
		n, err := conn.Read(buff)
		if err != nil {
			break
		}
		commandParser := parser.NewCommandParser(buff[:n])
		c, err := commandParser.Parse()
		if err != nil {
			log.Println(err)
			break
		}
		if c.Label == "psync" {
			isReplica = true
			s.subscribe(conn)
		}
		cmd := command.New(c.Label, c.Args)

		//TODO check if the command is mutable only "set" for now...
		if c.Label == "set" {
			s.notify(buff[:n])
		}
		cmd.Execute(conn)
	}
	if !isReplica {
		conn.Close()
	}
}

// Subscribe a replica to receive commands from master
func (s *Server) subscribe(conn net.Conn) {
	s.replicas = append(s.replicas, conn)
}

// Send cmd to all connected replicas
func (s *Server) notify(cmd []byte) {
	for _, replica := range s.replicas {
		replica.Write(cmd)
	}
}

func (s *Server) AsReplica(masterHostname string, masterPort int) {
	masterAddr := fmt.Sprintf("%s:%d", masterHostname, masterPort)
	configureReplica(s.port, masterAddr)
}
