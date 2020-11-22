package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func mkdirCommand(reader io.Reader, config Config) bool {
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	directoryName := readString(reader, tlvLength)

	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_INO {
		panic("Unexpected command")
	}
	var inodeNumber uint64
	readAndPanic(reader, &inodeNumber)

	if !config.dryRun {
		err := os.Mkdir(path.Join(config.root, directoryName), os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("mkdir %s (%d)\n", directoryName, inodeNumber)
	return true
}
