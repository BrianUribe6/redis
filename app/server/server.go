package server

import (
	"fmt"
	"log"
	"net"

	command "github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	parser "github.com/codecrafters-io/redis-starter-go/app/resp/parser"
)

type Server struct {
	hostname string
	port     int
	replicas []client.Client
}

func New(hostname string, port int) *Server {
	return &Server{
		hostname: hostname,
		port:     port,
		replicas: []client.Client{},
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
		c := client.New(conn)
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		go s.handleClient(c)
	}
}

func (s *Server) handleClient(cli client.Client) {
	buff := make([]byte, 1024)
	isReplica := false
	for {
		n, err := cli.Read(buff)
		if err != nil {
			break
		}
		commandParser := parser.NewCommandParser(buff[:n])
		parsedCmd, err := commandParser.Parse()
		if err != nil {
			log.Println(err)
			break
		}
		if parsedCmd.Label == "psync" {
			isReplica = true
			s.subscribe(cli)
		}
		cmd := command.New(parsedCmd.Label, parsedCmd.Args)

		//TODO check if the command is mutable only "set" for now...
		if parsedCmd.Label == "set" {
			s.notify(buff[:n])
		}
		cmd.Execute(cli)

		cli.Flush()
	}
	if !isReplica {
		cli.Close()
	}
}

// Subscribe a replica to receive commands from master
func (s *Server) subscribe(c client.Client) {
	s.replicas = append(s.replicas, c)
}

// Send cmd to all connected replicas
func (s *Server) notify(cmd []byte) {
	for _, replica := range s.replicas {
		replica.Write(cmd)
		replica.Flush()
	}
}

func (s *Server) AsReplica(masterHostname string, masterPort int) {
	masterAddr := fmt.Sprintf("%s:%d", masterHostname, masterPort)
	configureReplica(s.port, masterAddr)
}
