package parser_test

import (
	"bytes"
	"testing"

	parser "github.com/codecrafters-io/redis-starter-go/app/resp/parser"
)

type TestInput struct {
	input       string
	expected    any
	shouldError bool
}

type UnitTests map[string]TestInput

func TestParseNumber(t *testing.T) {
	tests := UnitTests{
		"Parsing a valid number": {
			input:       "556\r\n",
			expected:    556,
			shouldError: false,
		},
		"Parsing a number with invalid characters": {
			input:       "55a6\r\n",
			expected:    0,
			shouldError: true,
		},
		"Parsing a number without a CR": {
			input:       "1234\t\n",
			expected:    0,
			shouldError: true,
		},
		"Parsing a number without a LF": {
			input:       "1234\r\t",
			expected:    0,
			shouldError: true,
		},
		"Parsing a number with CRLF sequence in the middle": {
			// The next command should error after this input, but parseNumber
			// should gracefully parse up to CRLF
			input:       "12\r\n34",
			expected:    12,
			shouldError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(test.input))
			parser := parser.New(reader)
			value, err := parser.ParseNumber()
			if test.shouldError {
				if err == nil {
					t.Fatalf("Parsing number with input '%s', expected error, got %d", test.input, value)
				}
			} else if value != test.expected {
				t.Fatalf("Parsing number with input '%s', expected %d, got %d", test.input, test.expected, value)
			}
		})
	}
}

func TestParseBulkString(t *testing.T) {
	tests := UnitTests{
		"Parsing a valid string": {
			input:       "$5\r\nhello\r\n",
			expected:    "hello",
			shouldError: false,
		},
		"Parsing a valid string with mismatched length": {
			input:       "$3\r\nhello\r\n",
			expected:    "",
			shouldError: true,
		},
		"Parse empty string": {
			input:       "$0\r\n\r\n",
			expected:    "",
			shouldError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(test.input))
			parser := parser.New(reader)
			value, err := parser.ParseBulkString()
			if test.shouldError {
				if err == nil {
					t.Fatalf("parsing bulk string with input '%s', expected error, got %s", test.input, value)
				}
			} else if value != test.expected {
				t.Fatalf("parsing bulk string with input '%s', expected %s, got %s", test.input, test.expected, value)
			}
		})
	}
}

func TestParseSimpleString(t *testing.T) {
	tests := UnitTests{
		"Parsing a valid string": {
			input:       "+hello\r\n",
			expected:    "hello",
			shouldError: false,
		},
		"Parse empty string": {
			input:       "+\r\n",
			expected:    "",
			shouldError: false,
		},
		"Missing type marker": {
			input:       "hello\r\n",
			expected:    "",
			shouldError: true,
		},
		"Multiple CRLF": {
			input:       "+hello\r\nworld\r\n",
			expected:    "hello",
			shouldError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(test.input))
			parser := parser.New(reader)
			value, err := parser.ParseSimpleString()
			if test.shouldError {
				if err == nil {
					t.Fatalf("%s: '%s', expected error, got %s", name, test.input, value)
				}
			} else if value != test.expected {
				t.Fatalf("%s: '%s', expected %s, got %s", name, test.input, test.expected, value)
			}
		})
	}
}
