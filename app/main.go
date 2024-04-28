package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/server"
)

var portNumFlag = flag.Int("port", 6379, "the port at which the server will be listening to")
var masterHostname = flag.String("replicaof", "localhost", "the address of this server master")

func main() {
	flag.Parse()
	s := server.New("0.0.0.0", *portNumFlag)

	if isReplica() {
		masterPort := flag.Arg(0)
		_, err := strconv.Atoi(masterPort)
		if err != nil {
			log.Fatal("Invalid port number")
		}
		s.ListenAsReplica(*masterHostname, masterPort)
	}
	s.Listen()

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
