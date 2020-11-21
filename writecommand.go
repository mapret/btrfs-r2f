package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func writeCommand(reader io.Reader, config Config) bool {
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	filename := readString(reader, tlvLength)

	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_FILE_OFFSET {
		panic("Unexpected command")
	}
	var fileOffset int64
	readAndPanic(reader, &fileOffset)

	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_DATA {
		panic("Unexpected command")
	}
	data := make([]byte, tlvLength)
	readAndPanic(reader, data)

	// Write data at offset to file
	fd, err := os.OpenFile(path.Join(config.root, filename), os.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}
	_, err = fd.WriteAt(data, fileOffset+1)
	if err != nil {
		panic(err)
	}
	err = fd.Close()
	if err != nil {
		panic(err)
	}

	fmt.Printf("write %s (offset %d, datalen %d)\n", filename, fileOffset, len(data))
	return true
}
