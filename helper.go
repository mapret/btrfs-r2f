package btrfsr2f

import (
	"encoding/binary"
	"io"
	"os"
)

func readAndPanic(reader io.Reader, data interface{}) {
	err := binary.Read(reader, binary.LittleEndian, data)
	if err != nil {
		panic(err)
	}
}

func readTlvTypeAndLength(reader io.Reader) (sendAttribute, uint16) {
	var tlvType sendAttribute
	var tlvLength uint16
	readAndPanic(reader, &tlvType)
	readAndPanic(reader, &tlvLength)
	return tlvType, tlvLength
}

func readString(reader io.Reader, length uint16) string {
	data := make([]byte, length)
	readAndPanic(reader, data)
	return string(data)
}

func isStdinPipeConnected() bool {
	stat, _ := os.Stdin.Stat()
	return stat.Mode()&os.ModeNamedPipe != 0
}
