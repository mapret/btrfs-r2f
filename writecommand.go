package main

import (
	"fmt"
	"io"
)

func writeCommand(reader io.Reader) bool {
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	filename := readString(reader, tlvLength)

	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_FILE_OFFSET {
		panic("Unexpected command")
	}
	var fileOffset uint64
	readAndPanic(reader, &fileOffset)

	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_DATA {
		panic("Unexpected command")
	}
	data := make([]byte, tlvLength)
	readAndPanic(reader, data)

	fmt.Printf("write %s (offset %d, datalen %d)\n", filename, fileOffset, len(data))
	return true
}
