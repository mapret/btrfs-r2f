package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

/*const testFolder = "TestFilesystemRoot"

func TestMain(m *testing.M) {
	cleanTestFolder()
	code := m.Run()
	os.Exit(code)
}*/

func testCommandAndExpectStdout(t *testing.T, command func(io.Reader, Config) bool, data []byte, expectedResult string) {
	stdout := bytes.NewBufferString("")
	dataReader := bytes.NewReader(data)
	config := Config{
		DryRun:  true,
		Verbose: true,
		Stdout:  stdout,
	}

	returnValue := command(dataReader, config)
	if !returnValue {
		t.Fatalf("Command returned false")
	}
	if dataReader.Len() != 0 {
		t.Fatalf("Data stream not empty")
	}
	stdoutValue := stdout.String()
	if expectedResult != stdoutValue {
		extraNewline := "\n"
		if strings.HasSuffix(expectedResult, "\n") && strings.HasSuffix(stdoutValue, "\n") {
			extraNewline = ""
		}
		t.Fatalf("command stdout mismatch: expected\n\t%s%sbut was\n\t%s%s", expectedResult, extraNewline, stdoutValue, extraNewline)
	}
}

func TestLinkCommand(t *testing.T) {
	data := []byte("\x0f\x00\x04\x00dir1\x11\x00\x04\x00dir2")
	testCommandAndExpectStdout(t, linkCommand, data, "link dir1 to dir2\n")
}

func TestMkdirCommand(t *testing.T) {
	data := []byte("\x0f\x00\x04\x00dir1\x03\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01")
	testCommandAndExpectStdout(t, mkdirCommand, data, "mkdir dir1 (81985529216486895)\n")
}

func TestMkfileCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1\x03\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01")
	testCommandAndExpectStdout(t, mkfileCommand, data, "mkfile file1 (81985529216486895)\n")
}

func TestRenameCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1\x10\x00\x07\x00newname")
	testCommandAndExpectStdout(t, renameCommand, data, "rename file1 to newname\n")
}

func TestRmdirCommand(t *testing.T) {
	data := []byte("\x0f\x00\x04\x00dir1")
	testCommandAndExpectStdout(t, rmdirCommand, data, "rmdir dir1\n")
}

func TestSymlinkCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1\x03\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01\x11\x00\x05\x00file2")
	testCommandAndExpectStdout(t, symlinkCommand, data, "symlink file1 to file2 (81985529216486895)\n")
}

func TestTruncateCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1\x04\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01")
	testCommandAndExpectStdout(t, truncateCommand, data, "truncate file1 to 81985529216486895 bytes\n")
}

func TestUnlinkCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1")
	testCommandAndExpectStdout(t, unlinkCommand, data, "unlink file1\n")
}

func TestWriteCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1\x12\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01\x13\x00\x07\x00\xcd\xab\x89\x67\x45\x23\x01")
	testCommandAndExpectStdout(t, writeCommand, data, "write file1 (offset 81985529216486895, datalen 7)\n")
}
