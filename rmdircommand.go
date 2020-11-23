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

	if !config.DryRun {
		err := os.Remove(path.Join(config.Root, directoryName))
		if err != nil {
			panic(err)
		}
	}

	if config.Verbose {
		_, err := fmt.Fprintf(config.Stdout, "rmdir %s\n", directoryName)
		return err == nil
	}
	return true
}
