package main

import (
	"bytes"
	"encoding/binary"
	"log"
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
