package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func truncateCommand(reader io.Reader, config Config) bool {
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	filename := readString(reader, tlvLength)

	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_SIZE {
		panic("Unexpected command")
	}
	var size int64
	readAndPanic(reader, &size)

	if !config.DryRun {
		err := os.Truncate(path.Join(config.Root, filename), size)
		if err != nil {
			panic(err)
		}
	}

	if config.Verbose {
		_, err := fmt.Fprintf(config.Stdout, "truncate %s to %d bytes\n", filename, size)
		return err == nil
	}
	return true
}
