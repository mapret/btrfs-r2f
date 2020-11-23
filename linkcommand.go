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

	if !config.DryRun {
		err := os.Link(path.Join(config.Root, linkTarget), path.Join(config.Root, linkName))
		if err != nil {
			panic(err)
		}
	}

	if config.Verbose {
		_, err := fmt.Fprintf(config.Stdout, "link %s to %s\n", linkName, linkTarget)
		return err == nil
	}
	return true
}
