package src

import (
	"bytes"
	"encoding/binary"
	"log"
	"fmt"
)

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func TreatError(s string, err error) error {
	return fmt.Errorf(s + ": %s", err)
}
