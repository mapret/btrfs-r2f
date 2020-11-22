package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func rmdirCommand(reader io.Reader, config Config) bool {
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	directoryName := readString(reader, tlvLength)

	if !config.dryRun {
		err := os.Remove(path.Join(config.root, directoryName))
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("rmdir %s\n", directoryName)
	return true
}
