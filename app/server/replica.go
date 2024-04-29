package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	command "github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	"github.com/codecrafters-io/redis-starter-go/app/resp/parser"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

var (
	errInvalidAck = errors.New("invalid acknowledgement, expected FULLRESYNC")
	errUnexpected = errors.New("received unexpected response from master")
)

func (s *Server) ListenAsReplica(masterHostname string, masterPort string) {
	masterAddr := net.JoinHostPort(masterHostname, masterPort)

	log.Println("Initiating server in replica mode")
	store.Info.SetRole(store.SLAVE_ROLE)

	master := handshake(s.port, masterAddr)

	go handleMaster(master)
}

func handshake(listeningPort int, masterAddress string) client.Client {
	con, err := net.Dial("tcp", masterAddress)
	c := client.New(con)
	log.Printf("Attempting to connect to master at %s\n", masterAddress)
	if err != nil {
		log.Fatal("failed to connect" + err.Error())
	}

	commands := [][]string{
		{"ping"},
		{"replconf", "listening-port", fmt.Sprint(listeningPort)},
		{"replconf", "capa", "psync2"},
	}

	buffer := make([]byte, 1024)
	for _, cmd := range commands {
		c.Write(resp.EncodeArrayBulk(cmd...))
		// We want to know master's response immediately instead of buffering it
		c.Flush()
		_, err = c.Read(buffer)
		if err != nil {
			log.Panic("handshake failed, did not get a response from master")
		}
	}
	log.Println("Handshake successful. Waiting for full resync...")
	c.Write(resp.EncodeArrayBulk("psync", "?", "-1"))
	c.Flush()

	_, err = handleFullResyncResponse(c)
	if err != nil {
		log.Panic(err)
	}

	// After replying with a fullresync the master should be sending
	// an RDB file with the full database contents
	_, err = parseRDBFile(c)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Full resync done")

	return c
}

// Asserts buffer contents are formatted as a valid FULLRESYNC command.
//
// expected format: +FULLRESYNC <master_replid> <offset>\r\n
//
// On success, returns the master's replication id
func handleFullResyncResponse(c client.Client) (string, error) {
	p := parser.New(c)
	s, err := p.ParseSimpleString()
	if err != nil {
		return "", errInvalidAck
	}
	log.Println("PSYNC replied with:", s)

	message := strings.Split(s, " ")
	if len(message) != 3 {
		return "", errUnexpected
	}
	label := message[0]
	// masterIReplId := message[1]
	if label != "FULLRESYNC" {
		log.Fatal("expected to receive FULLRESYNC got: " + label)
	}

	_, err = strconv.Atoi(message[2]) // offset

	if err != nil {
		return "", errors.New("offset must be a valid integer")
	}

	return message[1], nil
}

func parseRDBFile(c client.Client) ([]byte, error) {
	// Expect master to respond with $<file_size>\r\n<file_contents>
	content := bytes.NewBuffer([]byte{})
	p := parser.New(c)
	token, err := p.Next()
	if err != nil || token != parser.BULK_STRING_TYPE {
		return nil, errUnexpected
	}

	fileSize, err := p.ParseNumber()
	if err != nil {
		return nil, errUnexpected
	}

	for i := 0; i < fileSize; i++ {
		b, err := c.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("expected a file with size %d bytes, got %d instead", fileSize, i)
		}

		content.WriteByte(b)
	}

	c.BytesRead += p.BytesRead() + fileSize
	return content.Bytes(), nil
}

func handleMaster(c client.Client) {
	defer c.Close()
	for {
		p := parser.New(c)
		decoded, err := p.Parse()
		if err != nil {
			if err == io.EOF {
				log.Printf("lost connection with client %s", c.Connection().RemoteAddr())
				break
			}
			log.Println("Unrecognized command", err)
			continue
		}
		cmd := command.New(decoded.Label, decoded.Args)
		response := cmd.Execute(c)
		if decoded.Label == "replconf" {
			c.Write(response)
			c.Flush()

		}
		c.BytesRead += p.BytesRead()
	}
}
