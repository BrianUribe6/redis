package parser_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func TestParseNumber(t *testing.T) {
	tests := map[string]struct {
		input       []byte
		expected    int
		shouldError bool
	}{
		"Parsing a valid number": {
			input:       []byte("556\r\n"),
			expected:    556,
			shouldError: false,
		},
		"Parsing a number with invalid characters": {
			input:       []byte("55a6\r\n"),
			expected:    0,
			shouldError: true,
		},
		"Parsing a number without a CR": {
			input:       []byte("1234\t\n"),
			expected:    0,
			shouldError: true,
		},
		"Parsing a number without a LF": {
			input:       []byte("1234\r\t"),
			expected:    0,
			shouldError: true,
		},
		"Parsing a number with CRLF sequence in the middle": {
			// The next command should error after this input, but parseNumber
			// should gracefully parse up to CRLF
			input:       []byte("12\r\n34"),
			expected:    12,
			shouldError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			parser := parser.New(test.input)
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
	tests := map[string]struct {
		input       []byte
		expected    string
		shouldError bool
	}{
		"Parsing a valid string": {
			input:       []byte("$5\r\nhello\r\n"),
			expected:    "hello",
			shouldError: false,
		},
		"Parsing a valid string with mismatched length": {
			input:       []byte("$3\r\nhello\r\n"),
			expected:    "",
			shouldError: true,
		},
		"Parse empty string": {
			input:       []byte("$0\r\n\r\n"),
			expected:    "",
			shouldError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			parser := parser.New(test.input)
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
