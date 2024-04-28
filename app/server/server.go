package server

import (
	"fmt"
	"io"
	"log"
	"net"

	command "github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	"github.com/codecrafters-io/redis-starter-go/app/resp/parser"
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
	address := net.JoinHostPort(s.hostname, fmt.Sprint(s.port))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	defer listener.Close()

	log.Printf("Server running at %s\n", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go s.handleClient(client.New(conn))
	}
}

func (s *Server) handleClient(cli client.Client) {
	defer cli.Close()
	for {
		p := parser.New(cli)
		decoded, err := p.Parse()
		if err != nil {
			if err == io.EOF {
				log.Printf("lost connection with client %s", cli.Connection().RemoteAddr())
				break
			}
			log.Println("Unrecognized command", err)
			continue
		}
		cmd := command.New(decoded.Label, decoded.Args)
		cmd.Execute(cli)
		cli.Flush()

		if decoded.Label == "psync" {
			s.subscribe(cli)
		}
		//TODO check if the command is mutable only "set" for now...
		if decoded.Label == "set" {
			s.propagate(decoded)
		}
	}
}

// Subscribe a replica to receive commands from master
func (s *Server) subscribe(c client.Client) {
	s.replicas = append(s.replicas, c)
}

// Send cmd to all connected replicas
func (s *Server) propagate(cmd *parser.Command) {
	for _, replica := range s.replicas {
		c := append([]string{cmd.Label}, cmd.Args...)
		replica.SendArrayBulk(c...)
		replica.Flush()
	}
}
