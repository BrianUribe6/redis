package command

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

const (
	RDB_FILENAME         = "dump.rdb"
	errWrongNumberOfArgs = "wrong number of arguments"
	errSyntax            = "syntax error"
)

type Executor interface {
	Execute(conn net.Conn)
}

type Command struct {
	label string
	args  []string
}

type PingCommand Command
type EchoCommand Command
type SetCommand Command
type GetCommand Command
type InfoCommand Command
type ReplConfCommand Command
type PSYNCCommand Command
type NotImplementedCommand Command

func New(label string, params []string) Executor {
	switch strings.ToLower(label) {
	case "ping":
		return &PingCommand{label, params}
	case "echo":
		return &EchoCommand{label, params}
	case "set":
		return &SetCommand{label, params}
	case "get":
		return &GetCommand{label, params}
	case "info":
		return &InfoCommand{label, params}
	case "replconf":
		return &ReplConfCommand{label, params}
	case "psync":
		return &PSYNCCommand{label, params}
	}
	return &NotImplementedCommand{}
}

func (cmd *PingCommand) Execute(con net.Conn) {
	if len(cmd.args) == 0 {
		resp.ReplySimpleString(con, "PONG")
	} else if len(cmd.args) == 1 {
		resp.ReplyBulkString(con, cmd.args[0])
	} else {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
	}
}

func (cmd *EchoCommand) Execute(con net.Conn) {
	if len(cmd.args) != 1 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	resp.ReplyBulkString(con, cmd.args[0])
}

func (cmd *NotImplementedCommand) Execute(con net.Conn) {
	resp.ReplySimpleError(con, "unknown command, may not be implemented yet")
}

func (cmd *SetCommand) Execute(con net.Conn) {
	numArgs := len(cmd.args)
	// TODO write a proper flag parser
	if numArgs != 2 && numArgs != 4 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	key := cmd.args[0]
	value := cmd.args[1]
	var expiry int64 = -1

	if numArgs == 4 {
		pxFlag := strings.ToLower(cmd.args[2])
		if pxFlag != "px" {
			resp.ReplySimpleError(con, errSyntax)
			return
		}

		exp, err := strconv.ParseInt(cmd.args[3], 10, 64)
		if err != nil || exp < 0 {
			resp.ReplySimpleError(con, "invalid expiry time")
			return
		}
		expiry = exp
	}

	store.Set(key, value, expiry)
	resp.ReplySuccess(con)
}

func (cmd *GetCommand) Execute(con net.Conn) {
	if len(cmd.args) != 1 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	value, exist := store.Get(cmd.args[0])
	if !exist {
		resp.ReplyNullBulkString(con)
	} else {
		resp.ReplyBulkString(con, value)
	}
}

func (cmd *InfoCommand) Execute(con net.Conn) {
	if len(cmd.args) == 0 {
		resp.ReplySimpleError(con, "unsupported option. Currently only info replication is supported")
		return
	}
	if len(cmd.args) > 1 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}

	if cmd.args[0] != "replication" {
		resp.ReplySimpleError(con, errSyntax)
		return
	}

	resp.ReplyBulkString(con, store.Info.String())
}

func (cmd *ReplConfCommand) Execute(con net.Conn) {
	if len(cmd.args) < 2 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	// TODO
	switch cmd.args[0] {
	case "listening-port":
		break
	case "capa":
		break
	default:
		resp.ReplySimpleError(con, errSyntax)
		return
	}

	resp.ReplySuccess(con)
}

func (cmd *PSYNCCommand) Execute(con net.Conn) {
	if len(cmd.args) != 2 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	log.Println("Received synchronization request from", con.RemoteAddr().String())
	// 1. Notify replica that it should expect a full copy of the database
	response := fmt.Sprintf("FULLRESYNC %s 0", store.Info.MasterReplId)
	resp.ReplySimpleString(con, response)

	// 2. Read the file dump of the database
	log.Println("Loading RDB...")
	rdbFile, err := os.Open(RDB_FILENAME)
	if err != nil {
		log.Println(err)
		return
	}
	defer rdbFile.Close()

	fileInfo, err := rdbFile.Stat()
	if err != nil {
		log.Println("unable to gather info from RDB")
		return
	}
	// 3. Format it as a RESP file syntax and send it in CHUNKS
	// RESP Syntax for sending files is $<length_of_file>\r\n<contents_of_file>
	con.Write([]byte(fmt.Sprintf("$%d\r\n", fileInfo.Size())))

	buf := make([]byte, 1024)
	for {
		n, err := rdbFile.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			break
		}
		if err == io.EOF {
			break
		}

		con.Write(buf[:n])
	}

	log.Printf("Syncronization with replica %s succeeded", con.RemoteAddr().String())
}
