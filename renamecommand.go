package btrfsr2f

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
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

	if !config.DryRun {
		oldPath := path.Join(config.Root, oldName)
		newPath := path.Join(config.Root, newName)

		if runtime.GOOS == "windows" {
			// Workaround for Windows: Shortcuts have the additional extension ".lnk", which is not shown in
			// Windows Explorer and also not present in the btrfs-send command stream
			_, err := os.Stat(oldPath)
			if os.IsNotExist(err) {
				oldPath += ".lnk"
				newPath += ".lnk"
			}
		}

		err := os.Rename(oldPath, newPath)
		if err != nil {
			panic(err)
		}
	}

	if config.Verbose {
		_, err := fmt.Fprintf(config.Stdout, "rename %s to %s\n", oldName, newName)
		return err == nil
	}
	return true
}
