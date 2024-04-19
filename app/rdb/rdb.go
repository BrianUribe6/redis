package rdb

import (
	"encoding/base64"
	"io"
	"log"
	"os"
)

const FILENAME = "dump.rdb"

type Reader struct {
	file *os.File
	Info os.FileInfo
}

func New() (*Reader, error) {
	CreateEmpty() // Always create an empty RDB file (for now)
	rdbFile, err := os.Open(FILENAME)
	if err != nil {
		return nil, err
	}
	fileInfo, err := rdbFile.Stat()
	if err != nil {
		return nil, err
	}

	return &Reader{rdbFile, fileInfo}, nil
}

func (r *Reader) Read(callback func(buffer []byte)) error {
	buf := make([]byte, 1024)
	for {
		n, err := r.file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		callback(buf[:n])
	}
	return nil
}

func CreateEmpty() {
	str := "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Panic(err)
	}

	file, err := os.Create(FILENAME)
	if err != nil {
		log.Panic(err)
	}
	file.Write(data)
}

func (r *Reader) Close() {
	r.file.Close()
}
