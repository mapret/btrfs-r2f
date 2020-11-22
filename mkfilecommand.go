package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func mkfileCommand(reader io.Reader, config Config) bool {
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

	if !config.dryRun {
		// Create empty file
		emptyFile, err := os.Create(path.Join(config.root, filename))
		if err != nil {
			panic(err)
		}
		err = emptyFile.Close()
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("mkfile %s (%d)\n", filename, inodeNumber)
	return true
}
