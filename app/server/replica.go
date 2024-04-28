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

	go s.handleClient(master)
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
		c.SendArrayBulk(cmd...)
		// We want to know master's response immediately instead of buffering it
		c.Flush()
		_, err = c.Read(buffer)
		if err != nil {
			log.Panic("handshake failed, did not get a response from master")
		}
	}
	log.Println("Handshake successful. Waiting for full resync...")
	c.SendArrayBulk("psync", "?", "-1")
	c.Flush()

	_, err = handleFullResyncResponse(c.Reader)
	if err != nil {
		log.Panic(err)
	}

	// After replying with a fullresync the master should be sending
	// an RDB file with the full database contents
	_, err = parseRDBFile(c.Reader)
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
func handleFullResyncResponse(reader io.ByteReader) (string, error) {
	p := parser.New(reader)
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

func parseRDBFile(reader io.ByteReader) ([]byte, error) {
	// Expect master to respond with $<file_size>\r\n<file_contents>
	content := bytes.NewBuffer([]byte{})
	p := parser.New(reader)
	token, err := p.ReadByte()
	if err != nil || token != parser.BULK_STRING_TYPE {
		return nil, errUnexpected
	}

	fileSize, err := p.ParseNumber()
	if err != nil {
		return nil, errUnexpected
	}

	for i := 0; i < fileSize; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("expected a file with size %d bytes, got %d instead", fileSize, i)
		}

		content.WriteByte(b)
	}
	return content.Bytes(), nil
}
