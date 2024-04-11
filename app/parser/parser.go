package parser

import (
	"errors"
	"fmt"
	"strings"

	command "github.com/codecrafters-io/redis-starter-go/app/commands"
)

type CommandParser struct {
	buffer []byte
	pos    int
}

var EOF error = errors.New("reached end of file")
var errSyntax error = errors.New("syntax error")

const (
	ARRAY_TYPE       = '*'
	BULK_STRING_TYPE = '$'
	CR               = '\r'
	LF               = '\n'
)

func New(buffer []byte) CommandParser {
	parser := CommandParser{
		buffer: buffer,
		pos:    0,
	}
	return parser
}

func (p *CommandParser) peek() byte {
	return p.buffer[p.pos]
}

func (p *CommandParser) next() (byte, error) {
	if p.pos < len(p.buffer) {
		val := p.buffer[p.pos]
		p.pos++
		return val, nil
	}
	return 0, EOF
}

func (p *CommandParser) Parse() (*command.Executor, error) {
	token, _ := p.next()
	if token != ARRAY_TYPE {
		return nil, fmt.Errorf("expected '%c' got %c instead", ARRAY_TYPE, token)
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
		args = append(args, s)
	}

	c := command.New(args[0], args[1:])

	return &c, nil
}

func (p *CommandParser) ParseBulkString() (string, error) {
	if rune(p.peek()) != BULK_STRING_TYPE {
		return "", errors.New("a bulk string must start with $")
	}
	p.next()
	length, err := p.ParseNumber()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for i := 0; i < length; i++ {
		token, err := p.next()
		if err != nil {
			return "", errSyntax
		}
		sb.WriteByte(token)
	}
	err = p.parseCRLF()
	return sb.String(), err
}

func (p *CommandParser) parseCRLF() error {
	token, _ := p.next()
	if rune(token) != CR {
		return errSyntax
	}
	token, _ = p.next()
	if token != LF {
		return errSyntax
	}
	return nil
}

func (p *CommandParser) ParseNumber() (int, error) {
	arrLength := 0
	var token byte
	var err error
	for token, err = p.next(); isDigit(token); {
		if err != nil {
			return 0, errSyntax
		}
		arrLength = arrLength*10 + int(token) - '0'
		token, err = p.next()
	}

	if next, _ := p.next(); token != CR || next != LF {
		return 0, errSyntax
	}
	return arrLength, nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
