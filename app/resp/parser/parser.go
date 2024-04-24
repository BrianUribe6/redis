package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type CommandParser struct {
	Reader io.ByteReader
}

type Command struct {
	Label string
	Args  []string
}

var errSyntax error = errors.New("syntax error")

const (
	ARRAY_TYPE         = '*'
	BULK_STRING_TYPE   = '$'
	SIMPLE_STRING_TYPE = '+'
	CR                 = '\r'
	LF                 = '\n'
)

func NewCommandParser(buffer []byte) CommandParser {
	return CommandParser{
		Reader: bytes.NewReader(buffer),
	}
}

func FromReader(r io.ByteReader) CommandParser {
	return CommandParser{
		Reader: r,
	}
}

func (p *CommandParser) Parse() (*Command, error) {
	token, _ := p.Reader.ReadByte()
	if token != ARRAY_TYPE {
		return nil, errSyntax
	}
	arrLength, err := p.ParseNumber()
	if err != nil {
		return nil, err
	}

	args := make([]string, 0, arrLength)
	for i := 0; i < arrLength; i++ {
		s, err := p.ParseBulkString()
		if err != nil {
			return nil, err
		}
		//commands are case-insensitive
		args = append(args, strings.ToLower(s))
	}
	cmd := &Command{args[0], args[1:]}

	return cmd, nil
}

func (p *CommandParser) ParseSimpleString() (string, error) {
	token, err := p.Reader.ReadByte()
	if err != nil {
		return "", err
	}
	if token != SIMPLE_STRING_TYPE {
		return "", syntaxError(SIMPLE_STRING_TYPE, rune(token))
	}
	s, err := p.readUntilCRLF()

	return string(s), err

}

func (p *CommandParser) ParseBulkString() (string, error) {
	token, err := p.Reader.ReadByte()
	if err != nil {
		return "", err
	}
	if token != BULK_STRING_TYPE {
		return "", syntaxError(BULK_STRING_TYPE, rune(token))
	}
	length, err := p.ParseNumber()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for i := 0; i < length; i++ {
		token, err := p.Reader.ReadByte()
		if err != nil {
			return sb.String(), fmt.Errorf("expected to read %d bytes, only got %d", length, i)
		}
		sb.WriteByte(token)
	}
	return sb.String(), p.readCRLF()
}

func (p *CommandParser) ParseNumber() (int, error) {
	rawNumber, err := p.readUntilCRLF()
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(rawNumber))
}

func (p *CommandParser) readUntilCRLF() ([]byte, error) {
	var result []byte
	for {
		b, err := p.Reader.ReadByte()
		if err == io.EOF {
			return result, err
		}
		result = append(result, b)

		last := len(result) - 2
		if last >= 0 && result[last] == CR && b == LF {
			return result[:last], nil
		}
	}
}

func (p *CommandParser) readCRLF() error {
	token, err := p.Reader.ReadByte()
	if err != nil || token != CR {
		return syntaxError('\r', rune(token))
	}
	token, err = p.Reader.ReadByte()
	if err != nil || token != LF {
		return syntaxError('\n', rune(token))
	}
	return nil
}

func syntaxError(expected rune, got rune) error {
	return fmt.Errorf("syntax error: expected '%c', got '%c ", expected, got)
}
