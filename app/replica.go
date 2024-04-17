package main

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
	errMissingPort   = "You must provide the master's port number"
	errInvalidPort   = "Invalid port numbers"
	errHandshakeFail = "Failed to establish handshake with master"
)

func configureReplica(masterHostname string, masterPort string) {
	log.Println("Initiating server in replica mode")
	store.Info.SetRole(store.SLAVE_ROLE)
	if len(masterPort) == 0 {
		log.Panic(errMissingPort)
	}
	_, err := strconv.Atoi(masterPort)
	if err != nil {
		log.Panic(errInvalidPort)
	}

	handshake(masterHostname + ":" + masterPort)
}

func handshake(masterAddress string) {
	con, err := net.Dial("tcp", masterAddress)
	log.Printf("Attempting to connect to master at %s\n", masterAddress)
	if err != nil {
		log.Panic("failed to connect" + err.Error())
	}
	defer con.Close()

	commands := [][]string{
		{"PING"},
		{"REPLCONF", "listening-port", fmt.Sprint(*portNumFlag)},
		{"REPLCONF", "capa", "psync2"},
		{"PSYNC", "?", "-1"},
	}

	log.Println("Initiating handshake")
	buffer := make([]byte, 1024)
	var n int
	for _, cmd := range commands {
		resp.ReplyArrayBulk(con, cmd)
		n, err = con.Read(buffer)
		if err != nil {
			log.Panic(errHandshakeFail)
		}
	}
	log.Println("Handshake successful. Waiting for full resync...")
	assertFullResync(buffer[:n])
	resyncWithMaster(con)

	log.Println("Full resync done.")
}

// Asserts buffer contens are formatted as a valid RESP FULLRESYNC command
//
// On success, returns the master's replication id
func assertFullResync(buffer []byte) string {
	message := strings.Split(string(buffer), " ")
	log.Println("PSYNC replied", message)
	// expect +FULLRESYNC <master_replid> 0 encoded as a SIMPLE string
	if len(message) != 3 {
		log.Fatal(errUnexpected)
	}
	if strings.ToLower(message[0]) != "+fullresync" {
		log.Fatal(errUnexpected)
	}
	if message[2] != "0\r\n" {
		log.Fatal(errUnexpected)
	}

	return message[1]
}

func resyncWithMaster(con net.Conn) {
	// Expect master to respond with $<file_size>\r\n<file_contents>
	buf := make([]byte, 1024)
	_, err := con.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	p := parser.NewCommandParser(buf)
	token, err := p.Next()
	if token != parser.BULK_STRING_TYPE || err != nil {
		log.Fatal(errUnexpected)
	}

	_, err = p.ParseNumber()
	if err != nil {
		log.Fatal(errUnexpected)
	}

}
