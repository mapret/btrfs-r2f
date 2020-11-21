package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func renameCommand(reader io.Reader, config Config) bool {
	// First BTRFS_SEND_A_PATH
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	oldName := readString(reader, tlvLength)

	// Followed by BTRFS_SEND_A_PATH_TO
	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH_TO {
		panic("Unexpected command")
	}
	newName := readString(reader, tlvLength)

	err := os.Rename(path.Join(config.root, oldName), path.Join(config.root, newName))
	if err != nil {
		panic(err)
	}

	fmt.Printf("rename %s to %s\n", oldName, newName)
	return true
}
