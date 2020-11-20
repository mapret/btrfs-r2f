package main

import (
	"fmt"
	"io"
)

func mkfileCommand(reader io.Reader) bool {
	// First BTRFS_SEND_A_PATH
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	filename := readString(reader, tlvLength)

	// Followed by BTRFS_SEND_A_INO
	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_INO || tlvLength != 8 {
		panic("Unexpected command")
	}
	var inodeNumber uint64
	readAndPanic(reader, &inodeNumber)

	fmt.Printf("mkfile %s (%d)\n", filename, inodeNumber)
	return true
}
