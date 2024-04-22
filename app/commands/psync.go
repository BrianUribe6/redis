package command

import (
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type PSYNCCommand Command

func (cmd *PSYNCCommand) Execute(con net.Conn) {
	if len(cmd.args) != 2 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	log.Println("Received synchronization request from", con.RemoteAddr().String())
	// 1. Notify replica that it should expect a full copy of the database
	resp.ReplySimpleString(con, fmt.Sprintf("FULLRESYNC %s 0", store.Info.MasterReplId))

	// 2. Read the file dump of the database
	log.Println("Loading RDB...")
	reader, err := rdb.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer reader.Close()

	// 3. Format it as a RESP file syntax and send it in CHUNKS
	// RESP Syntax for sending files is $<length_of_file>\r\n<contents_of_file>
	file := []byte(fmt.Sprintf("$%d\r\n", reader.Info.Size()))
	err = reader.Read(func(buffer []byte) {

		// FIXME I'm sending the whole thing to satisfy codecrafter's unit test
		// but the right thing instead is to write it in chunks i.e con.Write(buffer)
		// (Imagine if the file was 16GB)
		file = append(file, buffer...)
	})

	con.Write(file)

	if err != nil {
		log.Println("Sync failed:", err.Error())
		return
	}

	log.Printf("Syncronization with replica %s succeeded", con.RemoteAddr().String())
}
