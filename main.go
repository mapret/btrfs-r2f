package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const magicString = "btrfs-stream\000"

func main() {
	config := readConfig()

	if isStdinPipeConnected() {
		ExecuteProgram(os.Stdin, config)
	} else {
		fd, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)
		if err != nil {
			panic(err)
		}
		ExecuteProgram(fd, config)
		_ = fd.Close()
	}
}

func ExecuteProgram(reader io.Reader, config Config) {
	buffer := make([]byte, 13)
	_, err := io.ReadFull(reader, buffer)
	if err != nil {
		panic(err)
	}
	if string(buffer) != magicString {
		panic(fmt.Sprintf("Magic string mismatch, was \"%s\" but should be \"%s\"", string(buffer), magicString))
	}

	var sendVersion uint32
	err = binary.Read(reader, binary.LittleEndian, &sendVersion)
	if err != nil {
		panic(err)
	}
	if sendVersion != 1 {
		panic(fmt.Sprintf("Illegal send version, was %d but only %d is allowed", sendVersion, 1))
	}

	for readCommand(reader, config) {
	}
}

// https://en.wikipedia.org/wiki/Box-drawing_character
// https://en.wikipedia.org/wiki/Code_page_437
// Box drawing characters for reference
// ┌─────┬·┐
// │     │ ·
// ├─────┼·┤
// └─────┴·┘

func readCommand(reader io.Reader, config Config) bool {
	//                            Send command header
	// ┌────────────────────────────┬──────────────┬────────────────────────────┐
	// │       Command size         │ Command type │           CRC32            │
	// │          uint32            │    uint16    │          byte[4]           │
	// └────────────────────────────┴──────────────┴────────────────────────────┘
	var commandSize, crc32 uint32
	var commandType sendCommand
	readAndPanic(reader, &commandSize)
	readAndPanic(reader, &commandType)
	readAndPanic(reader, &crc32)
	//fmt.Printf("Command: %d\n", commandType)

	//                    Send attribute header
	// ┌──────────────┬──────────────┬───────··············───────┐
	// │    Type      │    Length    │            Data            │
	// │   uint16     │    uint16    │          {Length}          │
	// └──────────────┴──────────────┴───────··············───────┘

	if commandType == BTRFS_SEND_C_END {
		return false
	} else if commandType == BTRFS_SEND_C_MKFILE {
		return mkfileCommand(reader, config)
	} else if commandType == BTRFS_SEND_C_RENAME {
		return renameCommand(reader, config)
	} else if commandType == BTRFS_SEND_C_WRITE {
		return writeCommand(reader, config)
	} else if commandType == BTRFS_SEND_C_MKDIR {
		return mkdirCommand(reader, config)
	} else if commandType == BTRFS_SEND_C_UNLINK {
		return unlinkCommand(reader, config)
	} else if commandType == BTRFS_SEND_C_RMDIR {
		return rmdirCommand(reader, config)
	} else if commandType == BTRFS_SEND_C_SYMLINK {
		return symlinkCommand(reader, config)
	} else if commandType == BTRFS_SEND_C_LINK {
		return linkCommand(reader, config)
	} else if commandType == BTRFS_SEND_C_TRUNCATE {
		return truncateCommand(reader, config)
	}

	data := make([]byte, commandSize)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		panic(err)
	}

	return true
}
