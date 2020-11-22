package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func linkCommand(reader io.Reader, config Config) bool {
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	linkName := readString(reader, tlvLength)

	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH_LINK {
		panic("Unexpected command")
	}
	linkTarget := readString(reader, tlvLength)

	if !config.dryRun {
		err := os.Link(path.Join(config.root, linkTarget), path.Join(config.root, linkName))
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("link %s to %s\n", linkName, linkTarget)
	return true
}
