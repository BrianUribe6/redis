package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/parser"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

const (
	errInvalidAck    = "Invalid acknowledgement, expected FULLRESYNC"
	errUnexpected    = "Received unexpected response from master"
	errHandshakeFail = "Failed to establish handshake with master"
)

func configureReplica(listeningPort int, masterAddress string) {
	log.Println("Initiating server in replica mode")
	store.Info.SetRole(store.SLAVE_ROLE)

	client := handshake(listeningPort, masterAddress)

	log.Println("Listening for master commands")
	listenMasterCommands(client)
}

func handshake(listeningPort int, masterAddress string) net.Conn {
	con, err := net.Dial("tcp", masterAddress)
	log.Printf("Attempting to connect to master at %s\n", masterAddress)
	if err != nil {
		log.Fatal("failed to connect" + err.Error())
	}

	commands := [][]string{
		{"ping"},
		{"replconf", "listening-port", fmt.Sprint(listeningPort)},
		{"replconf", "capa", "psync2"},
		{"psync", "?", "-1"},
	}

	log.Println("Initiating handshake")
	buffer := make([]byte, 1024)
	var n int
	for _, cmd := range commands {
		resp.ReplyArrayBulk(con, cmd)
		n, err = con.Read(buffer)
		if err != nil {
			log.Fatal(errHandshakeFail)
		}
	}
	log.Println("Handshake successful. Waiting for full resync...")

	assertFullResyncReceived(buffer[:n])
	// After replying with a fullresync the master should be sending
	// an RDB file with the full database contents
	resyncWithMaster(con, buffer)
	log.Println("Full resync done.")

	return con
}

// Asserts buffer contents are formatted as valid FULLRESYNC command.
//
// expected format: +FULLRESYNC <master_replid> <offset>\r\n
//
// On success, returns the master's replication id
func assertFullResyncReceived(buffer []byte) string {
	p := parser.NewCommandParser(buffer)

	fmt.Println(string(buffer))

	s, err := p.ParseSimpleString()
	if err != nil {
		log.Fatal(errInvalidAck)
	}
	log.Println("PSYNC replied with:", s)

	message := strings.Split(s, " ")
	if len(message) != 3 {
		log.Fatal(errUnexpected)
	}
	label := message[0]
	// masterIReplId := message[1]
	if label != "FULLRESYNC" {
		log.Fatal("invalid command")
	}

	_, err = strconv.Atoi(message[2]) // offset

	if err != nil {
		log.Fatalf("invalid offset")
	}

	return message[1]
}

func resyncWithMaster(con net.Conn, buffer []byte) {
	// Expect master to respond with $<file_size>\r\n<file_contents>
	n, err := con.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	p := parser.NewCommandParser(buffer[:n])
	token, err := p.Next()
	if err != nil || token != parser.BULK_STRING_TYPE {
		log.Fatal(errUnexpected)
	}

	_, err = p.ParseNumber()
	if err != nil {
		log.Fatal(errUnexpected)
	}
	//TODO do something with the file
}

func listenMasterCommands(con net.Conn) {
	defer con.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := con.Read(buffer)
		if err != nil {
			log.Println("")
		}
		commandParser := parser.NewCommandParser(buffer[:n])
		c, err := commandParser.Parse()
		if err != nil {
			break
		}
		log.Println(c.Label, strings.Join(c.Args, " "))
	}
}
