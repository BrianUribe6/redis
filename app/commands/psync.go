package command

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type PSYNCCommand Command

func (cmd *PSYNCCommand) Execute(c client.Client) RESPValue {
	if len(cmd.args) != 2 {
		return resp.EncodeSimpleError(errWrongNumberOfArgs)
	}
	conn := c.Connection()
	log.Println("Received synchronization request from", conn.RemoteAddr().String())

	// 1. Notify replica that it should expect a full copy of the database
	c.Write(resp.EncodeSimpleString(fmt.Sprintf("FULLRESYNC %s 0", store.Info.MasterReplId)))

	// 2. Read the file dump of the database
	log.Println("Loading RDB...")

	data := CreateEmptyRDB()
	// 3. Format it as a RESP file syntax and send it in CHUNKS
	// RESP Syntax for sending files is $<length_of_file>\r\n<contents_of_file>
	file := []byte(fmt.Sprintf("$%d\r\n", len(data)))

	file = append(file, data...)

	log.Printf("Syncronization with replica %s succeeded", conn.RemoteAddr().String())

	return file
}

func CreateEmptyRDB() []byte {
	str := "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Panic(err)
	}

	// os.WriteFile(FILENAME, data, 0666)
	return data
}
