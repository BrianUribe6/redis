package command

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

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
	rdb.CreateEmpty()

	rdbFile, err := os.Open(rdb.FILENAME)
	if err != nil {
		log.Fatal("Failed to open RDB:", err)
	}
	defer rdbFile.Close()

	rdbFileInfo, err := rdbFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(rdbFile)

	// 3. Format it as a RESP file syntax and send it in CHUNKS
	// RESP Syntax for sending files is $<length_of_file>\r\n<contents_of_file>
	file := []byte(fmt.Sprintf("$%d\r\n", rdbFileInfo.Size()))

	data := make([]byte, 4096)
	for {
		n, err := reader.Read(data)
		if err == io.EOF {
			break
		}
		file = append(file, data[:n]...)
	}
	con.Write(file)

	if err != nil {
		log.Println("Sync failed:", err.Error())
		return
	}

	log.Printf("Syncronization with replica %s succeeded", con.RemoteAddr().String())
}
