package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
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
		Root:    testFolder,
		DryRun:  false,
		Verbose: true,
		Stdout:  stdout,
	}
	prepareTestFolder(t, false)

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
	data := []byte("\x0f\x00\x05\x00file2\x11\x00\x05\x00file1")
	testCommandAndExpectStdout(t, linkCommand, data, "link file2 to file1\n")

	stat1, _ := os.Stat(path.Join(testFolder, "file1"))
	stat2, _ := os.Stat(path.Join(testFolder, "file2"))
	if !os.SameFile(stat1, stat2) {
		t.Fatal("file1 and file2 are not the same")
	}
}

func TestMkdirCommand(t *testing.T) {
	data := []byte("\x0f\x00\x04\x00dir2\x03\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01")
	testCommandAndExpectStdout(t, mkdirCommand, data, "mkdir dir2 (81985529216486895)\n")

	stat, err := os.Stat(path.Join(testFolder, "dir2"))
	if os.IsNotExist(err) || !stat.IsDir() {
		t.Fatal("dir2 was not created")
	}
}

func TestMkfileCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file2\x03\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01")
	testCommandAndExpectStdout(t, mkfileCommand, data, "mkfile file2 (81985529216486895)\n")

	stat, err := os.Stat(path.Join(testFolder, "file2"))
	if os.IsNotExist(err) || stat.IsDir() {
		t.Fatal("file2 was not created")
	}
}

func TestRenameCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1\x10\x00\x07\x00newname")
	testCommandAndExpectStdout(t, renameCommand, data, "rename file1 to newname\n")

	stat, err := os.Stat(path.Join(testFolder, "newname"))
	if os.IsNotExist(err) || stat.IsDir() {
		t.Fatal("file1 was not renamed to newname")
	}
}

func TestRmdirCommand(t *testing.T) {
	data := []byte("\x0f\x00\x04\x00dir1")
	testCommandAndExpectStdout(t, rmdirCommand, data, "rmdir dir1\n")

	_, err := os.Stat(path.Join(testFolder, "newname"))
	if !os.IsNotExist(err) {
		t.Fatal("dir1 was not deleted")
	}
}

func TestSymlinkCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file2\x03\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01\x11\x00\x05\x00file1")
	testCommandAndExpectStdout(t, symlinkCommand, data, "symlink file2 to file1 (81985529216486895)\n")

	if runtime.GOOS == "windows" {
		_, err := os.Stat(path.Join(testFolder, "file2.lnk"))
		if os.IsNotExist(err) {
			t.Fatal("file2 was not created")
		}
	} else {
		stat, err := os.Lstat(path.Join(testFolder, "file2"))
		if os.IsNotExist(err) {
			t.Fatal("file2 was not created")
		}
		if stat.Mode()&os.ModeSymlink == 0 {
			t.Fatal("file2 is not a symlink")
		}
	}
}

func TestTruncateCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1\x04\x00\x08\x00\x10\x00\x00\x00\x00\x00\x00\x00")
	testCommandAndExpectStdout(t, truncateCommand, data, "truncate file1 to 16 bytes\n")

	stat, _ := os.Stat(path.Join(testFolder, "file1"))
	if stat.Size() != 16 {
		t.Fatal("file1 was not truncated")
	}
	file1bytes, _ := ioutil.ReadFile(path.Join(testFolder, "file1"))
	if string(file1bytes) != "The quick brown " {
		t.Fatal("file1 was not truncated correctly")
	}
}

func TestUnlinkCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1")
	testCommandAndExpectStdout(t, unlinkCommand, data, "unlink file1\n")

	_, err := os.Stat(path.Join(testFolder, "file1"))
	if !os.IsNotExist(err) {
		t.Fatal("file1 was not deleted")
	}
}

func TestWriteCommand(t *testing.T) {
	data := []byte("\x0f\x00\x05\x00file1\x12\x00\x08\x00\x04\x00\x00\x00\x00\x00\x00\x00\x13\x00\x05\x00QUICK")
	testCommandAndExpectStdout(t, writeCommand, data, "write file1 (offset 4, datalen 5)\n")

	stat, _ := os.Stat(path.Join(testFolder, "file1"))
	if stat.Size() != int64(len(quickBrownFox)) {
		t.Fatal("size of file1 changed")
	}
	file1bytes, _ := ioutil.ReadFile(path.Join(testFolder, "file1"))
	if string(file1bytes) != strings.ReplaceAll(quickBrownFox, "quick", "QUICK") {
		t.Fatal("file1 was not written to correctly")
	}
}
