package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

func symlinkCommand(reader io.Reader, config Config) bool {
	tlvType, tlvLength := readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH {
		panic("Unexpected command")
	}
	linkName := readString(reader, tlvLength)

	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_INO {
		panic("Unexpected command")
	}
	var inodeNumber uint64
	readAndPanic(reader, &inodeNumber)

	tlvType, tlvLength = readTlvTypeAndLength(reader)
	if tlvType != BTRFS_SEND_A_PATH_LINK {
		panic("Unexpected command")
	}
	linkTarget := readString(reader, tlvLength)

	if !config.DryRun {
		if runtime.GOOS == "windows" {
			// Workaround for Windows: Create shortcut (via powershell), since administrator privileges are required
			// to creating a symlink
			command := exec.Command("powershell")
			command.Dir = config.Root

			absoluteTargetPath, _ := filepath.Abs(path.Join(config.Root, linkTarget))
			buffer := bytes.Buffer{}
			buffer.WriteString("$WshShell = New-Object -comObject WScript.Shell\n")
			buffer.WriteString(fmt.Sprintf("$Shortcut = $WshShell.CreateShortcut('%s.lnk')\n", linkName))
			buffer.WriteString(fmt.Sprintf("$Shortcut.TargetPath = '%s'\n", absoluteTargetPath))
			buffer.WriteString("$Shortcut.Save()\n")
			buffer.WriteString("exit\n")
			command.Stdin = &buffer

			err := command.Run()
			if err != nil {
				panic(err)
			}
		} else {
			err := os.Symlink(linkTarget, path.Join(config.Root, linkName))
			if err != nil {
				panic(err)
			}
		}
	}

	if config.Verbose {
		_, err := fmt.Fprintf(config.Stdout, "symlink %s to %s (%d)\n", linkName, linkTarget, inodeNumber)
		return err == nil
	}
	return true
}
