package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

// IntToHex converts an int64 to a slice of bytes
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)

	// Big Endian = Lower -> Higher
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DBExists checks if a Database already exists
func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
