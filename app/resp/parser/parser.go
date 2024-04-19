package parser

import (
	"errors"
	"strings"
)

type CommandParser struct {
	buffer []byte
	pos    int
}

type Command struct {
	Label string
	Args  []string
}

var EOF error = errors.New("reached end of file")
var errSyntax error = errors.New("syntax error")

const (
	ARRAY_TYPE         = '*'
	BULK_STRING_TYPE   = '$'
	SIMPLE_STRING_TYPE = '+'
	CR                 = '\r'
	LF                 = '\n'
)

func NewCommandParser(buffer []byte) CommandParser {
	parser := CommandParser{
		buffer: buffer,
		pos:    0,
	}
	return parser
}

func (p *CommandParser) Peek() byte {
	return p.buffer[p.pos]
}

func (p *CommandParser) Next() (byte, error) {
	if p.pos < len(p.buffer) {
		val := p.buffer[p.pos]
		p.pos++
		return val, nil
	}
	return 0, EOF
}

func (p *CommandParser) Position() int {
	return p.pos
}

func (p *CommandParser) Parse() (*Command, error) {
	token, _ := p.Next()
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
		args = append(args, s)
	}
	cmd := &Command{args[0], args[1:]}

	return cmd, nil
}

func (p *CommandParser) ParseSimpleString() (string, error) {
	token, _ := p.Next()
	if token != SIMPLE_STRING_TYPE {
		return "", errSyntax
	}

	var sb strings.Builder
	for p.Peek() != CR {
		token, err := p.Next()
		if err != nil {
			return "", err
		}
		sb.WriteByte(token)
	}
	err := p.parseCRLF()
	return sb.String(), err
}

func (p *CommandParser) ParseBulkString() (string, error) {
	if rune(p.Peek()) != BULK_STRING_TYPE {
		return "", errors.New("a bulk string must start with $")
	}
	p.Next()
	length, err := p.ParseNumber()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for i := 0; i < length; i++ {
		token, err := p.Next()
		if err != nil {
			return "", errSyntax
		}
		sb.WriteByte(token)
	}
	err = p.parseCRLF()
	return sb.String(), err
}

func (p *CommandParser) parseCRLF() error {
	token, _ := p.Next()
	if rune(token) != CR {
		return errSyntax
	}
	token, _ = p.Next()
	if token != LF {
		return errSyntax
	}
	return nil
}

func (p *CommandParser) ParseNumber() (int, error) {
	arrLength := 0
	var token byte
	var err error
	for token, err = p.Next(); isDigit(token); token, err = p.Next() {
		if err != nil {
			return 0, errSyntax
		}
		arrLength = arrLength*10 + int(token) - '0'
	}

	if next, _ := p.Next(); token != CR || next != LF {
		return 0, errSyntax
	}
	return arrLength, nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
