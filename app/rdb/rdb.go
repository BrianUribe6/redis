package rdb

import (
	"encoding/base64"
	"log"
	"os"
)

const FILENAME = "dump.rdb"

func CreateEmpty() {
	str := "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Panic(err)
	}

	os.WriteFile(FILENAME, data, 0666)
}
