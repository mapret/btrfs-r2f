package btrfsr2f

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

	if !config.DryRun {
		err := os.Remove(path.Join(config.Root, filename))
		if err != nil {
			panic(err)
		}
	}

	if config.Verbose {
		_, err := fmt.Fprintf(config.Stdout, "unlink %s\n", filename)
		return err == nil
	}
	return true
}
