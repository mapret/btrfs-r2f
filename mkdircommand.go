package btrfsr2f

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

	if !config.DryRun {
		err := os.Mkdir(path.Join(config.Root, directoryName), 0700)
		if err != nil {
			panic(err)
		}
	}

	if config.Verbose {
		_, err := fmt.Fprintf(config.Stdout, "mkdir %s (%d)\n", directoryName, inodeNumber)
		return err == nil
	}
	return true
}
