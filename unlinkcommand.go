package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func unlinkCommand(reader io.Reader, config Config) bool {
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	filename := readString(reader, tlvLength)

	if !config.dryRun {
		err := os.Remove(path.Join(config.root, filename))
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("unlink %s\n", filename)
	return true
}
