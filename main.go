package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const magicString = "btrfs-stream\000"

func main() {
	filename := os.Args[1]

	fd, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	doStuff(fd)
	_ = fd.Close()
}

func doStuff(reader io.Reader) {
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

	for readCommand(reader) {
	}
}

// https://en.wikipedia.org/wiki/Box-drawing_character
// https://en.wikipedia.org/wiki/Code_page_437
// Box drawing characters for reference
// ┌─────┬·┐
// │     │ ·
// ├─────┼·┤
// └─────┴·┘

func readCommand(reader io.Reader) bool {
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

	fmt.Printf("Command: %d\n", commandType)

	//                    Send attribute header
	// ┌──────────────┬──────────────┬───────··············───────┐
	// │    Type      │    Length    │            Data            │
	// │   uint16     │    uint16    │          {Length}          │
	// └──────────────┴──────────────┴───────··············───────┘

	if commandType == BTRFS_SEND_C_END {
		return false
	} else if commandType == BTRFS_SEND_C_MKFILE {
		return mkfileCommand(reader)
	} else if commandType == BTRFS_SEND_C_RENAME {
		return renameCommand(reader)
	}

	data := make([]byte, commandSize)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		panic(err)
	}

	return true
}
